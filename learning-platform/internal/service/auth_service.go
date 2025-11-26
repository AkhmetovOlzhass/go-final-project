package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"learning-platform/internal/models"
	"learning-platform/internal/repository"
)

type AuthService struct {
	users  repository.IUserRepository
	tokens repository.ITokenRepository
	secret string
}

func NewAuthService(u repository.IUserRepository, t repository.ITokenRepository, secret string) *AuthService {
	return &AuthService{users: u, tokens: t, secret: secret}
}

func (s *AuthService) generateHash(value string) string {
	hash := sha256.Sum256([]byte(value))
	return hex.EncodeToString(hash[:])
}

func (s *AuthService) createJWT(userID uuid.UUID, role models.UserRole) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userID.String(),
		"role":    string(role),
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	})
	return token.SignedString([]byte(s.secret))
}

func (s *AuthService) Register(email, password, displayName string) error {
	existing, err := s.users.FindByEmail(email)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("email already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		Email:        email,
		PasswordHash: string(hash),
		DisplayName:  displayName,
	}

	return s.users.Create(user)
}

func (s *AuthService) Login(email, password string) (string, string, error) {
	user, err := s.users.FindByEmail(email)
	if err != nil {
		return "", "", fmt.Errorf("find user failed: %w", err)
	}
	if user == nil {
		return "", "", errors.New("invalid credentials")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return "", "", errors.New("invalid credentials")
	}

	if err := s.tokens.RevokeAllForUser(user.ID); err != nil {
		return "", "", err
	}

	access, err := s.createJWT(user.ID, user.Role)
	if err != nil {
		return "", "", err
	}

	rawRefresh := uuid.New().String()
	saveModel := models.RefreshToken{
		UserID:    user.ID,
		TokenHash: s.generateHash(rawRefresh),
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}

	if err := s.tokens.Save(&saveModel); err != nil {
		return "", "", err
	}

	return access, rawRefresh, nil
}

func (s *AuthService) Refresh(oldRefresh string) (string, string, error) {
	hash := s.generateHash(oldRefresh)

	token, err := s.tokens.FindValid(hash)
	if err != nil {
		return "", "", fmt.Errorf("failed to read refresh token: %w", err)
	}
	if token == nil {
		return "", "", errors.New("invalid refresh token")
	}

	if err := s.tokens.Revoke(hash); err != nil {
		return "", "", err
	}

	user, err := s.users.FindByID(token.UserID.String())
	if err != nil {
		return "", "", fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return "", "", errors.New("user not found")
	}

	newAccess, err := s.createJWT(token.UserID, user.Role)
	if err != nil {
		return "", "", err
	}

	newRefreshRaw := uuid.NewString()
	newToken := &models.RefreshToken{
		UserID:    token.UserID,
		TokenHash: s.generateHash(newRefreshRaw),
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}

	if err := s.tokens.Save(newToken); err != nil {
		return "", "", err
	}

	return newAccess, newRefreshRaw, nil
}

func (s *AuthService) GetUserByID(id string) (*models.User, error) {
	return s.users.FindByID(id)
}
