package token

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"maps"
	"math/big"
	"strings"
	"time"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/validator"
)

type Alg string

const (
	RS256 Alg = "RS256"
	ES256 Alg = "ES256"
	HS512 Alg = "HS512"
)

const (
	typ = "JWT"
)

// Claims holds standard JWT fields. Service-specific claims go in Extra and are
// serialised as top-level keys in the JWT payload so the token stays RFC 7519
// compliant and the token package stays free of domain concepts.
type Claims struct {
	Sub string `json:"sub"` // subject
	Iss string `json:"iss"` // issuer
	Aud string `json:"aud"` // audience
	Exp int64  `json:"exp"` // expiration time
	Iat int64  `json:"iat"` // issued at
	Jti string `json:"jti"` // unique identifier for the token

	// Extra holds caller-defined claims that are merged into the flat JWT
	// payload. Keys must not collide with the standard fields above.
	Extra map[string]any `json:"-"`
}

// MarshalJSON produces a flat JWT payload: standard fields + every entry in Extra.
func (c Claims) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"sub": c.Sub,
		"iss": c.Iss,
		"aud": c.Aud,
		"exp": c.Exp,
		"iat": c.Iat,
		"jti": c.Jti,
	}
	maps.Copy(m, c.Extra)
	return json.Marshal(m)
}

// UnmarshalJSON reads standard fields and collects remaining keys into Extra.
func (c *Claims) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	standard := map[string]func(json.RawMessage) error{
		"sub": func(b json.RawMessage) error { return json.Unmarshal(b, &c.Sub) },
		"iss": func(b json.RawMessage) error { return json.Unmarshal(b, &c.Iss) },
		"aud": func(b json.RawMessage) error { return json.Unmarshal(b, &c.Aud) },
		"exp": func(b json.RawMessage) error { return json.Unmarshal(b, &c.Exp) },
		"iat": func(b json.RawMessage) error { return json.Unmarshal(b, &c.Iat) },
		"jti": func(b json.RawMessage) error { return json.Unmarshal(b, &c.Jti) },
	}

	for k, v := range raw {
		if fn, ok := standard[k]; ok {
			if err := fn(v); err != nil {
				return fmt.Errorf("claims field %q: %w", k, err)
			}
			continue
		}
		if c.Extra == nil {
			c.Extra = make(map[string]any)
		}
		var val any
		if err := json.Unmarshal(v, &val); err != nil {
			return fmt.Errorf("claims extra field %q: %w", k, err)
		}
		c.Extra[k] = val
	}
	return nil
}

type JWTSignerConfig struct {
	PrivateKey string `validate:"required"`
	Alg        string `validate:"required"`

	Issuer   string        `validate:"required"`
	Audience string        `validate:"required"`
	Expire   time.Duration `validate:"required"`
}

type JWTSigner interface {
	SignRS256(payload Claims) (string, error)
	SignES256(payload Claims) (string, error)
}

var _ JWTSigner = (*jwtSigner)(nil)

type jwtSigner struct {
	rsaPrivateKey   *rsa.PrivateKey
	ecdsaPrivateKey *ecdsa.PrivateKey
	alg             string

	issuer   string
	audience string
	expire   time.Duration
}

func NewJWTSigner(config JWTSignerConfig) (JWTSigner, error) {
	if err := validator.Validate(config); err != nil {
		return nil, err
	}

	switch config.Alg {
	case string(RS256):
		privKey, err := loadRS256PrivateKey([]byte(config.PrivateKey))
		if err != nil {
			return nil, fmt.Errorf("load private key: %w", err)
		}
		return &jwtSigner{
			rsaPrivateKey: privKey,
			alg:           config.Alg,
			issuer:        config.Issuer,
			audience:      config.Audience,
			expire:        config.Expire,
		}, nil
	case string(ES256):
		privKey, err := loadES256PrivateKey([]byte(config.PrivateKey))
		if err != nil {
			return nil, fmt.Errorf("load private key: %w", err)
		}
		return &jwtSigner{
			ecdsaPrivateKey: privKey,
			alg:             config.Alg,
			issuer:          config.Issuer,
			audience:        config.Audience,
			expire:          config.Expire,
		}, nil
	default:
		return nil, fmt.Errorf("invalid algorithm: %s", config.Alg)
	}
}

// MustNewJWTSigner is a convenience wrapper that panics on error.
func MustNewJWTSigner(config JWTSignerConfig) JWTSigner {
	s, err := NewJWTSigner(config)
	if err != nil {
		panic(err)
	}
	return s
}

