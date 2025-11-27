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

func (s *TaskService) GetAllTasks() ([]models.Task, error) {
	return s.taskRepo.GetAllTasks()
}

func (s *TaskService) GetDraftTasks() ([]models.Task, error) {
	return s.taskRepo.GetDraftTasks()
}

func (s *TaskService) PublishTask(id string) (*models.Task, error) {

	if err := s.taskRepo.UpdateStatus(id, models.TaskStatusPublished); err != nil {
        return nil, err
    }

    updatedTask, err := s.taskRepo.GetTaskByID(id)
    if err != nil {
        return nil, err
    }

    return updatedTask, nil
}

func (s *TaskService) CreateTask(task *models.Task) error {
	return s.taskRepo.CreateTask(task)
}

func (s *TaskService) GetTaskById(id string) (*models.Task, error) {
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