package service

import (
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"time"
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"go.opentelemetry.io/otel"

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

func (s *AuthService) Register(ctx context.Context, email, password, name string) error {
	ctx, span := otel.Tracer("auth").Start(ctx, "AuthService.Register")
	defer span.End()

	_, checkSpan := otel.Tracer("auth").Start(ctx, "CheckExistingUser")
	existing, _ := s.users.FindByEmail(ctx, email)
	checkSpan.End()

	if existing != nil && existing.Status == "ACTIVE" {
		err := errors.New("email already registered")
		span.RecordError(err)
		return err
	}

	if existing != nil && existing.Status != "ACTIVE" {
		_, codeSpan := otel.Tracer("auth").Start(ctx, "GenerateVerificationCode (resend)")
		code, err := generateVerificationCode()
		codeSpan.End()

		if err != nil {
			span.RecordError(err)
			return err
		}

		_, verSpan := otel.Tracer("auth").Start(ctx, "DB.SaveVerification (resend)")
		model := &models.EmailVerification{
			ID:        uuid.New(),
			UserID:    existing.ID,
			Code:      code,
			ExpiresAt: time.Now().Add(15 * time.Minute),
		}
		err = s.verifications.Create(ctx, model)
		verSpan.End()

		if err != nil {
			span.RecordError(err)
			return err
		}

		_, kSpan := otel.Tracer("kafka").Start(ctx, "Kafka.SendEmailCode (resend)")
		s.emailProducer.SendAsync(kafka.EmailMessage{
			Email:   email,
			Subject: "Verify your account",
			Code:    code,
		})
		kSpan.End()

		return nil
	}

	_, createSpan := otel.Tracer("auth").Start(ctx, "DB.CreateUser")
	pass, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &models.User{
		Email:        email,
		PasswordHash: string(pass),
		DisplayName:  name,
	}
	err := s.users.Create(ctx, user);
	createSpan.End()

	if err != nil {
		span.RecordError(err)
		return err
	}

	_, codeSpan := otel.Tracer("auth").Start(ctx, "GenerateVerificationCode")
	code, err := generateVerificationCode()
	codeSpan.End()

	if err != nil {
		span.RecordError(err)
		return err
	}

	_, verSpan := otel.Tracer("auth").Start(ctx, "DB.SaveVerification")
	model := &models.EmailVerification{
		ID:        uuid.New(),
		UserID:    user.ID,
		Code:      code,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}
	err = s.verifications.Create(ctx, model);
	verSpan.End()

	if err != nil {
		span.RecordError(err)
		return err
	}

	_, kSpan := otel.Tracer("kafka").Start(ctx, "Kafka.SendEmailCode")
	s.emailProducer.SendAsync(kafka.EmailMessage{
		Email:   email,
		Subject: "Verify your account",
		Code:    code,
	})
	kSpan.End()

	return nil
}

func (s *AuthService) VerifyEmail(ctx context.Context, email, code string) error {
	ctx, span := otel.Tracer("auth").Start(ctx, "AuthService.VerifyEmail")
	defer span.End()

	_, findSpan := otel.Tracer("auth").Start(ctx, "Verification.FindValid")
	rec, err := s.verifications.FindValid(ctx, email, code)
	findSpan.End()

	if err != nil {
		span.RecordError(err)
		return errors.New("invalid or expired code")
	}

	_, markSpan := otel.Tracer("auth").Start(ctx, "Verification.MarkUsed")
	err = s.verifications.MarkUsed(ctx, rec.ID)
	markSpan.End()

	if err != nil {
		span.RecordError(err)
		return err
	}

	_, userSpan := otel.Tracer("auth").Start(ctx, "User.FindByID")
	user, err := s.users.FindByID(ctx, rec.UserID.String())
	userSpan.End()

	if err != nil || user == nil {
		err = errors.New("user not found")
		span.RecordError(err)
		return err
	}

	_, updateSpan := otel.Tracer("auth").Start(ctx, "User.UpdateStatus")
	err = s.users.Update(ctx, rec.UserID.String(), map[string]interface{}{"status": "ACTIVE"})
	updateSpan.End()

	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, string, error) {
	ctx, span := otel.Tracer("auth").Start(ctx, "AuthService.Login")
	defer span.End()

	_, findSpan := otel.Tracer("auth").Start(ctx, "User.FindByEmail")
	user, err := s.users.FindByEmail(ctx, email)
	findSpan.End()

	if err != nil || user == nil {
		err := errors.New("invalid credentials")
		span.RecordError(err)
		return "", "", err
	}

	_, checkSpan := otel.Tracer("auth").Start(ctx, "Password.Verify")
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	checkSpan.End()

	if err != nil {
		err2 := errors.New("invalid credentials")
		span.RecordError(err2)
		return "", "", err2
	}

	if user.Status != "ACTIVE" {
		err := errors.New("email not verified")
		span.RecordError(err)
		return "", "", err
	}

	_, revokeSpan := otel.Tracer("auth").Start(ctx, "RefreshTokens.RevokeAll")
	err = s.tokens.RevokeAllForUser(ctx, user.ID)
	revokeSpan.End()

	if err != nil {
		span.RecordError(err)
		return "", "", err
	}

	_, jwtSpan := otel.Tracer("auth").Start(ctx, "JWT.CreateAccess")
	access, err := s.createJWT(user.ID, user.Role)
	jwtSpan.End()

	if err != nil {
		span.RecordError(err)
		return "", "", err
	}

	rawRefresh := uuid.New().String()

	refreshModel := models.RefreshToken{
		UserID:    user.ID,
		TokenHash: s.generateHash(rawRefresh),
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}

	_, saveSpan := otel.Tracer("auth").Start(ctx, "RefreshTokens.Save")
	err = s.tokens.Save(ctx, &refreshModel)
	saveSpan.End()

	if err != nil {
		span.RecordError(err)
		return "", "", err
	}

	return access, rawRefresh, nil
}

func (s *AuthService) Refresh(ctx context.Context, oldRefresh string) (string, string, error) {
	ctx, span := otel.Tracer("auth").Start(ctx, "AuthService.Refresh")
	defer span.End()

	hash := s.generateHash(oldRefresh)

	_, findSpan := otel.Tracer("auth").Start(ctx, "RefreshTokens.FindValid")
	token, err := s.tokens.FindValid(ctx, hash)
	findSpan.End()

	if err != nil || token == nil {
		err = errors.New("invalid refresh token")
		span.RecordError(err)
		return "", "", err
	}

	_, revokeSpan := otel.Tracer("auth").Start(ctx, "RefreshTokens.Revoke")
	err = s.tokens.Revoke(ctx, hash)
	revokeSpan.End()

	if err != nil {
		span.RecordError(err)
		return "", "", err
	}

	_, userSpan := otel.Tracer("auth").Start(ctx, "User.FindByID")
	user, err := s.users.FindByID(ctx, token.UserID.String())
	userSpan.End()

	if err != nil || user == nil {
		err = errors.New("user not found")
		span.RecordError(err)
		return "", "", err
	}

	_, accessSpan := otel.Tracer("auth").Start(ctx, "JWT.CreateAccess")
	newAccess, err := s.createJWT(token.UserID, user.Role)
	accessSpan.End()

	if err != nil {
		span.RecordError(err)
		return "", "", err
	}

	newRefreshRaw := uuid.NewString()
	newModel := &models.RefreshToken{
		UserID:    token.UserID,
		TokenHash: s.generateHash(newRefreshRaw),
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}

	_, saveSpan := otel.Tracer("auth").Start(ctx, "RefreshTokens.SaveNew")
	err = s.tokens.Save(ctx, newModel)
	saveSpan.End()

	if err != nil {
		span.RecordError(err)
		return "", "", err
	}

	return newAccess, newRefreshRaw, nil
}

func (s *AuthService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	ctx, span := otel.Tracer("auth").Start(ctx, "AuthService.GetUserByID")
	defer span.End()

	user, err := s.users.FindByID(ctx, id)

	if err != nil {
		span.RecordError(err)
	}

	return user, err
}
