package token

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"strings"
	"testing"
	"time"
)

// — Helpers —

func generateES256Keys(t *testing.T) (privPEM, pubPEM string) {
	t.Helper()
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate EC key: %v", err)
	}
	privBytes, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		t.Fatalf("marshal EC private key: %v", err)
	}
	privPEM = string(pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes}))

	pubBytes, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	if err != nil {
		t.Fatalf("marshal EC public key: %v", err)
	}
	pubPEM = string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes}))
	return
}

func generateRS256Keys(t *testing.T) (privPEM, pubPEM string) {
	t.Helper()
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate RSA key: %v", err)
	}
	privPEM = string(pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)}))

	pubBytes, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	if err != nil {
		t.Fatalf("marshal RSA public key: %v", err)
	}
	pubPEM = string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes}))
	return
}

// — Claims JSON —

func TestClaimsMarshalJSON(t *testing.T) {
	c := Claims{
		Sub:   "user-1",
		Iss:   "test",
		Aud:   "api",
		Exp:   9999999999,
		Iat:   1000000000,
		Jti:   "jti-1",
		Extra: map[string]any{"role": "admin"},
	}

	data, err := c.MarshalJSON()
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	s := string(data)
	for _, want := range []string{`"sub":"user-1"`, `"role":"admin"`, `"jti":"jti-1"`} {
		if !strings.Contains(s, want) {
			t.Errorf("expected %s in %s", want, s)
		}
	}
}

func TestClaimsUnmarshalJSON(t *testing.T) {
	raw := `{"sub":"u1","iss":"iss","aud":"aud","exp":99,"iat":10,"jti":"j1","role":"admin","org":"org1"}`

	var c Claims
	if err := c.UnmarshalJSON([]byte(raw)); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if c.Sub != "u1" {
		t.Errorf("expected sub=u1, got %s", c.Sub)
	}
	if c.Extra["role"] != "admin" {
		t.Errorf("expected extra role=admin, got %v", c.Extra["role"])
	}
	if c.Extra["org"] != "org1" {
		t.Errorf("expected extra org=org1, got %v", c.Extra["org"])
	}
}

func TestClaimsUnmarshalJSON_InvalidJSON(t *testing.T) {
	var c Claims
	if err := c.UnmarshalJSON([]byte(`not json`)); err == nil {
		t.Fatal("expected error")
	}
}

// — NewJWTSigner —

func TestNewJWTSigner_ES256(t *testing.T) {
	privPEM, _ := generateES256Keys(t)

	signer, err := NewJWTSigner(JWTSignerConfig{
		PrivateKey: privPEM,
		Alg:        "ES256",
		Issuer:     "test",
		Audience:   "api",
		Expire:     time.Hour,
	})
	if err != nil {
		t.Fatalf("new signer: %v", err)
	}
	if signer == nil {
		t.Fatal("expected non-nil signer")
	}
}

func TestNewJWTSigner_RS256(t *testing.T) {
	privPEM, _ := generateRS256Keys(t)

	signer, err := NewJWTSigner(JWTSignerConfig{
		PrivateKey: privPEM,
		Alg:        "RS256",
		Issuer:     "test",
		Audience:   "api",
		Expire:     time.Hour,
	})
	if err != nil {
		t.Fatalf("new signer: %v", err)
	}
	if signer == nil {
		t.Fatal("expected non-nil signer")
	}
}

func TestNewJWTSigner_InvalidAlg(t *testing.T) {
	_, err := NewJWTSigner(JWTSignerConfig{
		PrivateKey: "key",
		Alg:        "UNKNOWN",
		Issuer:     "test",
		Audience:   "api",
		Expire:     time.Hour,
	})
	if err == nil {
		t.Fatal("expected error for unknown alg")
	}
}

func TestNewJWTSigner_MissingConfig(t *testing.T) {
	_, err := NewJWTSigner(JWTSignerConfig{})
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestNewJWTSigner_InvalidKey(t *testing.T) {
	_, err := NewJWTSigner(JWTSignerConfig{
		PrivateKey: "not-a-pem-key",
		Alg:        "ES256",
		Issuer:     "test",
		Audience:   "api",
		Expire:     time.Hour,
	})
	if err == nil {
		t.Fatal("expected error for invalid key")
	}
}

func TestMustNewJWTSigner_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	MustNewJWTSigner(JWTSignerConfig{})
}

// — Sign + Verify roundtrip ES256 —

