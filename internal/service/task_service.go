package service

import (
	"context"

	"learning-platform/internal/models"
	"learning-platform/internal/repository"
	"go.opentelemetry.io/otel"
)

type TaskService struct {
	taskRepo repository.ITaskRepository
}

func NewTaskService(taskRepo repository.ITaskRepository) *TaskService {
	return &TaskService{taskRepo: taskRepo}
}

func (s *TaskService) GetAllTasks(ctx context.Context) ([]models.Task, error) {
	ctx, span := otel.Tracer("task").Start(ctx, "TaskService.GetAllTasks")
	defer span.End()

	tasks, err := s.taskRepo.GetAll(ctx)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return tasks, nil
}

func (s *TaskService) GetDraftTasks(ctx context.Context) ([]models.Task, error) {
	ctx, span := otel.Tracer("task").Start(ctx, "TaskService.GetDraftTasks")
	defer span.End()

	tasks, err := s.taskRepo.GetDraft(ctx)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return tasks, nil
}

func (s *TaskService) PublishTask(ctx context.Context, id string) (*models.Task, error) {
	ctx, span := otel.Tracer("task").Start(ctx, "TaskService.PublishTask")
	defer span.End()

	err := s.taskRepo.UpdateStatus(ctx, id, models.TaskStatusPublished)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return task, nil
}

func (s *TaskService) CreateTask(ctx context.Context, task *models.Task) error {
	ctx, span := otel.Tracer("task").Start(ctx, "TaskService.CreateTask")
	defer span.End()

	err := s.taskRepo.Create(ctx, task)
	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (s *TaskService) GetTaskById(ctx context.Context, id string) (*models.Task, error) {
	ctx, span := otel.Tracer("task").Start(ctx, "TaskService.GetTaskById")
	defer span.End()

	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return task, nil
}

func (s *TaskService) GetTasksByTopic(ctx context.Context, topicID string) ([]models.Task, error) {
	ctx, span := otel.Tracer("task").Start(ctx, "TaskService.GetTasksByTopic")
	defer span.End()

	tasks, err := s.taskRepo.GetByTopic(ctx, topicID)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return tasks, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, task *models.Task) error {
	ctx, span := otel.Tracer("task").Start(ctx, "TaskService.UpdateTask")
	defer span.End()

	err := s.taskRepo.Update(ctx, task)
	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (s *TaskService) DeleteTask(ctx context.Context, id string) error {
	ctx, span := otel.Tracer("task").Start(ctx, "TaskService.DeleteTask")
	defer span.End()

	err := s.taskRepo.Delete(ctx, id)
	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (s *TaskService) GetTasksByAuthor(ctx context.Context, authorID string) ([]models.Task, error) {
	ctx, span := otel.Tracer("task").Start(ctx, "TaskService.GetTasksByAuthor")
	defer span.End()

	tasks, err := s.taskRepo.GetByAuthor(ctx, authorID)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return tasks, nil
}
