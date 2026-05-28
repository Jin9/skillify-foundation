package safe

import "testing"

func TestDeref_NonNil(t *testing.T) {
	s := "hello"
	if got := Deref(&s); got != "hello" {
		t.Errorf("Deref(&%q) = %q, want %q", s, got, "hello")
	}
}

func TestDeref_Nil(t *testing.T) {
	var p *string
	if got := Deref(p); got != "" {
		t.Errorf("Deref(nil) = %q, want empty string", got)
	}
}

func TestDeref_NilInt(t *testing.T) {
	var p *int
	if got := Deref(p); got != 0 {
		t.Errorf("Deref(nil int) = %d, want 0", got)
	}
}

func TestDerefOr_NonNil(t *testing.T) {
	n := 42
	if got := DerefOr(&n, 99); got != 42 {
		t.Errorf("DerefOr(&42, 99) = %d, want 42", got)
	}
}

func TestDerefOr_Nil(t *testing.T) {
	var p *int
	if got := DerefOr(p, 99); got != 99 {
		t.Errorf("DerefOr(nil, 99) = %d, want 99", got)
	}
}

func TestPtr(t *testing.T) {
	p := Ptr("hello")
	if p == nil {
		t.Fatal("Ptr() returned nil")
	}
	if *p != "hello" {
		t.Errorf("*Ptr(%q) = %q", "hello", *p)
	}
}

func TestPtr_Int(t *testing.T) {
	p := Ptr(42)
	if p == nil {
		t.Fatal("Ptr() returned nil")
	}
	if *p != 42 {
		t.Errorf("*Ptr(%d) = %d", 42, *p)
	}
}

func TestIsNil_True(t *testing.T) {
	var p *string
	if !IsNil(p) {
		t.Error("IsNil(nil) should be true")
	}
}

func TestIsNil_False(t *testing.T) {
	s := "hello"
	if IsNil(&s) {
		t.Error("IsNil(&s) should be false")
	}
}
