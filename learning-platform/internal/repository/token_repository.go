package repository

import (
	"context"
	"errors"
	"time"

	"learning-platform/internal/models"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
)

type ITokenRepository interface {
	Save(ctx context.Context, token *models.RefreshToken) error
	FindValid(ctx context.Context, hash string) (*models.RefreshToken, error)
	Revoke(ctx context.Context, hash string) error
	RevokeAllForUser(ctx context.Context, userID uuid.UUID) error
	DeleteExpired(ctx context.Context) error
}

type TokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) Save(ctx context.Context, token *models.RefreshToken) error {
	ctx, span := otel.Tracer("db").Start(ctx, "TokenRepository.Save")
	defer span.End()

	err := r.db.WithContext(ctx).Create(token).Error
	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (r *TokenRepository) FindValid(ctx context.Context, hash string) (*models.RefreshToken, error) {
	ctx, span := otel.Tracer("db").Start(ctx, "TokenRepository.FindValid")
	defer span.End()

	var token models.RefreshToken

	err := r.db.WithContext(ctx).
		Where("token_hash = ? AND revoked = ? AND expires_at > ?", hash, false, time.Now()).
		First(&token).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return &token, nil
}

func (r *TokenRepository) Revoke(ctx context.Context, hash string) error {
	ctx, span := otel.Tracer("db").Start(ctx, "TokenRepository.Revoke")
	defer span.End()

	err := r.db.WithContext(ctx).
		Model(&models.RefreshToken{}).
		Where("token_hash = ? AND revoked = ?", hash, false).
		Update("revoked", true).Error

	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (r *TokenRepository) RevokeAllForUser(ctx context.Context, userID uuid.UUID) error {
	ctx, span := otel.Tracer("db").Start(ctx, "TokenRepository.RevokeAllForUser")
	defer span.End()

	err := r.db.WithContext(ctx).
		Model(&models.RefreshToken{}).
		Where("user_id = ? AND revoked = ?", userID, false).
		Update("revoked", true).Error

	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (r *TokenRepository) DeleteExpired(ctx context.Context) error {
	ctx, span := otel.Tracer("db").Start(ctx, "TokenRepository.DeleteExpired")
	defer span.End()

	err := r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&models.RefreshToken{}).Error

	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}
