package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/token"
	mocks "gitlab.com/b2c-e-commerce-platform/platform/backend/common/token/mocks"
)

func TestJWT(t *testing.T) {
	gin.SetMode(gin.TestMode)

	errMock := errors.New("mock error")

	tests := []struct {
		name           string
		headerAuth     string
		mockVerifyErr  error
		mockParseErr   error
		mockParseClaim *token.Claims
		mockValErr     error
		expectedStatus int
	}{
		{
			name:           "Missing Header",
			headerAuth:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Verify Signature Fails",
			headerAuth:     "Bearer invalid-token",
			mockVerifyErr:  errMock,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Parse Token Fails",
			headerAuth:     "Bearer valid-signature-bad-payload",
			mockVerifyErr:  nil,
			mockParseErr:   errMock,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Validate Claims Fails",
			headerAuth:     "Bearer valid-parse-bad-claims",
			mockVerifyErr:  nil,
			mockParseErr:   nil,
			mockParseClaim: &token.Claims{},
			mockValErr:     errMock,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Refresh Token Rejected",
			headerAuth:     "Bearer refresh-token",
			mockVerifyErr:  nil,
			mockParseErr:   nil,
			mockParseClaim: &token.Claims{Extra: map[string]any{"tokenType": "refresh"}},
			mockValErr:     nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:          "Missing tokenType Rejected",
			headerAuth:    "Bearer legacy-token",
			mockVerifyErr: nil,
			mockParseErr:  nil,
			mockParseClaim: &token.Claims{
				Extra: map[string]any{},
			},
			mockValErr:     nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:          "Success",
			headerAuth:    "Bearer super-valid-token",
			mockVerifyErr: nil,
			mockParseErr:  nil,
			mockParseClaim: &token.Claims{
				Extra: map[string]any{"tokenType": "access"},
			},
			mockValErr:     nil,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockVerifier := new(mocks.JWTVerifierMock)
			mockParser := new(mocks.JWTParserMock)

			if tt.headerAuth != "" {
				rawToken := extractBearerToken(tt.headerAuth)

				mockVerifier.On("VerifySignatureES256", rawToken).Return(tt.mockVerifyErr)

				if tt.mockVerifyErr == nil {
					mockParser.On("ParseToken", rawToken).Return(tt.mockParseClaim, tt.mockParseErr)

					if tt.mockParseErr == nil {
						mockParser.On("ValidateClaims", tt.mockParseClaim).Return(tt.mockValErr)
					}
				}
			}

			r := gin.New()
			r.Use(JWT(mockParser, mockVerifier))

			r.GET("/protected", func(c *gin.Context) {
				_, err := token.ClaimsFromContext(c.Request.Context())
				if err != nil {
					c.Status(http.StatusInternalServerError)
					return
				}
				c.Status(http.StatusOK)
			})

			req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
			if tt.headerAuth != "" {
				req.Header.Set("Authorization", tt.headerAuth)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockVerifier.AssertExpectations(t)
			mockParser.AssertExpectations(t)
		})
	}
}
