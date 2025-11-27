package service

import (
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"learning-platform/internal/kafka"
	"learning-platform/internal/models"
	"learning-platform/internal/repository"
)

func generateVerificationCode() (string, error) {
	max := big.NewInt(1000000)
	n, err := crand.Int(crand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

type AuthService struct {
	users         repository.IUserRepository
	verifications repository.IVerificationRepository
	tokens        repository.ITokenRepository
	emailProducer *kafka.EmailProducer
	secret        string
}

func NewAuthService(u repository.IUserRepository, v repository.IVerificationRepository, t repository.ITokenRepository, p *kafka.EmailProducer, secret string) *AuthService {
	return &AuthService{
		users:         u,
		verifications: v,
		tokens:        t,
		emailProducer: p,
		secret:        secret,
	}
}

func (s *AuthService) generateHash(value string) string {
	hash := sha256.Sum256([]byte(value))
	return hex.EncodeToString(hash[:])
}

func (s *AuthService) createJWT(userID uuid.UUID, role models.UserRole) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userID.String(),
		"role":   string(role),
		"exp":    time.Now().Add(15 * time.Minute).Unix(),
	})
	return token.SignedString([]byte(s.secret))
}

func (s *AuthService) Register(email, password, name string) error {
	existing, _ := s.users.FindByEmail(email)
	if existing != nil {
		return errors.New("email already registered")
	}

	pass, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &models.User{
		Email:        email,
		PasswordHash: string(pass),
		DisplayName:  name,
	}

	if err := s.users.Create(user); err != nil {
		return err
	}

	code, err := generateVerificationCode()
	if err != nil {
		return err
	}

	model := &models.EmailVerification{
		ID:        uuid.New(),
		UserID:    user.ID,
		Code:      code,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	if err := s.verifications.Create(model); err != nil {
		return err
	}

	s.emailProducer.SendAsync(kafka.EmailMessage{
		Email:   email,
		Subject: "Verify your account",
		Code:    code,
	})

	return nil
}

func (s *AuthService) VerifyEmail(email, code string) error {
	rec, err := s.verifications.FindValid(email, code)
	if err != nil {
		return errors.New("invalid or expired code")
	}

	if err := s.verifications.MarkUsed(rec.ID); err != nil {
		return err
	}

	userId := rec.UserID.String()

	user, err := s.users.FindByID(userId)
	if err != nil || user == nil {
		return errors.New("user not found")
	}

	updateFields := map[string]interface{}{
		"status": "ACTIVE",
	}

	return s.users.Update(userId, updateFields)
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

	if user.Status != "ACTIVE" {
		return "", "", errors.New("email not verified")
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
