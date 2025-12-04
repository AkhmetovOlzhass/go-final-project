package service

import (
	"context"

	"encoding/json"
	"learning-platform/internal/models"
	"learning-platform/internal/repository"
	"time"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
)

type TaskService struct {
	taskRepo repository.ITaskRepository
	redis    *redis.Client
}

func NewTaskService(repo repository.ITaskRepository, rdb *redis.Client) *TaskService {
	return &TaskService{
		taskRepo: repo,
		redis:    rdb,
	}
}

func (s *TaskService) GetAllTasks(ctx context.Context) ([]models.Task, error) {
	ctx, span := otel.Tracer("task").Start(ctx, "TaskService.GetAllTasks")
	defer span.End()

	cacheKey := "tasks:all"
	if cached, err := s.redis.Get(ctx, cacheKey).Result(); err == nil {
		var tasks []models.Task
		if err := json.Unmarshal([]byte(cached), &tasks); err == nil {
			return tasks, nil
		}
	}

	tasks, err := s.taskRepo.GetAll(ctx)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	data, _ := json.Marshal(tasks)
	s.redis.Set(ctx, cacheKey, data, 10*time.Minute)

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
	s.redis.Del(context.Background(), "tasks:all")
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
	s.redis.Del(context.Background(), "tasks:all")
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
	s.redis.Del(context.Background(), "tasks:all")
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
	s.redis.Del(context.Background(), "tasks:all")
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