func TestSignAndVerifyES256(t *testing.T) {
	privPEM, pubPEM := generateES256Keys(t)

	signer, err := NewJWTSigner(JWTSignerConfig{
		PrivateKey: privPEM,
		Alg:        "ES256",
		Issuer:     "test-issuer",
		Audience:   "test-aud",
		Expire:     time.Hour,
	})
	if err != nil {
		t.Fatalf("new signer: %v", err)
	}

	token, err := signer.SignES256(Claims{Sub: "user-1", Jti: "jti-1"})
	if err != nil {
		t.Fatalf("sign: %v", err)
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Fatalf("expected 3 parts, got %d", len(parts))
	}

	verifier, err := NewJWTVerifier(JWTVerifierConfig{PublicKey: pubPEM, Alg: "ES256"})
	if err != nil {
		t.Fatalf("new verifier: %v", err)
	}

	if err := verifier.VerifySignatureES256(token); err != nil {
		t.Fatalf("verify: %v", err)
	}
}

func TestSignES256_NilKey(t *testing.T) {
	s := &jwtSigner{ecdsaPrivateKey: nil}
	_, err := s.SignES256(Claims{})
	if err == nil {
		t.Fatal("expected error for nil key")
	}
}

// — Sign + Verify roundtrip RS256 —

func TestSignAndVerifyRS256(t *testing.T) {
	privPEM, pubPEM := generateRS256Keys(t)

	signer, err := NewJWTSigner(JWTSignerConfig{
		PrivateKey: privPEM,
		Alg:        "RS256",
		Issuer:     "test-issuer",
		Audience:   "test-aud",
		Expire:     time.Hour,
	})
	if err != nil {
		t.Fatalf("new signer: %v", err)
	}

	token, err := signer.SignRS256(Claims{Sub: "user-1"})
	if err != nil {
		t.Fatalf("sign: %v", err)
	}

	verifier, err := NewJWTVerifier(JWTVerifierConfig{PublicKey: pubPEM, Alg: "RS256"})
	if err != nil {
		t.Fatalf("new verifier: %v", err)
	}

	if err := verifier.VerifySignatureRS256(token); err != nil {
		t.Fatalf("verify: %v", err)
	}
}

func TestSignRS256_NilKey(t *testing.T) {
	s := &jwtSigner{rsaPrivateKey: nil}
	_, err := s.SignRS256(Claims{})
	if err == nil {
		t.Fatal("expected error for nil key")
	}
}

// — NewJWTVerifier —

func TestNewJWTVerifier_HS512(t *testing.T) {
	v, err := NewJWTVerifier(JWTVerifierConfig{PublicKey: "my-secret-key", Alg: "HS512"})
	if err != nil {
		t.Fatalf("new verifier: %v", err)
	}
	if v == nil {
		t.Fatal("expected non-nil verifier")
	}
}

func TestNewJWTVerifier_InvalidAlg(t *testing.T) {
	_, err := NewJWTVerifier(JWTVerifierConfig{PublicKey: "key", Alg: "UNKNOWN"})
	if err == nil {
		t.Fatal("expected error for unknown alg")
	}
}

func TestNewJWTVerifier_MissingConfig(t *testing.T) {
	_, err := NewJWTVerifier(JWTVerifierConfig{})
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestNewJWTVerifier_InvalidKey(t *testing.T) {
	_, err := NewJWTVerifier(JWTVerifierConfig{PublicKey: "not-pem", Alg: "ES256"})
	if err == nil {
		t.Fatal("expected error for invalid key")
	}
}

func TestMustNewJWTVerifier_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	MustNewJWTVerifier(JWTVerifierConfig{})
}

// — Verify edge cases —

func TestVerifyES256_InvalidFormat(t *testing.T) {
	_, pubPEM := generateES256Keys(t)
	v, _ := NewJWTVerifier(JWTVerifierConfig{PublicKey: pubPEM, Alg: "ES256"})

	if err := v.VerifySignatureES256("not.valid"); err == nil {
		t.Fatal("expected error for 2-part token")
	}
}

func TestVerifyES256_InvalidSignatureEncoding(t *testing.T) {
	_, pubPEM := generateES256Keys(t)
	v, _ := NewJWTVerifier(JWTVerifierConfig{PublicKey: pubPEM, Alg: "ES256"})

	if err := v.VerifySignatureES256("a.b.!!!invalid!!!"); err == nil {
		t.Fatal("expected error for invalid base64 signature")
	}
}

