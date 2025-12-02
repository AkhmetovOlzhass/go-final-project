package service

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"learning-platform/internal/models"
)



type fakeUserRepo struct {
	byEmail map[string]*models.User
	byID    map[string]*models.User

	findByEmailErr error
	findByIDErr    error
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{
		byEmail: make(map[string]*models.User),
		byID:    make(map[string]*models.User),
	}
}

func (f *fakeUserRepo) Create(ctx context.Context, user *models.User) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	f.byEmail[user.Email] = user
	f.byID[user.ID.String()] = user
	return nil
}

func (f *fakeUserRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	if f.findByEmailErr != nil {
		return nil, f.findByEmailErr
	}
	return f.byEmail[email], nil
}

func (f *fakeUserRepo) FindByID(ctx context.Context, id string) (*models.User, error) {
	if f.findByIDErr != nil {
		return nil, f.findByIDErr
	}
	return f.byID[id], nil
}


func (f *fakeUserRepo) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	u, ok := f.byID[id]
	if !ok || u == nil {
		return nil
	}

	if status, ok := updates["status"].(models.UserStatus); ok {
		u.Status = status
	}
	if role, ok := updates["role"].(models.UserRole); ok {
		u.Role = role
	}
	if name, ok := updates["display_name"].(string); ok {
		u.DisplayName = name
	}

	if email, ok := updates["email"].(string); ok {
		u.Email = email
		f.byEmail[email] = u
	}

	if name2, ok := updates["displayName"].(string); ok {
		u.DisplayName = name2
	}

	if avatar, ok := updates["avatarUrl"].(string); ok {
		u.AvatarURL = &avatar
	}

	return nil
}

func (f *fakeUserRepo) GetAll(ctx context.Context) ([]models.User, error) {
	res := make([]models.User, 0, len(f.byID))
	for _, u := range f.byID {
		if u != nil {
			res = append(res, *u)
		}
	}
	return res, nil
}



type fakeTokenRepo struct {
	byHash map[string]*models.RefreshToken

	saveErr          error
	findErr          error
	revokeErr        error
	revokeAllErr     error
	deleteExpiredErr error

	revokeAllFor []uuid.UUID
}

func newFakeTokenRepo() *fakeTokenRepo {
	return &fakeTokenRepo{
		byHash: make(map[string]*models.RefreshToken),
	}
}

func (f *fakeTokenRepo) Save(ctx context.Context, token *models.RefreshToken) error {
	if f.saveErr != nil {
		return f.saveErr
	}
	f.byHash[token.TokenHash] = token
	return nil
}

func (f *fakeTokenRepo) FindValid(ctx context.Context, hash string) (*models.RefreshToken, error) {
	if f.findErr != nil {
		return nil, f.findErr
	}
	tk, ok := f.byHash[hash]
	if !ok {
		return nil, nil
	}
	return tk, nil
}

func (f *fakeTokenRepo) Revoke(ctx context.Context, hash string) error {
	if f.revokeErr != nil {
		return f.revokeErr
	}
	delete(f.byHash, hash)
	return nil
}

func (f *fakeTokenRepo) RevokeAllForUser(ctx context.Context, userID uuid.UUID) error {
	if f.revokeAllErr != nil {
		return f.revokeAllErr
	}
	f.revokeAllFor = append(f.revokeAllFor, userID)

	for h, t := range f.byHash {
		if t.UserID == userID {
			delete(f.byHash, h)
		}
	}
	return nil
}

func (f *fakeTokenRepo) DeleteExpired(ctx context.Context) error {
	if f.deleteExpiredErr != nil {
		return f.deleteExpiredErr
	}
	return nil
}

type fakeVerificationRepo struct{}

func (f *fakeVerificationRepo) Create(ctx context.Context, v *models.EmailVerification) error {
	return nil
}
func (f *fakeVerificationRepo) FindValid(ctx context.Context, email, code string) (*models.EmailVerification, error) {
	return nil, nil
}
func (f *fakeVerificationRepo) MarkUsed(ctx context.Context, id uuid.UUID) error {
	return nil
}

func TestAuthService_generateHash(t *testing.T) {
	svc := &AuthService{
		secret: "test-secret",
	}

	h1 := svc.generateHash("value")
	h2 := svc.generateHash("value")
	h3 := svc.generateHash("other")

	require.NotEmpty(t, h1)
	require.NotEmpty(t, h2)
	require.NotEmpty(t, h3)

	assert.Equal(t, h1, h2)
	assert.NotEqual(t, h1, h3)
}

func TestAuthService_createJWT(t *testing.T) {
	secret := "super-secret"

	svc := &AuthService{
		secret: secret,
	}

	userID := uuid.New()
	role := models.UserRole("STUDENT")

	tokenStr, err := svc.createJWT(userID, role)
	require.NoError(t, err)
	require.NotEmpty(t, tokenStr)

	parsed, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	require.NoError(t, err)
	require.True(t, parsed.Valid)

	claims, ok := parsed.Claims.(jwt.MapClaims)
	require.True(t, ok)

	assert.Equal(t, userID.String(), claims["userId"])
	assert.Equal(t, string(role), claims["role"])

	exp, ok := claims["exp"].(float64)
	require.True(t, ok)
	assert.Greater(t, int64(exp), time.Now().Unix())
}