func (j *jwtSigner) buildSigningInput(payload *Claims) string {
	header := map[string]string{
		"alg": j.alg,
		"typ": typ,
	}

	now := time.Now()
	payload.Iat = now.Unix()
	payload.Exp = now.Add(j.expire).Unix()
	payload.Iss = j.issuer
	payload.Aud = j.audience

	headerJSON, _ := json.Marshal(header)
	payloadJSON, _ := json.Marshal(payload)

	return base64URLEncode(headerJSON) + "." + base64URLEncode(payloadJSON)
}

func (j *jwtSigner) SignRS256(payload Claims) (string, error) {
	if j.rsaPrivateKey == nil {
		return "", errors.New("private key not set")
	}

	signingInput := j.buildSigningInput(&payload)

	hashed := sha256.Sum256([]byte(signingInput))
	sig, err := rsa.SignPKCS1v15(rand.Reader, j.rsaPrivateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", fmt.Errorf("sign error: %w", err)
	}

	return signingInput + "." + base64URLEncode(sig), nil
}

func (j *jwtSigner) SignES256(payload Claims) (string, error) {
	if j.ecdsaPrivateKey == nil {
		return "", errors.New("private key not set")
	}

	signingInput := j.buildSigningInput(&payload)

	// Sign using SHA256
	hash := sha256.Sum256([]byte(signingInput))
	r, s, err := ecdsa.Sign(rand.Reader, j.ecdsaPrivateKey, hash[:])
	if err != nil {
		return "", err
	}

	// Encode R||S as 64-byte raw signature
	rBytes := r.Bytes()
	sBytes := s.Bytes()
	rPadded := append(make([]byte, 32-len(rBytes)), rBytes...)
	sPadded := append(make([]byte, 32-len(sBytes)), sBytes...)
	signature := append(rPadded, sPadded...)

	return signingInput + "." + base64URLEncode(signature), nil
}

type JWTVerifierConfig struct {
	PublicKey string `validate:"required"`
	Alg       string `validate:"required"`
}

type JWTVerifier interface {
	VerifySignatureRS256(token string) error
	VerifySignatureES256(token string) error
	VerifySignatureHS512(token string) error
}

var _ JWTVerifier = (*jwtVerifier)(nil)

type jwtVerifier struct {
	rsaPublicKey   *rsa.PublicKey
	ecdsaPublicKey *ecdsa.PublicKey
	hmacKey        []byte
}

func NewJWTVerifier(config JWTVerifierConfig) (JWTVerifier, error) {
	if err := validator.Validate(config); err != nil {
		return nil, err
	}

	switch config.Alg {
	case string(RS256):
		pubKey, err := loadRS256PublicKey([]byte(config.PublicKey))
		if err != nil {
			return nil, fmt.Errorf("load public key: %w", err)
		}
		return &jwtVerifier{
			rsaPublicKey:   pubKey,
			ecdsaPublicKey: nil,
			hmacKey:        nil,
		}, nil
	case string(ES256):
		pubKey, err := loadES256PublicKey([]byte(config.PublicKey))
		if err != nil {
			return nil, fmt.Errorf("load public key: %w", err)
		}
		return &jwtVerifier{
			rsaPublicKey:   nil,
			ecdsaPublicKey: pubKey,
			hmacKey:        nil,
		}, nil
	case string(HS512):
		return &jwtVerifier{
			rsaPublicKey:   nil,
			ecdsaPublicKey: nil,
			hmacKey:        []byte(config.PublicKey),
		}, nil
	default:
		return nil, fmt.Errorf("invalid algorithm: %s", config.Alg)
	}
}

// MustNewJWTVerifier is a convenience wrapper that panics on error.
func MustNewJWTVerifier(config JWTVerifierConfig) JWTVerifier {
	v, err := NewJWTVerifier(config)
	if err != nil {
		panic(err)
	}
	return v
}

func (j *jwtVerifier) VerifySignatureHS512(token string) error {
	if j.hmacKey == nil {
		return errors.New("HMAC key not set")
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return errors.New("invalid JWT format")
	}

	message := parts[0] + "." + parts[1]
	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return fmt.Errorf("invalid signature encoding: %w", err)
	}

	mac := hmac.New(sha512.New, j.hmacKey)
	mac.Write([]byte(message))
	expectedSignature := mac.Sum(nil)

	if !hmac.Equal(signature, expectedSignature) {
		return errors.New("invalid signature")
	}

	return nil
}

