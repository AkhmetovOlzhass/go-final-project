package repository

import (
	"context"

	"learning-platform/internal/models"

	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
)

type ITaskRepository interface {
	GetAll(ctx context.Context) ([]models.Task, error)
	GetDraft(ctx context.Context) ([]models.Task, error)
	UpdateStatus(ctx context.Context, id string, status models.TaskStatus) error
	Create(ctx context.Context, task *models.Task) error
	GetByID(ctx context.Context, id string) (*models.Task, error)
	GetByTopic(ctx context.Context, topicID string) ([]models.Task, error)
	Update(ctx context.Context, task *models.Task) error
	Delete(ctx context.Context, id string) error
	GetByAuthor(ctx context.Context, authorID string) ([]models.Task, error)
}

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) GetAll(ctx context.Context) ([]models.Task, error) {
	ctx, span := otel.Tracer("db").Start(ctx, "TaskRepository.GetAll")
	defer span.End()

	var tasks []models.Task
	err := r.db.WithContext(ctx).Find(&tasks).Error
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return tasks, nil
}

func (r *TaskRepository) GetDraft(ctx context.Context) ([]models.Task, error) {
	ctx, span := otel.Tracer("db").Start(ctx, "TaskRepository.GetDraft")
	defer span.End()

	var tasks []models.Task
	err := r.db.WithContext(ctx).
		Where("status = ?", models.TaskStatusDraft).
		Find(&tasks).Error

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return tasks, nil
}

func (r *TaskRepository) UpdateStatus(ctx context.Context, id string, status models.TaskStatus) error {
	ctx, span := otel.Tracer("db").Start(ctx, "TaskRepository.UpdateStatus")
	defer span.End()

	err := r.db.WithContext(ctx).
		Model(&models.Task{}).
		Where("id = ?", id).
		Update("status", status).
		Error

	if err != nil {
		span.RecordError(err)
	}

	return err
}

func (r *TaskRepository) Create(ctx context.Context, task *models.Task) error {
	ctx, span := otel.Tracer("db").Start(ctx, "TaskRepository.Create")
	defer span.End()

	err := r.db.WithContext(ctx).Create(task).Error
	if err != nil {
		span.RecordError(err)
	}

	return err
}

func (r *TaskRepository) GetByID(ctx context.Context, id string) (*models.Task, error) {
	ctx, span := otel.Tracer("db").Start(ctx, "TaskRepository.GetByID")
	defer span.End()

	var task models.Task
	err := r.db.WithContext(ctx).
		Preload("Topic").
		Preload("Author").
		Where("id = ?", id).
		First(&task).Error

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return &task, nil
}

func (r *TaskRepository) GetByTopic(ctx context.Context, topicID string) ([]models.Task, error) {
	ctx, span := otel.Tracer("db").Start(ctx, "TaskRepository.GetByTopic")
	defer span.End()

	var tasks []models.Task
	err := r.db.WithContext(ctx).
		Where("topic_id = ? AND status = ?", topicID, models.TaskStatusPublished).
		Order("created_at DESC").
		Find(&tasks).Error

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return tasks, nil
}

func (r *TaskRepository) Update(ctx context.Context, task *models.Task) error {
	ctx, span := otel.Tracer("db").Start(ctx, "TaskRepository.Update")
	defer span.End()

	err := r.db.WithContext(ctx).Save(task).Error
	if err != nil {
		span.RecordError(err)
	}

	return err
}

func (r *TaskRepository) Delete(ctx context.Context, id string) error {
	ctx, span := otel.Tracer("db").Start(ctx, "TaskRepository.Delete")
	defer span.End()

	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&models.Task{}).Error

	if err != nil {
		span.RecordError(err)
	}

	return err
}

func (r *TaskRepository) GetByAuthor(ctx context.Context, authorID string) ([]models.Task, error) {
	ctx, span := otel.Tracer("db").Start(ctx, "TaskRepository.GetByAuthor")
	defer span.End()

	var tasks []models.Task
	err := r.db.WithContext(ctx).
		Where("author_id = ?", authorID).
		Order("created_at DESC").
		Find(&tasks).Error

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return tasks, nil
}
