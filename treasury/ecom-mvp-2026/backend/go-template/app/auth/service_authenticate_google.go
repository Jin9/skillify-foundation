package auth

import (
	"context"
	"fmt"
	"log/slog"
)

// errUnauthorized signals the caller to respond with 401.
var errUnauthorized = fmt.Errorf("unauthorized")

// googleAuthData holds the validated Google profile information.
type googleAuthData struct {
	email        string
	username     string
	profileImage string
}

// authenticateGoogle validates the access token, fetches the user profile,
// and revokes the token. Returns errUnauthorized for auth failures.
func (h *handler) authenticateGoogle(ctx context.Context, accessToken string) (*googleAuthData, error) {
	validateTokenResp, err := h.googleClient.ValidateAccessToken(ctx, accessToken)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errUnauthorized, err)
	}
	if validateTokenResp.Code >= 400 {
		return nil, errUnauthorized
	}

	userProfileResp, err := h.googleClient.GetUserProfile(ctx, accessToken)
	if err != nil {
		return nil, fmt.Errorf("get user profile: %w", err)
	}
	if userProfileResp.Code >= 400 {
		return nil, fmt.Errorf("get user profile: unexpected status %d", userProfileResp.Code)
	}

	revokeResp, err := h.googleClient.RevokeToken(ctx, accessToken)
	if err != nil {
		slog.Warn("fail to revoke google access token with client error", slog.String("err", err.Error()), slog.String("tag", "resolve identify"))
	}
	if revokeResp.Code >= 400 {
		slog.Warn("fail to revoke google access token with http error", slog.Int("httpCode", revokeResp.Code), slog.String("tag", "resolve identify"))
	}

	return &googleAuthData{
		email:        validateTokenResp.Response.Email,
		username:     userProfileResp.Response.Name,
		profileImage: userProfileResp.Response.Picture,
	}, nil
}