func (j *jwtVerifier) VerifySignatureRS256(token string) error {
	if j.rsaPublicKey == nil {
		return errors.New("public key not set")
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return errors.New("invalid JWT format")
	}

	signed := parts[0] + "." + parts[1]
	sig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return fmt.Errorf("invalid signature encoding: %w", err)
	}

	hashed := sha256.Sum256([]byte(signed))
	if err := rsa.VerifyPKCS1v15(j.rsaPublicKey, crypto.SHA256, hashed[:], sig); err != nil {
		return fmt.Errorf("invalid signature: %w", err)
	}

	return nil
}

func (j *jwtVerifier) VerifySignatureES256(token string) error {
	if j.ecdsaPublicKey == nil {
		return errors.New("public key not set")
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return errors.New("invalid token format")
	}

	message := parts[0] + "." + parts[1]
	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return fmt.Errorf("invalid signature encoding: %w", err)
	}
	if len(signature) != 64 {
		return errors.New("invalid signature length")
	}

	// Split R||S
	r := new(big.Int).SetBytes(signature[:32])
	s := new(big.Int).SetBytes(signature[32:])
	hash := sha256.Sum256([]byte(message))

	if valid := ecdsa.Verify(j.ecdsaPublicKey, hash[:], r, s); !valid {
		return errors.New("invalid signature")
	}

	return nil
}

type JWTParserConfig struct {
	Issuer   string `validate:"required"`
	Audience string `validate:"required"`
}

type JWTParser interface {
	ParseToken(token string) (*Claims, error)
	ValidateClaims(payload *Claims) error
}

var _ JWTParser = (*jwtParser)(nil)

type jwtParser struct {
	issuer   string
	audience string
}

func NewJWTParser(config JWTParserConfig) (JWTParser, error) {
	if err := validator.Validate(config); err != nil {
		return nil, err
	}
	return &jwtParser{
		issuer:   config.Issuer,
		audience: config.Audience,
	}, nil
}

// MustNewJWTParser is a convenience wrapper that panics on error.
func MustNewJWTParser(config JWTParserConfig) JWTParser {
	p, err := NewJWTParser(config)
	if err != nil {
		panic(err)
	}
	return p
}

func (j *jwtParser) ParseToken(token string) (*Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid JWT format")
	}

	payloadSegment := parts[1]
	payloadBytes, err := base64.RawURLEncoding.DecodeString(payloadSegment)
	if err != nil {
		return nil, fmt.Errorf("base64 decode error: %w", err)
	}

	payload := &Claims{}
	if err := json.Unmarshal(payloadBytes, payload); err != nil {
		return nil, fmt.Errorf("JSON unmarshal error: %w", err)
	}

	return payload, nil
}

func (j *jwtParser) ValidateClaims(payload *Claims) error {
	now := time.Now().Unix()

	if payload.Exp < now {
		return errors.New("token expired")
	}
	if payload.Iat > now {
		return errors.New("token not yet valid")
	}
	if payload.Iss != j.issuer {
		return errors.New("invalid issuer")
	}
	if j.audience != "" && payload.Aud != j.audience {
		return errors.New("invalid audience")
	}
	// Add any other claim validations you need

	return nil
}

func loadRS256PrivateKey(data []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(data)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("invalid PEM block for private key")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

// LoadES256PrivateKey takes a PEM-encoded ECDSA private key (ES256) string and returns an *ecdsa.PrivateKey.
func loadES256PrivateKey(data []byte) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing the private key")
	}

	key, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse EC private key: %w", err)
	}

	// Optional: validate the key is actually using P-256
	if key.Curve.Params().Name != "P-256" {
		return nil, fmt.Errorf("invalid curve: expected P-256 but got %s", key.Curve.Params().Name)
	}

	return key, nil
}

func loadRS256PublicKey(data []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(data)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("invalid PEM block for public key")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}
	return rsaPub, nil
}

func loadES256PublicKey(data []byte) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing the public key")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	pubKey, ok := pubInterface.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("not an ECDSA public key")
	}

	// Optional: verify it's using P-256
	if pubKey.Curve.Params().Name != "P-256" {
		return nil, fmt.Errorf("invalid curve: expected P-256 but got %s", pubKey.Curve.Params().Name)
	}

	return pubKey, nil
}

func base64URLEncode(data []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(data), "=")
}
