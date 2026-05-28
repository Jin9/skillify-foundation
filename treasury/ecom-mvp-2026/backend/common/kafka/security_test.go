package kafka

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"math/big"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/codec"
)

func TestApplyTLSConfig(t *testing.T) {
	t.Run("nil config", func(t *testing.T) {
		err := applyTLSConfig(nil, true, "ca", "", "")
		if err == nil {
			t.Fatalf("expected error for nil config")
		}
	})

	t.Run("disabled", func(t *testing.T) {
		config := newConfig()

		err := applyTLSConfig(config, false, "", "", "")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if config.Net.TLS.Enable {
			t.Fatalf("expected tls disabled")
		}
	})

	t.Run("missing ca cert", func(t *testing.T) {
		config := newConfig()

		err := applyTLSConfig(config, true, "", "", "")
		if err == nil {
			t.Fatalf("expected error for missing ca cert")
		}
	})

	t.Run("invalid ca cert", func(t *testing.T) {
		config := newConfig()

		err := applyTLSConfig(config, true, "invalid", "", "")
		if err == nil {
			t.Fatalf("expected error for invalid ca cert")
		}
	})

	t.Run("valid ca cert", func(t *testing.T) {
		config := newConfig()
		caPEM := createSelfSignedCertPEM(t)

		err := applyTLSConfig(config, true, caPEM, "", "")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if !config.Net.TLS.Enable {
			t.Fatalf("expected tls enabled")
		}
		if config.Net.TLS.Config == nil {
			t.Fatalf("expected tls config")
		}
	})

	t.Run("invalid client cert pair", func(t *testing.T) {
		config := newConfig()
		caPEM := createSelfSignedCertPEM(t)

		err := applyTLSConfig(config, true, caPEM, "invalid", "invalid")
		if err == nil {
			t.Fatalf("expected error for invalid client cert pair")
		}
	})
}

func TestApplySASLConfig(t *testing.T) {
	t.Run("nil config", func(t *testing.T) {
		err := applySASLConfig(nil, true, sarama.SASLTypePlaintext, "user", "pass")
		if err == nil {
			t.Fatalf("expected error for nil config")
		}
	})

	t.Run("plaintext", func(t *testing.T) {
		config := newConfig()

		err := applySASLConfig(config, true, sarama.SASLTypePlaintext, "user", "pass")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if !config.Net.SASL.Enable {
			t.Fatalf("expected sasl enabled")
		}
		if config.Net.SASL.User != "user" || config.Net.SASL.Password != "pass" {
			t.Fatalf("expected credentials to be set")
		}
		if config.Net.SASL.SCRAMClientGeneratorFunc != nil {
			t.Fatalf("expected scram generator to remain nil for plaintext")
		}
	})

	t.Run("scram sha 256", func(t *testing.T) {
		config := newConfig()

		err := applySASLConfig(config, true, sarama.SASLTypeSCRAMSHA256, "user", "pass")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if config.Net.SASL.SCRAMClientGeneratorFunc == nil {
			t.Fatalf("expected scram generator to be set")
		}

		client, ok := config.Net.SASL.SCRAMClientGeneratorFunc().(*SCRAMClient)
		if !ok {
			t.Fatalf("expected scram client type")
		}

		if client.HashGeneratorFcn == nil {
			t.Fatalf("expected hash generator to be set")
		}
	})

	t.Run("unsupported mechanism", func(t *testing.T) {
		config := newConfig()

		err := applySASLConfig(config, true, sarama.SASLMechanism("oauthbearer"), "user", "pass")
		if err == nil {
			t.Fatalf("expected error for unsupported mechanism")
		}
	})
}

func TestDecodeKafkaPEM(t *testing.T) {
	t.Run("codec path", func(t *testing.T) {
		coder := codec.NewBase64Coder()
		encoded := coder.EncodeBase64("pem-data")

		decoded := DecodeKafkaPEM(coder, encoded)
		if decoded != "pem-data" {
			t.Fatalf("expected decoded pem-data, got %q", decoded)
		}
	})

	t.Run("codec invalid input returns empty string", func(t *testing.T) {
		coder := codec.NewBase64Coder()

		decoded := DecodeKafkaPEM(coder, "not-base64")
		if decoded != "" {
			t.Fatalf("expected empty string, got %q", decoded)
		}
	})

	t.Run("byte path", func(t *testing.T) {
		encoded := base64.StdEncoding.EncodeToString([]byte("pem-data"))

		decoded := DecodeKafkaPEMBytes(encoded)
		if decoded != "pem-data" {
			t.Fatalf("expected decoded pem-data, got %q", decoded)
		}
	})

	t.Run("byte invalid input returns empty string", func(t *testing.T) {
		decoded := DecodeKafkaPEMBytes("not-base64")
		if decoded != "" {
			t.Fatalf("expected empty string, got %q", decoded)
		}
	})
}

func TestNormalizePEM(t *testing.T) {
	input := "line1\\nline2"
	if got := normalizePEM(input); got != "line1\nline2" {
		t.Fatalf("expected newline normalized string, got %q", got)
	}
}

func TestSCRAMClientBegin(t *testing.T) {
	client := &SCRAMClient{HashGeneratorFcn: SHA256}

	if err := client.Begin("user", "password", ""); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if client.Client == nil {
		t.Fatalf("expected client to be initialized")
	}

	if client.ClientConversation == nil {
		t.Fatalf("expected client conversation to be initialized")
	}

	if client.Done() {
		t.Fatalf("expected conversation not done immediately after begin")
	}
}

func createSelfSignedCertPEM(t *testing.T) string {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate rsa key: %v", err)
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "test-ca",
		},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
	if err != nil {
		t.Fatalf("create certificate: %v", err)
	}

	return string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes}))
}
