package service

import (
	"learning-platform/internal/models"
	"learning-platform/internal/repository"
)

type TaskService struct {
	taskRepo repository.ITaskRepository
}

func NewTaskService(taskRepo repository.ITaskRepository) *TaskService {
	return &TaskService{taskRepo: taskRepo}
}

func (s *TaskService) GetAllTasks() ([]models.Task, error) {
	return s.taskRepo.GetAll()
}

func (s *TaskService) GetDraftTasks() ([]models.Task, error) {
	return s.taskRepo.GetDraft()
}

func (s *TaskService) PublishTask(id string) (*models.Task, error) {

	if err := s.taskRepo.UpdateStatus(id, models.TaskStatusPublished); err != nil {
		return nil, err
	}

	updatedTask, err := s.taskRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return updatedTask, nil
}

func (s *TaskService) CreateTask(task *models.Task) error {
	return s.taskRepo.Create(task)
}

func (s *TaskService) GetTaskById(id string) (*models.Task, error) {
	return s.taskRepo.GetByID(id)
}

func (s *TaskService) GetTasksByTopic(topicID string) ([]models.Task, error) {
	return s.taskRepo.GetByTopic(topicID)
}

func (s *TaskService) UpdateTask(task *models.Task) error {
	return s.taskRepo.Update(task)
}

func (s *TaskService) DeleteTask(id string) error {
	return s.taskRepo.Delete(id)
}

func (s *TaskService) GetTasksByAuthor(authorID string) ([]models.Task, error) {
	return s.taskRepo.GetByAuthor(authorID)
}
