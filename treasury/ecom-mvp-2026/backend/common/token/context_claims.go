package token

import (
	"context"
	"errors"
)

type contextKey string

const claimsContextKey contextKey = "auth_claims"

var (
	ErrClaimNotFound = errors.New("claim not found")
	ErrInvalidClaim  = errors.New("invalid claim")
)

func WithClaims(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, claimsContextKey, claims)
}

func ClaimsFromContext(ctx context.Context) (*Claims, error) {
	if ctx == nil {
		return nil, ErrClaimNotFound
	}

	claim := ctx.Value(claimsContextKey)
	if claim == nil {
		return nil, ErrClaimNotFound
	}

	parsedClaim, ok := claim.(*Claims)
	if !ok {
		return nil, ErrInvalidClaim
	}

	if parsedClaim == nil {
		return nil, ErrClaimNotFound
	}

	return parsedClaim, nil
}