func TestAuthService_Login(t *testing.T) {
	ctx := context.Background()

	password := "password123"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)

	activeUserID := uuid.New()

	tests := []struct {
		name          string
		setupUsers    func() *fakeUserRepo
		setupTokens   func() *fakeTokenRepo
		inputEmail    string
		inputPassword string
		wantErr       bool
		errContains   string
	}{
		{
			name: "user not found",
			setupUsers: func() *fakeUserRepo {
				return newFakeUserRepo()
			},
			setupTokens:   func() *fakeTokenRepo { return newFakeTokenRepo() },
			inputEmail:    "unknown@example.com",
			inputPassword: password,
			wantErr:       true,
			errContains:   "invalid credentials",
		},
		{
			name: "wrong password",
			setupUsers: func() *fakeUserRepo {
				repo := newFakeUserRepo()
				u := &models.User{
					ID:           activeUserID,
					Email:        "user@example.com",
					PasswordHash: string(hash),
					Status:       models.UserStatus("ACTIVE"),
					Role:         models.UserRole("STUDENT"),
				}
				repo.byEmail[u.Email] = u
				repo.byID[u.ID.String()] = u
				return repo
			},
			setupTokens:   func() *fakeTokenRepo { return newFakeTokenRepo() },
			inputEmail:    "user@example.com",
			inputPassword: "wrong",
			wantErr:       true,
			errContains:   "invalid credentials",
		},
		{
			name: "email not verified",
			setupUsers: func() *fakeUserRepo {
				repo := newFakeUserRepo()
				u := &models.User{
					ID:           activeUserID,
					Email:        "user@example.com",
					PasswordHash: string(hash),
					Status:       models.UserStatus("PENDING"),
					Role:         models.UserRole("STUDENT"),
				}
				repo.byEmail[u.Email] = u
				repo.byID[u.ID.String()] = u
				return repo
			},
			setupTokens:   func() *fakeTokenRepo { return newFakeTokenRepo() },
			inputEmail:    "user@example.com",
			inputPassword: password,
			wantErr:       true,
			errContains:   "email not verified",
		},
		{
			name: "success login",
			setupUsers: func() *fakeUserRepo {
				repo := newFakeUserRepo()
				u := &models.User{
					ID:           activeUserID,
					Email:        "user@example.com",
					PasswordHash: string(hash),
					Status:       models.UserStatus("ACTIVE"),
					Role:         models.UserRole("STUDENT"),
				}
				repo.byEmail[u.Email] = u
				repo.byID[u.ID.String()] = u
				return repo
			},
			setupTokens:   func() *fakeTokenRepo { return newFakeTokenRepo() },
			inputEmail:    "user@example.com",
			inputPassword: password,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := tt.setupUsers()
			tokenRepo := tt.setupTokens()

			svc := &AuthService{
				users:         userRepo,
				verifications: nil,
				tokens:        tokenRepo,
				secret:        "test-secret",
			}

			access, refresh, err := svc.Login(ctx, tt.inputEmail, tt.inputPassword)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Empty(t, access)
				assert.Empty(t, refresh)
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, access)
			assert.NotEmpty(t, refresh)

			require.Len(t, tokenRepo.byHash, 1)
			require.Len(t, tokenRepo.revokeAllFor, 1)
			assert.Equal(t, activeUserID, tokenRepo.revokeAllFor[0])
		})
	}
}


func TestAuthService_Refresh_InvalidToken(t *testing.T) {
	ctx := context.Background()

	userRepo := newFakeUserRepo()
	tokenRepo := newFakeTokenRepo()

	svc := &AuthService{
		users:         userRepo,
		verifications: nil,
		tokens:        tokenRepo,
		secret:        "test-secret",
	}

	_, _, err := svc.Refresh(ctx, "non-existent")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid refresh token")
}

func TestAuthService_Refresh_Success(t *testing.T) {
	ctx := context.Background()

	userRepo := newFakeUserRepo()
	tokenRepo := newFakeTokenRepo()

	secret := "test-secret"

	svc := &AuthService{
		users:         userRepo,
		verifications: nil,
		tokens:        tokenRepo,
		secret:        secret,
	}

	userID := uuid.New()
	user := &models.User{
		ID:     userID,
		Email:  "user@example.com",
		Status: models.UserStatus("ACTIVE"),
		Role:   models.UserRole("STUDENT"),
	}
	userRepo.byID[userID.String()] = user

	oldRaw := "old-refresh"
	oldHash := svc.generateHash(oldRaw)
	tokenRepo.byHash[oldHash] = &models.RefreshToken{
		UserID:    userID,
		TokenHash: oldHash,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	newAccess, newRefresh, err := svc.Refresh(ctx, oldRaw)
	require.NoError(t, err)

	assert.NotEmpty(t, newAccess)
	assert.NotEmpty(t, newRefresh)

	_, ok := tokenRepo.byHash[oldHash]
	assert.False(t, ok)

	require.Equal(t, 1, len(tokenRepo.byHash))
}
