package repository

import (
	"learning-platform/internal/models"

	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IVerificationRepository interface {
	Create(model *models.EmailVerification) error
	FindValid(email, code string) (*models.EmailVerification, error)
	MarkUsed(id uuid.UUID) error
}

type VerificationRepository struct {
	db *gorm.DB
}

func NewVerificationRepository(db *gorm.DB) *VerificationRepository {
	return &VerificationRepository{db: db}
}

func (r *VerificationRepository) Create(model *models.EmailVerification) error {
	return r.db.Create(model).Error
}

func (r *VerificationRepository) FindValid(email, code string) (*models.EmailVerification, error) {
	var v models.EmailVerification
	return &v, r.db.
		Joins("JOIN users u ON u.id = email_verifications.user_id").
		Where("u.email = ? AND email_verifications.code = ? AND email_verifications.used = FALSE AND email_verifications.expires_at > ?", email, code, time.Now()).
		First(&v).Error
}

func (r *VerificationRepository) MarkUsed(id uuid.UUID) error {
	return r.db.Model(&models.EmailVerification{}).
		Where("id = ?", id).
		Update("used", true).Error
}
