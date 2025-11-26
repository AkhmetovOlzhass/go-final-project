package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"learning-platform/internal/models"
)

type ITokenRepository interface {
	Save(token *models.RefreshToken) error
	FindValid(hash string) (*models.RefreshToken, error)
	Revoke(hash string) error
	RevokeAllForUser(userID uuid.UUID) error
	DeleteExpired() error
}

type TokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) Save(token *models.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *TokenRepository) FindValid(hash string) (*models.RefreshToken, error) {
	var token models.RefreshToken

	err := r.db.
		Where("token_hash = ? AND revoked = ? AND expires_at > ?", hash, false, time.Now()).
		First(&token).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &token, err
}

func (r *TokenRepository) Revoke(hash string) error {
	return r.db.Model(&models.RefreshToken{}).
		Where("token_hash = ? AND revoked = ?", hash, false).
		Update("revoked", true).Error
}

func (r *TokenRepository) RevokeAllForUser(userID uuid.UUID) error {
	return r.db.Model(&models.RefreshToken{}).
		Where("user_id = ? AND revoked = ?", userID, false).
		Update("revoked", true).Error
}

func (r *TokenRepository) DeleteExpired() error {
	return r.db.Where("expires_at < ?", time.Now()).
		Delete(&models.RefreshToken{}).Error
}
