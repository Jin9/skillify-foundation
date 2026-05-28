package generator

import "github.com/google/uuid"

type AppUUID interface {
	GenerateUUID() string
	GenerateUUIDV7() uuid.UUID
	GenerateUUIDTypeStrict() uuid.UUID
}

type appUUID struct{}

func NewAppUUID() AppUUID {
	return &appUUID{}
}

func (a appUUID) GenerateUUID() string {
	return uuid.NewString()
}

func (a appUUID) GenerateUUIDTypeStrict() uuid.UUID {
	return uuid.New()
}

func (a appUUID) GenerateUUIDV7() uuid.UUID {
	return uuid.Must(uuid.NewV7())
}
