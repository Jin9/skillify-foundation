package generator

import (
	"crypto/sha512"
	"fmt"
	"time"
)

// GeneratePwd matches your Java logic exactly
func GeneratePwd(key string, batchDate time.Time, length int) string {
	// Format date as yyyy-MM-dd
	dateString := batchDate.Format("2006-01-02")
	actualKey := key + dateString

	// SHA-512 hash
	hash := sha512.Sum512([]byte(actualKey))

	// Convert hash to hex string (same as your Java toHexString)
	hexString := toHexString(hash[:])

	// Return last `length` characters, or full string if shorter
	if len(hexString) > length {
		return hexString[len(hexString)-length:]
	}
	return hexString
}

// toHexString mimics your Java logic (not using encoding/hex)
func toHexString(bytes []byte) string {
	hexString := ""
	for _, b := range bytes {
		hex := fmt.Sprintf("%x", b)
		if len(hex) == 1 {
			hexString += "0"
		}
		hexString += hex
	}
	return hexString
}
