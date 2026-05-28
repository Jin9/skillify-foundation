package identity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	Name         string
	Phone        *string
	Role         string
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Address struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	ReceiverName string
	Phone        string
	AddressLine1 string
	AddressLine2 *string
	Province     string
	District     string
	SubDistrict  *string
	PostalCode   string
	IsDefault    bool
	DeletedAt    *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type RefreshToken struct {
	JTI           uuid.UUID
	UserID        uuid.UUID
	IssuedAt      time.Time
	ExpiresAt     time.Time
	RevokedAt     *time.Time
	ReplacedByJTI *uuid.UUID
}
