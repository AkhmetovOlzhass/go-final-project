package repository

import (
	"context"
	"time"

	"learning-platform/internal/models"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
)

type IVerificationRepository interface {
	Create(ctx context.Context, model *models.EmailVerification) error
	FindValid(ctx context.Context, email, code string) (*models.EmailVerification, error)
	MarkUsed(ctx context.Context, id uuid.UUID) error
}

type VerificationRepository struct {
	db *gorm.DB
}

func NewVerificationRepository(db *gorm.DB) *VerificationRepository {
	return &VerificationRepository{db: db}
}

func (r *VerificationRepository) Create(ctx context.Context, model *models.EmailVerification) error {
	ctx, span := otel.Tracer("db").Start(ctx, "VerificationRepository.Create")
	defer span.End()

	err := r.db.WithContext(ctx).Create(model).Error
	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (r *VerificationRepository) FindValid(ctx context.Context, email, code string) (*models.EmailVerification, error) {
	ctx, span := otel.Tracer("db").Start(ctx, "VerificationRepository.FindValid")
	defer span.End()

	var v models.EmailVerification
	err := r.db.WithContext(ctx).
		Joins("JOIN users u ON u.id = email_verifications.user_id").
		Where(
			"u.email = ? AND email_verifications.code = ? AND email_verifications.used = FALSE AND email_verifications.expires_at > ?",
			email, code, time.Now(),
		).
		First(&v).Error

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return &v, nil
}

func (r *VerificationRepository) MarkUsed(ctx context.Context, id uuid.UUID) error {
	ctx, span := otel.Tracer("db").Start(ctx, "VerificationRepository.MarkUsed")
	defer span.End()

	err := r.db.WithContext(ctx).
		Model(&models.EmailVerification{}).
		Where("id = ?", id).
		Update("used", true).Error

	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}