func TestVerifyES256_InvalidSignatureLength(t *testing.T) {
	_, pubPEM := generateES256Keys(t)
	v, _ := NewJWTVerifier(JWTVerifierConfig{PublicKey: pubPEM, Alg: "ES256"})

	short := base64.RawURLEncoding.EncodeToString([]byte("short"))
	if err := v.VerifySignatureES256("a.b." + short); err == nil {
		t.Fatal("expected error for short signature")
	}
}

func TestVerifyES256_NilKey(t *testing.T) {
	v := &jwtVerifier{ecdsaPublicKey: nil}
	if err := v.VerifySignatureES256("a.b.c"); err == nil {
		t.Fatal("expected error for nil key")
	}
}

func TestVerifyRS256_InvalidFormat(t *testing.T) {
	_, pubPEM := generateRS256Keys(t)
	v, _ := NewJWTVerifier(JWTVerifierConfig{PublicKey: pubPEM, Alg: "RS256"})

	if err := v.VerifySignatureRS256("only-one-part"); err == nil {
		t.Fatal("expected error for 1-part token")
	}
}

func TestVerifyRS256_InvalidSig(t *testing.T) {
	_, pubPEM := generateRS256Keys(t)
	v, _ := NewJWTVerifier(JWTVerifierConfig{PublicKey: pubPEM, Alg: "RS256"})

	fakeSig := base64.RawURLEncoding.EncodeToString([]byte("fake"))
	if err := v.VerifySignatureRS256("a.b." + fakeSig); err == nil {
		t.Fatal("expected error for invalid signature")
	}
}

func TestVerifyRS256_NilKey(t *testing.T) {
	v := &jwtVerifier{rsaPublicKey: nil}
	if err := v.VerifySignatureRS256("a.b.c"); err == nil {
		t.Fatal("expected error for nil key")
	}
}

func TestVerifyHS512_NilKey(t *testing.T) {
	v := &jwtVerifier{hmacKey: nil}
	if err := v.VerifySignatureHS512("a.b.c"); err == nil {
		t.Fatal("expected error for nil key")
	}
}

func TestVerifyHS512_InvalidFormat(t *testing.T) {
	v, _ := NewJWTVerifier(JWTVerifierConfig{PublicKey: "secret", Alg: "HS512"})
	if err := v.VerifySignatureHS512("only.two"); err == nil {
		t.Fatal("expected error for 2-part token")
	}
}

func TestVerifyHS512_InvalidSignature(t *testing.T) {
	v, _ := NewJWTVerifier(JWTVerifierConfig{PublicKey: "secret", Alg: "HS512"})
	fakeSig := base64.RawURLEncoding.EncodeToString([]byte("wrong"))
	if err := v.VerifySignatureHS512("a.b." + fakeSig); err == nil {
		t.Fatal("expected error for wrong signature")
	}
}

func TestVerifyHS512_BadBase64(t *testing.T) {
	v, _ := NewJWTVerifier(JWTVerifierConfig{PublicKey: "secret", Alg: "HS512"})
	if err := v.VerifySignatureHS512("a.b.!!!"); err == nil {
		t.Fatal("expected error for bad base64")
	}
}

// — JWTParser —

func TestNewJWTParser(t *testing.T) {
	p, err := NewJWTParser(JWTParserConfig{Issuer: "iss", Audience: "aud"})
	if err != nil {
		t.Fatalf("new parser: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil parser")
	}
}

func TestNewJWTParser_MissingConfig(t *testing.T) {
	_, err := NewJWTParser(JWTParserConfig{})
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestMustNewJWTParser_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	MustNewJWTParser(JWTParserConfig{})
}

func TestParseToken(t *testing.T) {
	privPEM, _ := generateES256Keys(t)
	signer, _ := NewJWTSigner(JWTSignerConfig{
		PrivateKey: privPEM,
		Alg:        "ES256",
		Issuer:     "iss",
		Audience:   "aud",
		Expire:     time.Hour,
	})

	token, _ := signer.SignES256(Claims{Sub: "user-1", Jti: "jti-1", Extra: map[string]any{"role": "admin"}})

	parser, _ := NewJWTParser(JWTParserConfig{Issuer: "iss", Audience: "aud"})
	claims, err := parser.ParseToken(token)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if claims.Sub != "user-1" {
		t.Errorf("expected sub=user-1, got %s", claims.Sub)
	}
	if claims.Iss != "iss" {
		t.Errorf("expected iss=iss, got %s", claims.Iss)
	}
	if claims.Extra["role"] != "admin" {
		t.Errorf("expected extra role=admin, got %v", claims.Extra["role"])
	}
}

