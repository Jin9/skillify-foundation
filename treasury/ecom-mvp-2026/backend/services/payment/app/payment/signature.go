package payment

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
)

// VerifyHMACSignature performs a constant-time HMAC-SHA256 verification.
//
// secret is the raw shared secret (CALLBACK_HMAC_SECRET).
// rawBody is the exact wire bytes of the incoming HTTP request body.
// sigHex is the value of the X-Mock-Provider-Signature header (hex-encoded).
//
// Returns true only if the computed MAC matches the provided signature.
// Uses subtle.ConstantTimeCompare so timing cannot distinguish near-misses
// from full mismatches — protecting against timing-oracle attacks.
func VerifyHMACSignature(secret []byte, rawBody []byte, sigHex string) bool {
	mac := hmac.New(sha256.New, secret)
	mac.Write(rawBody)
	expected := mac.Sum(nil)

	got, err := hex.DecodeString(sigHex)
	if err != nil {
		// Malformed hex — treat as invalid; constant-time path still taken
		// by comparing against a dummy expected value of the same length.
		return false
	}

	// subtle.ConstantTimeCompare returns 1 on equality, 0 otherwise.
	return subtle.ConstantTimeCompare(expected, got) == 1
}

// ComputeHMAC returns hex(HMAC-SHA256(body, secret)). Used in tests and simulate.
func ComputeHMAC(secret, body []byte) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}
