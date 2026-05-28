package token

import (
	"context"
	"errors"
	"testing"
)

func TestClaimsFromContext(t *testing.T) {
	t.Run("returns claim from context", func(t *testing.T) {
		expected := &Claims{Sub: "member-1"}
		ctx := WithClaims(context.Background(), expected)

		actual, err := ClaimsFromContext(ctx)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if actual != expected {
			t.Fatalf("expected same claim pointer")
		}
	})

	t.Run("returns not found when no claim in context", func(t *testing.T) {
		ctx := context.TODO()

		actual, err := ClaimsFromContext(ctx)
		if actual != nil {
			t.Fatalf("expected nil claim")
		}

		if !errors.Is(err, ErrClaimNotFound) {
			t.Fatalf("expected ErrClaimNotFound, got %v", err)
		}
	})

	t.Run("returns not found on empty context", func(t *testing.T) {
		actual, err := ClaimsFromContext(context.TODO())
		if actual != nil {
			t.Fatalf("expected nil claim")
		}

		if !errors.Is(err, ErrClaimNotFound) {
			t.Fatalf("expected ErrClaimNotFound, got %v", err)
		}
	})

	t.Run("returns not found when missing claim", func(t *testing.T) {
		actual, err := ClaimsFromContext(context.Background())
		if actual != nil {
			t.Fatalf("expected nil claim")
		}

		if !errors.Is(err, ErrClaimNotFound) {
			t.Fatalf("expected ErrClaimNotFound, got %v", err)
		}
	})

	t.Run("returns invalid claim for wrong type", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), claimsContextKey, "wrong")

		actual, err := ClaimsFromContext(ctx)
		if actual != nil {
			t.Fatalf("expected nil claim")
		}

		if !errors.Is(err, ErrInvalidClaim) {
			t.Fatalf("expected ErrInvalidClaim, got %v", err)
		}
	})

	t.Run("returns not found for typed nil claim", func(t *testing.T) {
		var empty *Claims
		ctx := WithClaims(context.Background(), empty)

		actual, err := ClaimsFromContext(ctx)
		if actual != nil {
			t.Fatalf("expected nil claim")
		}

		if !errors.Is(err, ErrClaimNotFound) {
			t.Fatalf("expected ErrClaimNotFound, got %v", err)
		}
	})
}
