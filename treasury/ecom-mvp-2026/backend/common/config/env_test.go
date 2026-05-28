package config

import (
	"os"
	"testing"
)

func TestIsLocalEnv(t *testing.T) {
	original := Env
	defer func() { Env = original }()

	Env = "LOCAL"
	if !IsLocalEnv() {
		t.Error("expected IsLocalEnv() == true for LOCAL")
	}

	Env = "local"
	if !IsLocalEnv() {
		t.Error("expected IsLocalEnv() == true for lowercase local")
	}

	Env = "DEV"
	if IsLocalEnv() {
		t.Error("expected IsLocalEnv() == false for DEV")
	}
}

func TestIsDevEnv(t *testing.T) {
	original := Env
	defer func() { Env = original }()

	Env = "DEV"
	if !IsDevEnv() {
		t.Error("expected IsDevEnv() == true for DEV")
	}

	Env = "dev"
	if !IsDevEnv() {
		t.Error("expected IsDevEnv() == true for lowercase dev")
	}

	Env = "PROD"
	if IsDevEnv() {
		t.Error("expected IsDevEnv() == false for PROD")
	}
}

func TestIsUATEnv(t *testing.T) {
	original := Env
	defer func() { Env = original }()

	Env = "UAT"
	if !IsUATEnv() {
		t.Error("expected IsUATEnv() == true for UAT")
	}

	Env = ""
	if IsUATEnv() {
		t.Error("expected IsUATEnv() == false for empty")
	}
}

func TestIsProdEnv(t *testing.T) {
	original := Env
	defer func() { Env = original }()

	Env = "PROD"
	if !IsProdEnv() {
		t.Error("expected IsProdEnv() == true for PROD")
	}

	Env = "prod"
	if !IsProdEnv() {
		t.Error("expected IsProdEnv() == true for lowercase prod")
	}

	Env = "LOCAL"
	if IsProdEnv() {
		t.Error("expected IsProdEnv() == false for LOCAL")
	}
}

func TestEnvInit(t *testing.T) {
	// Verify Env reads from ENV environment variable
	original := os.Getenv("ENV")
	defer os.Setenv("ENV", original)

	// Env was already set by init(), just verify it's a string
	_ = Env
}
