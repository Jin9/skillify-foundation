package identity

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	validAccessExp := time.Date(2026, 1, 1, 0, 15, 0, 0, time.UTC)
	validRefreshExp := time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		body           string
		loginTokens    loginTokenResult
		svcErr         error
		wantHTTPStatus int
		wantCode       string
	}{
		{
			name: "ok — valid credentials return tokens",
			body: `{"email":"alice@example.com","password":"correct-pw"}`,
			loginTokens: loginTokenResult{
				accessToken:  "access.tok.en",
				refreshToken: "refresh.tok.en",
				accessExp:    validAccessExp,
				refreshExp:   validRefreshExp,
				role:         "CUSTOMER",
			},
			wantHTTPStatus: http.StatusOK,
			wantCode:       "0000",
		},
		{
			name:           "400 on missing email field",
			body:           `{"password":"pw"}`,
			wantHTTPStatus: http.StatusBadRequest,
			wantCode:       "HM400",
		},
		{
			name:           "400 on missing password field",
			body:           `{"email":"alice@example.com"}`,
			wantHTTPStatus: http.StatusBadRequest,
			wantCode:       "HM400",
		},
		{
			name:           "400 on malformed JSON",
			body:           `{not-json}`,
			wantHTTPStatus: http.StatusBadRequest,
			wantCode:       "HM400",
		},
		{
			name:           "401 on invalid credentials (AUTH_INVALID)",
			body:           `{"email":"alice@example.com","password":"wrong-pw"}`,
			svcErr:         ErrAuthInvalid,
			wantHTTPStatus: http.StatusUnauthorized,
			wantCode:       "AUTH_INVALID",
		},
		{
			name:           "403 on suspended account (AUTH_SUSPENDED)",
			body:           `{"email":"suspended@example.com","password":"correct-pw"}`,
			svcErr:         ErrAuthSuspended,
			wantHTTPStatus: http.StatusForbidden,
			wantCode:       "AUTH_SUSPENDED",
		},
		{
			name:           "500 on unexpected service error",
			body:           `{"email":"alice@example.com","password":"correct-pw"}`,
			svcErr:         ErrUserNotFound, // non-auth sentinel → falls through to 500
			wantHTTPStatus: http.StatusInternalServerError,
			wantCode:       "HM500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &fakeService{
				err:         tt.svcErr,
				loginTokens: tt.loginTokens,
			}
			h := NewHandler(HandlerConfig{Service: svc})

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(
				http.MethodPost,
				"/api/v1/identity/auth/login",
				bytes.NewBufferString(tt.body),
			)
			c.Request.Header.Set("Content-Type", "application/json")

			h.Login(c)

			if w.Code != tt.wantHTTPStatus {
				t.Fatalf("http status: want %d got %d (body=%s)", tt.wantHTTPStatus, w.Code, w.Body.String())
			}

			var env struct {
				Code string `json:"code"`
				Data *struct {
					AccessToken  string `json:"accessToken"`
					RefreshToken string `json:"refreshToken"`
					Role         string `json:"role"`
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
					t.Fatal("expected data payload on 200 response")
				}
				if env.Data.AccessToken != tt.loginTokens.accessToken {
					t.Errorf("accessToken: want %q got %q", tt.loginTokens.accessToken, env.Data.AccessToken)
				}
				if env.Data.RefreshToken != tt.loginTokens.refreshToken {
					t.Errorf("refreshToken: want %q got %q", tt.loginTokens.refreshToken, env.Data.RefreshToken)
				}
				if env.Data.Role != tt.loginTokens.role {
					t.Errorf("role: want %q got %q", tt.loginTokens.role, env.Data.Role)
				}
			}
		})
	}
}
