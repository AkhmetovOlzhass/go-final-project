package service

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"learning-platform/internal/models"
	"learning-platform/internal/repository"
)

type userRepoMockForAuthService struct {
	mock.Mock
}

func (m *userRepoMockForAuthService) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *userRepoMockForAuthService) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if u, ok := args.Get(0).(*models.User); ok {
		return u, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *userRepoMockForAuthService) FindByID(id string) (*models.User, error) {
	args := m.Called(id)
	if u, ok := args.Get(0).(*models.User); ok {
		return u, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *userRepoMockForAuthService) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}


func (m *userRepoMockForAuthService) GetAll() ([]models.User, error) {
	args := m.Called()
	if users, ok := args.Get(0).([]models.User); ok {
		return users, args.Error(1)
	}
	return nil, args.Error(1)
}

func newTestTokenRepo(t *testing.T) *repository.TokenRepository {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite in memory: %v", err)
	}

	createTableSQL := `
        CREATE TABLE refresh_tokens (
            id TEXT PRIMARY KEY,
            user_id TEXT NOT NULL,
            token_hash TEXT NOT NULL,
            revoked BOOLEAN DEFAULT 0,
            expires_at DATETIME NOT NULL,
            created_at DATETIME
        );
    `
	if err := db.Exec(createTableSQL).Error; err != nil {
		t.Fatalf("failed to create refresh_tokens table: %v", err)
	}

	return repository.NewTokenRepository(db)
}


func TestAuthService_Register_Success(t *testing.T) {
	users := new(userRepoMockForAuthService)
	tokens := newTestTokenRepo(t)

	svc := NewAuthService(users, tokens, "test-secret")

	email := "test@example.com"
	password := "super-secret"
	displayName := "Tester"

	users.
		On("Create", mock.MatchedBy(func(u *models.User) bool {
			return u.Email == email &&
				u.DisplayName == displayName &&
				u.PasswordHash != "" &&
				u.PasswordHash != password
		})).
		Return(nil).
		Once()

	err := svc.Register(email, password, displayName)

	assert.NoError(t, err)
	users.AssertExpectations(t)
}


func TestAuthService_Login_InvalidCredentials_UserNil(t *testing.T) {
	users := new(userRepoMockForAuthService)
	tokens := newTestTokenRepo(t)

	svc := NewAuthService(users, tokens, "test-secret")

	email := "notfound@example.com"

	users.
		On("FindByEmail", email).
		Return((*models.User)(nil), nil).
		Once()

	access, refresh, err := svc.Login(email, "any-password")

	assert.Error(t, err)
	assert.EqualError(t, err, "invalid credentials")
	assert.Equal(t, "", access)
	assert.Equal(t, "", refresh)
	users.AssertExpectations(t)
}

func TestAuthService_Login_InvalidCredentials_WrongPassword(t *testing.T) {
	users := new(userRepoMockForAuthService)
	tokens := newTestTokenRepo(t)

	svc := NewAuthService(users, tokens, "test-secret")

	email := "user@example.com"
	hash, _ := bcrypt.GenerateFromPassword([]byte("correct-password"), 10)

	user := &models.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(hash),
		Role:         models.UserRoleStudent,
	}

	users.
		On("FindByEmail", email).
		Return(user, nil).
		Once()

	access, refresh, err := svc.Login(email, "wrong-password")

	assert.Error(t, err)
	assert.EqualError(t, err, "invalid credentials")
	assert.Equal(t, "", access)
	assert.Equal(t, "", refresh)
	users.AssertExpectations(t)
}

func TestAuthService_Login_Success(t *testing.T) {
	users := new(userRepoMockForAuthService)
	tokens := newTestTokenRepo(t)

	svc := NewAuthService(users, tokens, "test-secret")

	email := "user@example.com"
	password := "correct-password"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	user := &models.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(hash),
		Role:         models.UserRoleStudent,
	}

	users.
		On("FindByEmail", email).
		Return(user, nil).
		Once()

	access, refresh, err := svc.Login(email, password)

	assert.NoError(t, err)
	assert.NotEmpty(t, access)
	assert.NotEmpty(t, refresh)

	users.AssertExpectations(t)
}


func TestAuthService_Refresh_InvalidToken_NotFound(t *testing.T) {
	users := new(userRepoMockForAuthService)
	tokens := newTestTokenRepo(t)

	svc := NewAuthService(users, tokens, "test-secret")

	access, refresh, err := svc.Refresh("some-random-token")

	assert.Error(t, err)
	assert.EqualError(t, err, "invalid refresh token")
	assert.Equal(t, "", access)
	assert.Equal(t, "", refresh)
}

func TestAuthService_Refresh_ExpiredToken(t *testing.T) {
	users := new(userRepoMockForAuthService)
	tokens := newTestTokenRepo(t)

	svc := NewAuthService(users, tokens, "test-secret")

	userID := uuid.New()
	raw := "expired-token"
	hash := svc.generateHash(raw)

	expired := &models.RefreshToken{
		UserID:    userID,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(-time.Hour),
	}
	err := tokens.SaveToken(expired)
	assert.NoError(t, err)

	access, refresh, err := svc.Refresh(raw)

	assert.Error(t, err)
	assert.EqualError(t, err, "token expired")
	assert.Equal(t, "", access)
	assert.Equal(t, "", refresh)
}

func TestAuthService_Refresh_Success(t *testing.T) {
	users := new(userRepoMockForAuthService)
	tokens := newTestTokenRepo(t)

	svc := NewAuthService(users, tokens, "test-secret")

	userID := uuid.New()
	raw := "valid-refresh"
	hash := svc.generateHash(raw)

	rt := &models.RefreshToken{
		UserID:    userID,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(time.Hour),
	}
	err := tokens.SaveToken(rt)
	assert.NoError(t, err)

	user := &models.User{
		ID:   userID,
		Role: models.UserRoleStudent,
	}
	users.
		On("FindByID", userID.String()).
		Return(user, nil).
		Once()

	access, newRefresh, err := svc.Refresh(raw)

	assert.NoError(t, err)
	assert.NotEmpty(t, access)
	assert.NotEmpty(t, newRefresh)

	users.AssertExpectations(t)
}
