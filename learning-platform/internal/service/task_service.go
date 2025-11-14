package service

import (
	"learning-platform/internal/models"
	"learning-platform/internal/repository"
)

type TaskService struct {
	taskRepo *repository.TaskRepository
}

func NewTaskService(taskRepo *repository.TaskRepository) *TaskService {
	return &TaskService{taskRepo: taskRepo}
}

func (s *TaskService) CreateTask(task *models.Task) error {
	return s.taskRepo.CreateTask(task)
}

func (s *TaskService) GetTask(id string) (*models.Task, error) {
	return s.taskRepo.GetTaskByID(id)
}

func (s *TaskService) GetTasksByTopic(topicID string) ([]models.Task, error) {
	return s.taskRepo.GetTasksByTopic(topicID)
}

func (s *TaskService) UpdateTask(task *models.Task) error {
	return s.taskRepo.UpdateTask(task)
}

func (s *TaskService) DeleteTask(id string) error {
	return s.taskRepo.DeleteTask(id)
}

func (s *TaskService) GetTasksByAuthor(authorID string) ([]models.Task, error) {
	return s.taskRepo.GetTasksByAuthor(authorID)
}