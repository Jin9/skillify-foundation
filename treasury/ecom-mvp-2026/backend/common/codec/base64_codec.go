package codec

import (
	"encoding/base64"
	"fmt"
)

type Base64Coder interface {
	EncodeBase64(string) string
	DecodeBase64(string) (string, error)
}

type base64Coder struct{}

func NewBase64Coder() Base64Coder {
	return base64Coder{}
}

func (base64Coder) EncodeBase64(raw string) string {
	return base64.StdEncoding.EncodeToString([]byte(raw))
}

func (base64Coder) DecodeBase64(base64String string) (string, error) {
	if base64String == "" {
		return "", fmt.Errorf("empty base64 string")
	}
	data, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}
	return string(data), nil
}
