package identity

import (
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/token"
)

// Domain constants — avoids stringly-typed literals scattered across the package.
const (
	roleCustomer = "CUSTOMER"
	statusActive = "ACTIVE"

	tokenTypeAccess  = "access"
	tokenTypeRefresh = "refresh"
)

// Storage interfaces live here to avoid import cycles:
// access package imports identity (for domain types), identity package imports access (for interfaces).
// Defining interfaces in the identity package breaks the cycle.

type UserStorage interface {
	Insert(ctx context.Context, u User) error
	GetByEmail(ctx context.Context, email string) (User, error)
	GetByID(ctx context.Context, id uuid.UUID) (User, error)
}

type RefreshTokenStorage interface {
	Insert(ctx context.Context, rt RefreshToken) error
	RotateInTx(ctx context.Context, oldJTI, newJTI uuid.UUID, newRT RefreshToken) (int64, error)
	Revoke(ctx context.Context, jti uuid.UUID) error
}

// AddressStorage covers address persistence operations.
// SetAddressDefault is the only method needed for the set-default endpoint;
// all other address operations remain stubbed until post-MVP.
type AddressStorage interface {
	// SetDefault atomically clears is_default on all other live addresses for
	// the user, then sets is_default = true on the target. Returns the number
	// of rows affected by the second UPDATE (the target row). Returns 0 when
	// the address is not found, soft-deleted, or belongs to a different user.
	SetDefault(ctx context.Context, userID, addressID uuid.UUID) (int64, error)
}

// Service defines the operations used by handlers.
type Service interface {
	Register(ctx context.Context, email, password, name string, phone *string) (User, error)
	Login(ctx context.Context, email, password string) (accessToken, refreshToken string, accessExp, refreshExp time.Time, role string, err error)
	Refresh(ctx context.Context, rawRefreshToken string) (accessToken, refreshToken string, accessExp, refreshExp time.Time, err error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (User, error)
	SetAddressDefault(ctx context.Context, userID, addressID uuid.UUID) error
}

type ServiceConfig struct {
	UserStorage         UserStorage
	RefreshTokenStorage RefreshTokenStorage
	AddressStorage      AddressStorage
	AccessSigner        token.JWTSigner // signs access tokens
	RefreshSigner       token.JWTSigner // signs refresh tokens
	Verifier            token.JWTVerifier
	Parser              token.JWTParser
	AccessTokenTTL      time.Duration
	RefreshTokenTTL     time.Duration
}

type service struct {
	users         UserStorage
	rts           RefreshTokenStorage
	addrs         AddressStorage
	accessSigner  token.JWTSigner
	refreshSigner token.JWTSigner
	verifier      token.JWTVerifier
	parser        token.JWTParser
	accessTTL     time.Duration
	rtTTL         time.Duration
}

var _ Service = (*service)(nil)

func NewService(cfg ServiceConfig) Service {
	return &service{
		users:         cfg.UserStorage,
		rts:           cfg.RefreshTokenStorage,
		addrs:         cfg.AddressStorage,
		accessSigner:  cfg.AccessSigner,
		refreshSigner: cfg.RefreshSigner,
		verifier:      cfg.Verifier,
		parser:        cfg.Parser,
		accessTTL:     cfg.AccessTokenTTL,
		rtTTL:         cfg.RefreshTokenTTL,
	}
}

func (s *service) Register(ctx context.Context, email, password, name string, phone *string) (User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return User{}, err
	}

	u := User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(hash),
		Name:         name,
		Phone:        phone,
		Role:         roleCustomer,
		Status:       statusActive,
	}

	if err := s.users.Insert(ctx, u); err != nil {
		return User{}, err
	}
	return u, nil
}

func (s *service) Login(ctx context.Context, email, password string) (string, string, time.Time, time.Time, string, error) {
	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		// Run a dummy compare to preserve constant-time behavior (AUTH-007):
		// prevents timing oracle leaking whether the email exists.
		_ = bcrypt.CompareHashAndPassword([]byte("$2a$12$dummyhashfortimingnulltarget...."), []byte(password))
		return "", "", time.Time{}, time.Time{}, "", ErrAuthInvalid
	}

	// bcrypt.CompareHashAndPassword is constant-time-ish: it re-derives the
	// candidate hash from the stored salt+cost and does a fixed-time comparison.
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", time.Time{}, time.Time{}, "", ErrAuthInvalid
	}

	if user.Status == "SUSPENDED" {
		return "", "", time.Time{}, time.Time{}, "", ErrAuthSuspended
	}

	return s.issueTokenPair(ctx, user.ID, user.Role)
}