func TestParseToken_InvalidFormat(t *testing.T) {
	parser, _ := NewJWTParser(JWTParserConfig{Issuer: "iss", Audience: "aud"})
	_, err := parser.ParseToken("only-one-part")
	if err == nil {
		t.Fatal("expected error for invalid format")
	}
}

func TestParseToken_InvalidBase64(t *testing.T) {
	parser, _ := NewJWTParser(JWTParserConfig{Issuer: "iss", Audience: "aud"})
	_, err := parser.ParseToken("a.!!!.c")
	if err == nil {
		t.Fatal("expected error for bad base64")
	}
}

func TestParseToken_InvalidJSON(t *testing.T) {
	parser, _ := NewJWTParser(JWTParserConfig{Issuer: "iss", Audience: "aud"})
	notJSON := base64.RawURLEncoding.EncodeToString([]byte("not json"))
	_, err := parser.ParseToken("a." + notJSON + ".c")
	if err == nil {
		t.Fatal("expected error for bad JSON payload")
	}
}

// — ValidateClaims —

func TestValidateClaims(t *testing.T) {
	parser, _ := NewJWTParser(JWTParserConfig{Issuer: "iss", Audience: "aud"})

	t.Run("valid claims", func(t *testing.T) {
		c := &Claims{
			Sub: "user",
			Iss: "iss",
			Aud: "aud",
			Iat: time.Now().Unix() - 60,
			Exp: time.Now().Unix() + 3600,
		}
		if err := parser.ValidateClaims(c); err != nil {
			t.Fatalf("expected valid, got: %v", err)
		}
	})

	t.Run("expired", func(t *testing.T) {
		c := &Claims{Iss: "iss", Aud: "aud", Iat: 1000, Exp: 1001}
		if err := parser.ValidateClaims(c); err == nil {
			t.Fatal("expected expired error")
		}
	})

	t.Run("not yet valid (future iat)", func(t *testing.T) {
		c := &Claims{Iss: "iss", Aud: "aud", Iat: time.Now().Unix() + 9999, Exp: time.Now().Unix() + 99999}
		if err := parser.ValidateClaims(c); err == nil {
			t.Fatal("expected not-yet-valid error")
		}
	})

	t.Run("wrong issuer", func(t *testing.T) {
		c := &Claims{Iss: "wrong", Aud: "aud", Iat: time.Now().Unix() - 60, Exp: time.Now().Unix() + 3600}
		if err := parser.ValidateClaims(c); err == nil {
			t.Fatal("expected invalid issuer error")
		}
	})

	t.Run("wrong audience", func(t *testing.T) {
		c := &Claims{Iss: "iss", Aud: "wrong", Iat: time.Now().Unix() - 60, Exp: time.Now().Unix() + 3600}
		if err := parser.ValidateClaims(c); err == nil {
			t.Fatal("expected invalid audience error")
		}
	})
}

// — Key loading edge cases —

func TestLoadRS256PrivateKey_InvalidPEM(t *testing.T) {
	_, err := loadRS256PrivateKey([]byte("not pem"))
	if err == nil {
		t.Fatal("expected error for invalid PEM")
	}
}

func TestLoadES256PrivateKey_InvalidPEM(t *testing.T) {
	_, err := loadES256PrivateKey([]byte("not pem"))
	if err == nil {
		t.Fatal("expected error for invalid PEM")
	}
}

func TestLoadRS256PublicKey_InvalidPEM(t *testing.T) {
	_, err := loadRS256PublicKey([]byte("not pem"))
	if err == nil {
		t.Fatal("expected error for invalid PEM")
	}
}

func TestLoadES256PublicKey_InvalidPEM(t *testing.T) {
	_, err := loadES256PublicKey([]byte("not pem"))
	if err == nil {
		t.Fatal("expected error for invalid PEM")
	}
}

func TestLoadRS256PublicKey_WrongType(t *testing.T) {
	_, pubPEM := generateES256Keys(t)
	_, err := loadRS256PublicKey([]byte(pubPEM))
	if err == nil {
		t.Fatal("expected error for wrong key type under RS256 loader")
	}
}

func TestLoadES256PublicKey_WrongType(t *testing.T) {
	_, pubPEM := generateRS256Keys(t)
	_, err := loadES256PublicKey([]byte(pubPEM))
	if err == nil {
		t.Fatal("expected error for wrong key type under ES256 loader")
	}
}

func TestLoadES256PrivateKey_WrongType(t *testing.T) {
	privPEM, _ := generateRS256Keys(t)
	_, err := loadES256PrivateKey([]byte(privPEM))
	if err == nil {
		t.Fatal("expected error for wrong private key type under ES256 loader")
	}
}
