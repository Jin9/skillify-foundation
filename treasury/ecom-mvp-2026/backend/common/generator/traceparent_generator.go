package generator

import (
	crand "crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"strconv"
	"strings"
	"sync"
)

var (
	cryptoRandRead = crand.Read
	chaChaRead     = func(c *rand.ChaCha8, b []byte) (int, error) { return c.Read(b) }
)

// W3C Trace Context constants
const (
	TraceSupportedVersion = 0
	TraceFlagsSampled     = 1 // W3C spec: bit 0 (rightmost) for sampled flag
)

// W3C Trace Context types
type TraceID [16]byte
type SpanID [8]byte
type TraceFlags byte

func (id TraceID) String() string {
	return hex.EncodeToString(id[:])
}

func (id SpanID) String() string {
	return hex.EncodeToString(id[:])
}

// IsValid checks if TraceID is not all zeros.
func (id TraceID) IsValid() bool {
	for _, b := range id {
		if b != 0 {
			return true
		}
	}
	return false
}

// IsValid checks if SpanID is not all zeros.
func (id SpanID) IsValid() bool {
	for _, b := range id {
		if b != 0 {
			return true
		}
	}
	return false
}

// TraceParent represents a W3C-compliant trace context.
type TraceParent struct {
	TraceID    TraceID
	SpanID     SpanID
	TraceFlags TraceFlags
}

var (
	ErrEmptyTraceParent = errors.New("empty traceparent")
	ErrWrongTraceParent = errors.New("wrong traceparent")
	ErrInvalidTraceID   = errors.New("invalid trace ID (all zeros)")
	ErrInvalidSpanID    = errors.New("invalid span ID (all zeros)")
)

var emptyTrace = TraceParent{}

// Parse parses a W3C traceparent header value.
func Parse(parent string) (TraceParent, error) {
	if parent == "" {
		return emptyTrace, ErrEmptyTraceParent
	}

	token := strings.Split(parent, "-")
	if len(token) != 4 {
		return emptyTrace, ErrWrongTraceParent
	}
	if len(token[0]) != 2 {
		return emptyTrace, ErrWrongTraceParent
	}
	version, err := strconv.ParseUint(token[0], 16, 8)
	if err != nil {
		return emptyTrace, ErrWrongTraceParent
	}
	if int(version) != TraceSupportedVersion {
		return emptyTrace, ErrWrongTraceParent
	}

	traceIDBytes, err := hex.DecodeString(token[1])
	if err != nil {
		return emptyTrace, err
	}
	if len(traceIDBytes) != 16 {
		return emptyTrace, ErrWrongTraceParent
	}

	spanIDBytes, err := hex.DecodeString(token[2])
	if err != nil {
		return emptyTrace, err
	}
	if len(spanIDBytes) != 8 {
		return emptyTrace, ErrWrongTraceParent
	}

	traceFlagsBytes, err := hex.DecodeString(token[3])
	if err != nil {
		return emptyTrace, err
	}
	if len(traceFlagsBytes) != 1 {
		return emptyTrace, ErrWrongTraceParent
	}

	var sp TraceParent
	copy(sp.TraceID[:], traceIDBytes)
	copy(sp.SpanID[:], spanIDBytes)
	sp.TraceFlags = TraceFlags(traceFlagsBytes[0])
	if !sp.TraceID.IsValid() {
		return emptyTrace, ErrInvalidTraceID
	}
	if !sp.SpanID.IsValid() {
		return emptyTrace, ErrInvalidSpanID
	}

	return sp, nil
}

// =============================================================================
// OPTIMIZED ID GENERATOR (SINGLETON PATTERN)
// =============================================================================

type idGenerator struct {
	sync.Mutex
	randSource *rand.ChaCha8
}

var (
	globalGenerator *idGenerator
	generatorOnce   sync.Once
)

func getIDGenerator() *idGenerator {
	generatorOnce.Do(func() {
		globalGenerator = &idGenerator{}
		var seed [32]byte
		if _, err := cryptoRandRead(seed[:]); err != nil {
			slog.Error("failed to seed global random generator", "error", err)
			seed[0] = 1 // non-zero deterministic fallback
		}
		globalGenerator.randSource = rand.NewChaCha8(seed)
	})
	return globalGenerator
}

// NewSpanID generates a W3C-compliant span ID
func (gen *idGenerator) NewSpanID() SpanID {
	gen.Lock()
	defer gen.Unlock()

	var sid SpanID
	for {
		_, err := chaChaRead(gen.randSource, sid[:])
		if err != nil {
			slog.Error("failed to generate span ID", "error", err)
			sid[7] = 1 // fallback
			break
		}
		if sid.IsValid() {
			break
		}
	}
	return sid
}

// NewTraceID generates a W3C-compliant trace ID
func (gen *idGenerator) NewTraceID() TraceID {
	gen.Lock()
	defer gen.Unlock()

	var tid TraceID
	for {
		_, err := chaChaRead(gen.randSource, tid[:])
		if err != nil {
			slog.Error("failed to generate trace ID", "error", err)
			tid[15] = 1 // fallback
			break
		}
		if tid.IsValid() {
			break
		}
	}
	return tid
}

// =============================================================================
// TRACE CONTEXT CREATION
// =============================================================================

func NewTraceParent() TraceParent {
	return NewTraceParentWithFlags(TraceFlagsSampled)
}

func NewTraceParentWithFlags(flags TraceFlags) TraceParent {
	gen := getIDGenerator()
	return TraceParent{
		TraceID:    gen.NewTraceID(),
		SpanID:     gen.NewSpanID(),
		TraceFlags: flags,
	}
}

// =============================================================================
// TRACE CONTEXT METHODS
// =============================================================================

func (tp TraceParent) String() string {
	return fmt.Sprintf("%02x-%s-%s-%02x",
		TraceSupportedVersion,
		tp.TraceID.String(),
		tp.SpanID.String(),
		tp.TraceFlags)
}

func (tp TraceParent) CreateChild() TraceParent {
	gen := getIDGenerator()
	return TraceParent{
		TraceID:    tp.TraceID,
		SpanID:     gen.NewSpanID(),
		TraceFlags: tp.TraceFlags,
	}
}
