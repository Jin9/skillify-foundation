package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"syscall"
	"time"

	commonconfig "gitlab.com/b2c-e-commerce-platform/platform/backend/common/config"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/logger"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/checkout/app/checkout/access"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/checkout/config"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/checkout/router"

	_ "embed"
	_ "time/tzdata"
)

const (
	gracefulShutdownDuration = 10 * time.Second
	serverReadHeaderTimeout  = 5 * time.Second
	serverReadTimeout        = 5 * time.Second
	serverWriteTimeout       = 10 * time.Second
	handlerTimeout           = serverWriteTimeout - 100*time.Millisecond
)

// go build -ldflags "-X main.commit=abc123"
var commit string

//go:embed VERSION
var version string

func init() {
	if os.Getenv("GOMAXPROCS") != "" {
		runtime.GOMAXPROCS(0)
	} else {
		runtime.GOMAXPROCS(1)
	}
	if os.Getenv("GOMEMLIMIT") != "" {
		debug.SetMemoryLimit(-1)
	}
}

func main() {
	cfg := config.C(commonconfig.Env)
	_ = logger.New(logger.AWSKeyReplacer, logger.OpenSearchKeyReplacer, logger.CensorReplacer)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	d, cleanupDeps := router.NewDeps(ctx, cfg)
	defer cleanupDeps()

	// Background janitor: sweep stale INFLIGHT rows every 60s
	go runJanitor(ctx, d)

	// Background cleanup: delete expired rows every 5 minutes
	go runExpiredCleanup(ctx, d)

	r := router.New(d, version, commit, handlerTimeout)

	srv := &http.Server{
		Addr:              ":" + cfg.Server.Port,
		Handler:           r,
		ReadHeaderTimeout: serverReadHeaderTimeout,
		ReadTimeout:       serverReadTimeout,
		WriteTimeout:      serverWriteTimeout,
		MaxHeaderBytes:    1 << 20,
	}

	go func() {
		<-ctx.Done()
		slog.Info("checkout: shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), gracefulShutdownDuration)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
	}()

	slog.Info("checkout: listening", "port", cfg.Server.Port)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		slog.Error("checkout: ListenAndServe error", "err", err)
	}
	slog.Info("checkout: bye")
}

// runJanitor sweeps INFLIGHT rows older than 60s every 60s and marks them ABANDONED.
// Per td.json §idempotency_flow.janitor.
func runJanitor(ctx context.Context, d router.Deps) {
	store := access.NewIdempotencyStorage()
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	abandonedEnvelope, _ := json.Marshal(wrapper.Response[any]{
		Code:    wrapper.CodeTimeout,
		Message: "request timed out; retry with the same Idempotency-Key",
		Data:    nil,
	})

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			n, err := store.MarkStaleInflightAbandoned(ctx, d.Pool(), json.RawMessage(abandonedEnvelope), http.StatusGatewayTimeout)
			if err != nil {
				slog.ErrorContext(ctx, "checkout.janitor.error", slog.String("err", err.Error()))
			} else if n > 0 {
				slog.InfoContext(ctx, "checkout.janitor.abandoned", slog.Int64("count", n))
			}
		}
	}
}

// runExpiredCleanup deletes TTL-expired rows every 5 minutes.
func runExpiredCleanup(ctx context.Context, d router.Deps) {
	store := access.NewIdempotencyStorage()
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			n, err := store.CleanupExpired(ctx, d.Pool())
			if err != nil {
				slog.ErrorContext(ctx, "checkout.cleanup.error", slog.String("err", err.Error()))
			} else if n > 0 {
				slog.InfoContext(ctx, "checkout.cleanup.deleted", slog.Int64("count", n))
			}
		}
	}
}
