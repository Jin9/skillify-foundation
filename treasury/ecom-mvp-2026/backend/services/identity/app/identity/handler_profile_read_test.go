package identity

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/token"
)

// loginTokenResult holds the values returned by fakeService.Login.
type loginTokenResult struct {
	accessToken  string
	refreshToken string
	accessExp    time.Time
	refreshExp   time.Time
	role         string
}

// refreshTokenResult holds the values returned by fakeService.Refresh.
type refreshTokenResult struct {
	accessToken  string
	refreshToken string
	accessExp    time.Time
	refreshExp   time.Time
}

// fakeService is a minimal in-package stub that implements the Service
// interface for unit tests. Only the fields relevant to a given handler test
// need to be populated; all other methods return zero values.
type fakeService struct {
	// Shared error returned by methods that only signal success/failure.
	err error

	// GetUserByID result.
	user User

	// Login result (ignored when err != nil).
	loginTokens loginTokenResult

	// Refresh result (ignored when err != nil).
	refreshTokens refreshTokenResult
}

func (f *fakeService) Register(context.Context, string, string, string, *string) (User, error) {
	return User{}, nil
}

func (f *fakeService) Login(_ context.Context, _, _ string) (string, string, time.Time, time.Time, string, error) {
	if f.err != nil {
		return "", "", time.Time{}, time.Time{}, "", f.err
	}
	lt := f.loginTokens
	return lt.accessToken, lt.refreshToken, lt.accessExp, lt.refreshExp, lt.role, nil
}

func (f *fakeService) Refresh(_ context.Context, _ string) (string, string, time.Time, time.Time, error) {
	if f.err != nil {
		return "", "", time.Time{}, time.Time{}, f.err
	}
	rt := f.refreshTokens
	return rt.accessToken, rt.refreshToken, rt.accessExp, rt.refreshExp, nil
}

func (f *fakeService) GetUserByID(_ context.Context, _ uuid.UUID) (User, error) {
	return f.user, f.err
}

func (f *fakeService) SetAddressDefault(_ context.Context, _, _ uuid.UUID) error {
	return f.err
}

func phonePtr(s string) *string { return &s }

func TestProfileRead(t *testing.T) {
	gin.SetMode(gin.TestMode)

	validUserID := uuid.New()
	validUser := User{
		ID:    validUserID,
		Email: "alice@example.com",
		Name:  "Alice Tester",
		Phone: phonePtr("0801234567"),
		Role:  "CUSTOMER",
	}

	tests := []struct {
		name           string
		claimsSub      string
		injectClaims   bool
		svcUser        User
		svcErr         error
		wantHTTPStatus int
		wantCode       string
	}{
		{
			name:           "ok — returns profile for authenticated customer",
			claimsSub:      validUserID.String(),
			injectClaims:   true,
			svcUser:        validUser,
			wantHTTPStatus: http.StatusOK,
			wantCode:       "0000",
		},
		{
			name:           "401 when no claims in context",
			injectClaims:   false,
			wantHTTPStatus: http.StatusUnauthorized,
			wantCode:       "AUTH_INVALID",
		},
		{
			name:           "401 when claims.sub is not a UUID",
			claimsSub:      "not-a-uuid",
			injectClaims:   true,
			wantHTTPStatus: http.StatusUnauthorized,
			wantCode:       "AUTH_INVALID",
		},
		{
			name:           "404 when user no longer exists",
			claimsSub:      validUserID.String(),
			injectClaims:   true,
			svcErr:         ErrUserNotFound,
			wantHTTPStatus: http.StatusNotFound,
			wantCode:       "USER_NOT_FOUND",
		},
		{
			name:           "500 on unexpected storage error",
			claimsSub:      validUserID.String(),
			injectClaims:   true,
			svcErr:         errors.New("connection refused"),
			wantHTTPStatus: http.StatusInternalServerError,
			wantCode:       "HM500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &fakeService{user: tt.svcUser, err: tt.svcErr}
			h := NewHandler(HandlerConfig{Service: svc})

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/identity/profile/read", bytes.NewReader(nil))

			if tt.injectClaims {
				ctx := token.WithClaims(req.Context(), &token.Claims{Sub: tt.claimsSub})
				req = req.WithContext(ctx)
			}
			c.Request = req

			h.ProfileRead(c)

			if w.Code != tt.wantHTTPStatus {
				t.Fatalf("http status: want %d got %d (body=%s)", tt.wantHTTPStatus, w.Code, w.Body.String())
			}

			var env struct {
				Code    string `json:"code"`
				TraceID string `json:"traceId"`
				Data    *struct {
					UserID string  `json:"userId"`
					Email  string  `json:"email"`
					Name   string  `json:"name"`
					Phone  *string `json:"phone"`
					Role   string  `json:"role"`
				} `json:"data"`
			}
			if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
				t.Fatalf("unmarshal: %v (body=%s)", err, w.Body.String())
			}
			if env.Code != tt.wantCode {
				t.Errorf("code: want %q got %q", tt.wantCode, env.Code)
			}
			if tt.wantHTTPStatus == http.StatusOK {
				if env.Data == nil {
					t.Fatal("expected data on success")
				}
				if env.Data.UserID != validUser.ID.String() {
					t.Errorf("userId: want %s got %s", validUser.ID, env.Data.UserID)
				}
				if env.Data.Email != validUser.Email {
					t.Errorf("email: want %s got %s", validUser.Email, env.Data.Email)
				}
				if env.Data.Role != validUser.Role {
					t.Errorf("role: want %s got %s", validUser.Role, env.Data.Role)
				}
			}
		})
	}
}
