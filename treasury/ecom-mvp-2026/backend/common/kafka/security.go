package kafka

import (
	"crypto/sha256"
	"crypto/sha512"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"log/slog"
	"strings"

	"github.com/IBM/sarama"
	"github.com/xdg/scram"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/codec"
)

// PEM decoding
func DecodeKafkaPEM(codec codec.Base64Coder, base64String string) string {
	val, err := decodeKafkaPEM(codec, base64String)
	if err != nil {
		slog.Error("failed to decode kafka pem", "error", err)
		return ""
	}
	return val
}

func DecodeKafkaPEMBytes(base64String string) string {
	data, err := decodeKafkaPEMBytes(base64String)
	if err != nil {
		slog.Error("failed to decode base64", "error", err)
		return ""
	}
	return data
}

func decodeKafkaPEM(codec codec.Base64Coder, base64String string) (string, error) {
	val, err := codec.DecodeBase64(base64String)
	if err != nil {
		return "", fmt.Errorf("decode kafka pem: %w", err)
	}

	return val, nil
}

func decodeKafkaPEMBytes(base64String string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return "", fmt.Errorf("decode kafka pem bytes: %w", err)
	}

	return string(data), nil
}

// TLS
func newTLSConfig(caCertPEM string, clientCertPEM string, clientKeyPEM string) (*tls.Config, error) {
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM([]byte(normalizePEM(caCertPEM))) {
		return nil, fmt.Errorf("error adding CA cert to pool; please check that the \"CA_CERT_PEM\" value is correct")
	}

	var cert tls.Certificate
	if clientCertPEM != "" && clientKeyPEM != "" {
		var err error
		cert, err = tls.X509KeyPair([]byte(normalizePEM(clientCertPEM)), []byte(normalizePEM(clientKeyPEM)))
		if err != nil {
			return nil, fmt.Errorf("error parsing client cert/key pair: %w", err)
		}
	}

	tlsConfig := &tls.Config{
		RootCAs:            caCertPool,
		InsecureSkipVerify: false,
	}

	if cert.Certificate != nil {
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	return tlsConfig, nil
}

func normalizePEM(str string) string {
	return strings.ReplaceAll(str, "\\n", "\n")
}

func applyTLSConfig(conf *sarama.Config, enable bool, caCertPEM, clientCertPEM, clientKeyPEM string) error {
	if conf == nil {
		return fmt.Errorf("kafka config must not be nil")
	}
	if !enable {
		return nil
	}
	if strings.TrimSpace(caCertPEM) == "" {
		return fmt.Errorf("environment variable \"CA_CERT_PEM\" must not be empty when SSL_ENABLE=true")
	}
	tlsConfig, err := newTLSConfig(caCertPEM, clientCertPEM, clientKeyPEM)
	if err != nil {
		return err
	}
	conf.Net.TLS.Enable = true
	conf.Net.TLS.Config = tlsConfig
	return nil
}

// SASL / SCRAM
var (
	SHA256 scram.HashGeneratorFcn = sha256.New
	SHA512 scram.HashGeneratorFcn = sha512.New
)

type SCRAMClient struct {
	*scram.Client
	*scram.ClientConversation
	scram.HashGeneratorFcn
}

func (x *SCRAMClient) Begin(userName, password, authzID string) (err error) {
	x.Client, err = x.HashGeneratorFcn.NewClient(userName, password, authzID)
	if err != nil {
		return err
	}
	x.ClientConversation = x.Client.NewConversation()
	return nil
}

func (x *SCRAMClient) Step(challenge string) (response string, err error) {
	response, err = x.ClientConversation.Step(challenge)
	return
}

func (x *SCRAMClient) Done() bool {
	return x.ClientConversation.Done()
}

func applySASLConfig(conf *sarama.Config, enable bool, mechanism sarama.SASLMechanism, username, password string) error {
	if conf == nil {
		return fmt.Errorf("kafka config must not be nil")
	}
	if !enable {
		return nil
	}
	conf.Net.SASL.Enable = true
	switch mechanism {
	case sarama.SASLTypePlaintext:
		conf.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		conf.Net.SASL.User = username
		conf.Net.SASL.Password = password
	case sarama.SASLTypeSCRAMSHA512:
		conf.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
		conf.Net.SASL.User = username
		conf.Net.SASL.Password = password
		conf.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
			return &SCRAMClient{HashGeneratorFcn: SHA512}
		}
	case sarama.SASLTypeSCRAMSHA256:
		conf.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA256
		conf.Net.SASL.User = username
		conf.Net.SASL.Password = password
		conf.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
			return &SCRAMClient{HashGeneratorFcn: SHA256}
		}
	default:
		return fmt.Errorf("unsupported sasl mechanism: %q", mechanism)
	}

	return nil
}