func (s *service) Refresh(ctx context.Context, rawRefreshToken string) (string, string, time.Time, time.Time, error) {
	if err := s.verifier.VerifySignatureES256(rawRefreshToken); err != nil {
		return "", "", time.Time{}, time.Time{}, ErrAuthInvalid
	}

	claims, err := s.parser.ParseToken(rawRefreshToken)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, ErrAuthInvalid
	}
	if err := s.parser.ValidateClaims(claims); err != nil {
		return "", "", time.Time{}, time.Time{}, ErrAuthInvalid
	}

	oldJTI, err := uuid.Parse(claims.Jti)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, ErrAuthInvalid
	}
	userID, err := uuid.Parse(claims.Sub)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, ErrAuthInvalid
	}

	newJTI := uuid.New()
	now := time.Now().UTC()
	newRT := RefreshToken{
		JTI:       newJTI,
		UserID:    userID,
		IssuedAt:  now,
		ExpiresAt: now.Add(s.rtTTL),
	}

	// Extract role from extra claims.
	role := ""
	if claims.Extra != nil {
		if r, ok := claims.Extra["role"].(string); ok {
			role = r
		}
	}

	rows, err := s.rts.RotateInTx(ctx, oldJTI, newJTI, newRT)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, err
	}
	if rows == 0 {
		return "", "", time.Time{}, time.Time{}, ErrAuthRevoked
	}

	accessToken, accessExp, err := s.signAccess(userID.String(), newJTI.String(), role)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, err
	}

	refreshToken, err := s.signRefresh(userID.String(), newJTI.String(), role)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, err
	}

	return accessToken, refreshToken, accessExp, newRT.ExpiresAt, nil
}

func (s *service) GetUserByID(ctx context.Context, userID uuid.UUID) (User, error) {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		if isUserNotFound(err) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}
	return user, nil
}

// isUserNotFound matches the storage layer's not-found sentinel without
// creating an import cycle from identity → access. The storage error message
// is stable; if storage adds new not-found codes they should be added here.
func isUserNotFound(err error) bool {
	return err != nil && err.Error() == "user not found"
}

func (s *service) issueTokenPair(ctx context.Context, userID uuid.UUID, role string) (string, string, time.Time, time.Time, string, error) {
	jti := uuid.New()
	now := time.Now().UTC()

	rt := RefreshToken{
		JTI:       jti,
		UserID:    userID,
		IssuedAt:  now,
		ExpiresAt: now.Add(s.rtTTL),
	}

	if err := s.rts.Insert(ctx, rt); err != nil {
		return "", "", time.Time{}, time.Time{}, "", err
	}

	accessToken, accessExp, err := s.signAccess(userID.String(), jti.String(), role)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, "", err
	}

	refreshToken, err := s.signRefresh(userID.String(), jti.String(), role)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, "", err
	}

	return accessToken, refreshToken, accessExp, rt.ExpiresAt, role, nil
}

func (s *service) signAccess(sub, jti, role string) (string, time.Time, error) {
	claims := token.Claims{
		Sub: sub,
		Jti: jti,
		Extra: map[string]any{
			"role":      role,
			"tokenType": tokenTypeAccess,
		},
	}
	// Capture now before signing so our exp matches the signer's time.Now() call
	// to within a single clock read. The signer stores exp as Unix seconds, so
	// Truncate mirrors that precision without a round-trip parse of the signed token.
	issuedAt := time.Now().UTC()
	tok, err := s.accessSigner.SignES256(claims)
	if err != nil {
		return "", time.Time{}, err
	}
	exp := issuedAt.Add(s.accessTTL).Truncate(time.Second)
	return tok, exp, nil
}

func (s *service) signRefresh(sub, jti, role string) (string, error) {
	claims := token.Claims{
		Sub: sub,
		Jti: jti,
		Extra: map[string]any{
			"role":      role,
			"tokenType": tokenTypeRefresh,
		},
	}
	tok, err := s.refreshSigner.SignES256(claims)
	return tok, err
}

// SetAddressDefault delegates the two-step atomic UPDATE to the storage layer.
// Returns ErrAddressNotFound when the address does not exist, is soft-deleted,
// or belongs to a different user (storage returns 0 rows affected).
func (s *service) SetAddressDefault(ctx context.Context, userID, addressID uuid.UUID) error {
	rows, err := s.addrs.SetDefault(ctx, userID, addressID)
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrAddressNotFound
	}
	return nil
}
