package repository

import (
	"context"

	"learning-platform/internal/models"

	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
)

type ITopicRepository interface {
	Create(ctx context.Context, topic *models.Topic) error
	FindAll(ctx context.Context) ([]models.Topic, error)
	FindByID(ctx context.Context, id string) (*models.Topic, error)
	Update(ctx context.Context, topic *models.Topic) error
	Delete(ctx context.Context, id string) error
}

type TopicRepository struct {
	db *gorm.DB
}

func NewTopicRepository(db *gorm.DB) *TopicRepository {
	return &TopicRepository{db: db}
}

func (r *TopicRepository) Create(ctx context.Context, topic *models.Topic) error {
	ctx, span := otel.Tracer("db").Start(ctx, "TopicRepository.Create")
	defer span.End()

	err := r.db.WithContext(ctx).Create(topic).Error
	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (r *TopicRepository) FindAll(ctx context.Context) ([]models.Topic, error) {
	ctx, span := otel.Tracer("db").Start(ctx, "TopicRepository.FindAll")
	defer span.End()

	var topics []models.Topic
	err := r.db.WithContext(ctx).Find(&topics).Error
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return topics, nil
}

func (r *TopicRepository) FindByID(ctx context.Context, id string) (*models.Topic, error) {
	ctx, span := otel.Tracer("db").Start(ctx, "TopicRepository.FindByID")
	defer span.End()

	var topic models.Topic
	err := r.db.WithContext(ctx).First(&topic, "id = ?", id).Error
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return &topic, nil
}

func (r *TopicRepository) Update(ctx context.Context, topic *models.Topic) error {
	ctx, span := otel.Tracer("db").Start(ctx, "TopicRepository.Update")
	defer span.End()

	err := r.db.WithContext(ctx).Save(topic).Error
	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (r *TopicRepository) Delete(ctx context.Context, id string) error {
	ctx, span := otel.Tracer("db").Start(ctx, "TopicRepository.Delete")
	defer span.End()

	err := r.db.WithContext(ctx).Delete(&models.Topic{}, "id = ?", id).Error
	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}
