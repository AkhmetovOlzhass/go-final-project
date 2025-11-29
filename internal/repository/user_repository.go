package repository

import (
	"context"
	"errors"

	"learning-platform/internal/models"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
)

type IUserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id string) (*models.User, error)
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	GetAll(ctx context.Context) ([]models.User, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	ctx, span := otel.Tracer("db").Start(ctx, "UserRepository.GetAll")
	defer span.End()

	var users []models.User
	err := r.db.WithContext(ctx).Find(&users).Error
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	ctx, span := otel.Tracer("db").Start(ctx, "UserRepository.Create")
	defer span.End()

	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	ctx, span := otel.Tracer("db").Start(ctx, "UserRepository.FindByEmail")
	defer span.End()

	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	ctx, span := otel.Tracer("db").Start(ctx, "UserRepository.FindByID")
	defer span.End()

	uid, err := uuid.Parse(id)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	var user models.User
	err = r.db.WithContext(ctx).First(&user, "id = ?", uid).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	ctx, span := otel.Tracer("db").Start(ctx, "UserRepository.Update")
	defer span.End()

	err := r.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", id).
		Updates(updates).Error

	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}
