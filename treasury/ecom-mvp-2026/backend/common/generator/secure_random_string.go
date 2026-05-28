package generator

import (
	"crypto/rand"
	"errors"
)

var randReader = rand.Reader

func GenerateSecureString(length int, charset string) (string, error) {
	if length <= 0 {
		return "", errors.New("length must be positive")
	}
	if len(charset) == 0 {
		return "", errors.New("charset cannot be empty")
	}

	b := make([]byte, length)
	charsetLen := byte(len(charset))
	maxByte := 255 - (255 % charsetLen)

	randBuf := make([]byte, length*2) // buffer twice the size to reduce retries
	n := 0

	for i := 0; i < length; {
		if n == 0 {
			if _, err := randReader.Read(randBuf); err != nil {
				return "", err
			}
			n = len(randBuf)
		}

		r := randBuf[len(randBuf)-n]
		n--

		if r > maxByte {
			continue
		}

		b[i] = charset[r%charsetLen]
		i++
	}

	return string(b), nil
}
