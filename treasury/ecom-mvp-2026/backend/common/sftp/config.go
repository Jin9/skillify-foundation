package sftp

import "time"

// Config defines the connection details to an SFTP server.
type Config struct {
	Host           string
	Port           string
	Username       string
	Password       string        // Optional if PrivateKeyAuth is provided
	PrivateKeyAuth []byte        // Optional (PEM data)
	Timeout        time.Duration // Default to 10s if zero
}
