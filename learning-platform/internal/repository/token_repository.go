package repository

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"github.com/google/uuid"
	"learning-platform/internal/models"
)

type TokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) SaveToken(token *models.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *TokenRepository) FindByHash(hash string) (*models.RefreshToken, error) {
	var token models.RefreshToken
	err := r.db.Where("token_hash = ? AND revoked = FALSE", hash).First(&token).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &token, err
}

func (r *TokenRepository) RevokeToken(hash string) error {
	return r.db.Model(&models.RefreshToken{}).
		Where("token_hash = ?", hash).
		Update("revoked", true).Error
}

func (r *TokenRepository) DeleteExpired() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&models.RefreshToken{}).Error
}

func (r *TokenRepository) RevokeAllForUser(userID uuid.UUID) error {
	return r.db.Model(&models.RefreshToken{}).
		Where("user_id = ? AND revoked = FALSE", userID).
		Update("revoked", true).Error
}