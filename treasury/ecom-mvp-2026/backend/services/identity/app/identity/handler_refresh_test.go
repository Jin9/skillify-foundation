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

func TestRefresh(t *testing.T) {
	gin.SetMode(gin.TestMode)

	validAccessExp := time.Date(2026, 1, 1, 0, 15, 0, 0, time.UTC)
	validRefreshExp := time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		body           string
		refreshTokens  refreshTokenResult
		svcErr         error
		wantHTTPStatus int
		wantCode       string
	}{
		{
			name: "ok — valid refresh token returns new token pair",
			body: `{"refreshToken":"valid.refresh.token"}`,
			refreshTokens: refreshTokenResult{
				accessToken:  "new.access.tok",
				refreshToken: "new.refresh.tok",
				accessExp:    validAccessExp,
				refreshExp:   validRefreshExp,
			},
			wantHTTPStatus: http.StatusOK,
			wantCode:       "0000",
		},
		{
			name:           "400 on missing refreshToken field",
			body:           `{}`,
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
			name:           "401 on invalid or expired token (AUTH_INVALID)",
			body:           `{"refreshToken":"expired.or.bad.token"}`,
			svcErr:         ErrAuthInvalid,
			wantHTTPStatus: http.StatusUnauthorized,
			wantCode:       "AUTH_INVALID",
		},
		{
			name:           "401 on revoked token (AUTH_REVOKED)",
			body:           `{"refreshToken":"revoked.token"}`,
			svcErr:         ErrAuthRevoked,
			wantHTTPStatus: http.StatusUnauthorized,
			wantCode:       "AUTH_REVOKED",
		},
		{
			name:           "500 on unexpected service error",
			body:           `{"refreshToken":"valid.refresh.token"}`,
			svcErr:         ErrUserNotFound, // non-auth sentinel → falls through to 500
			wantHTTPStatus: http.StatusInternalServerError,
			wantCode:       "HM500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &fakeService{
				err:           tt.svcErr,
				refreshTokens: tt.refreshTokens,
			}
			h := NewHandler(HandlerConfig{Service: svc})

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(
				http.MethodPost,
				"/api/v1/identity/auth/refresh",
				bytes.NewBufferString(tt.body),
			)
			c.Request.Header.Set("Content-Type", "application/json")

			h.Refresh(c)

			if w.Code != tt.wantHTTPStatus {
				t.Fatalf("http status: want %d got %d (body=%s)", tt.wantHTTPStatus, w.Code, w.Body.String())
			}

			var env struct {
				Code string `json:"code"`
				Data *struct {
					AccessToken  string `json:"accessToken"`
					RefreshToken string `json:"refreshToken"`
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
				if env.Data.AccessToken != tt.refreshTokens.accessToken {
					t.Errorf("accessToken: want %q got %q", tt.refreshTokens.accessToken, env.Data.AccessToken)
				}
				if env.Data.RefreshToken != tt.refreshTokens.refreshToken {
					t.Errorf("refreshToken: want %q got %q", tt.refreshTokens.refreshToken, env.Data.RefreshToken)
				}
			}
		})
	}
}
