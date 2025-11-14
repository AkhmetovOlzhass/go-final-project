package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"learning-platform/internal/models"
	"learning-platform/internal/repository"
)

type AuthService struct {
	users  repository.IUserRepository
	tokens *repository.TokenRepository
	secret string
}

func NewAuthService(u repository.IUserRepository, t *repository.TokenRepository, secret string) *AuthService {
	return &AuthService{users: u, tokens: t, secret: secret}
}


func (s *AuthService) generateHash(value string) string {
	hash := sha256.Sum256([]byte(value))
	return hex.EncodeToString(hash[:])
}

func (s *AuthService) createJWT(userID uuid.UUID, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID.String(),
		"role":    role,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	})
	return token.SignedString([]byte(s.secret))
}

func (s *AuthService) Register(email, password, displayName string) error {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	user := &models.User{Email: email, PasswordHash: string(hash), DisplayName:  displayName,}
	return s.users.Create(user)
}

func (s *AuthService) Login(email, password string) (string, string, error) {
	user, err := s.users.FindByEmail(email)
	if err != nil {
		return "", "", err
	}
	if user == nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return "", "", errors.New("invalid credentials")
	}

	_ = s.tokens.RevokeAllForUser(user.ID)

	access, err := s.createJWT(user.ID, string(user.Role))
	if err != nil {
		return "", "", err
	}

	rawRefresh := uuid.NewString()
	refresh := &models.RefreshToken{
		UserID:    user.ID,
		TokenHash: s.generateHash(rawRefresh),
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}
	if err := s.tokens.SaveToken(refresh); err != nil {
		return "", "", err
	}

	return access, rawRefresh, nil
}


func (s *AuthService) Refresh(oldRefresh string) (string, string, error) {
	hash := s.generateHash(oldRefresh)
	token, err := s.tokens.FindByHash(hash)
	if err != nil || token == nil {
		return "", "", errors.New("invalid refresh token")
	}

	if token.ExpiresAt.Before(time.Now()) {
		_ = s.tokens.RevokeToken(hash)
		return "", "", errors.New("token expired")
	}

	_ = s.tokens.RevokeToken(hash)

	user, err := s.users.FindByID(token.UserID.String())
	if err != nil || user == nil {
		return "", "", errors.New("user not found")
	}

	newAccess, err := s.createJWT(token.UserID, string(user.Role))
	if err != nil {
		return "", "", err
	}

	newRefreshValue := uuid.NewString()
	newRefresh := &models.RefreshToken{
		UserID:    token.UserID,
		TokenHash: s.generateHash(newRefreshValue),
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}

	if err := s.tokens.SaveToken(newRefresh); err != nil {
		return "", "", err
	}

	return newAccess, newRefreshValue, nil
}

