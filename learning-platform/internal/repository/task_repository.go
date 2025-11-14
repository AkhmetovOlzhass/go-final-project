package repository

import (
	"learning-platform/internal/models"

	"gorm.io/gorm"
)

type TaskRepository struct {
	db *gorm.DB 
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) GetAllTasks() ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.Find(&tasks).Error
	return tasks, err
}

func (r *TaskRepository) GetDraftTasks() ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.Where("status = ?", models.TaskStatusDraft).Find(&tasks).Error
	return tasks, err
}

func (r *TaskRepository) UpdateStatus(id string, status models.TaskStatus) error {
    return r.db.Model(&models.Task{}).
        Where("id = ?", id).
        Update("status", status).
        Error
}



func (r *TaskRepository) CreateTask(task *models.Task) error {
	return r.db.Create(task).Error
}

func (r *TaskRepository) GetTaskByID(id string) (*models.Task, error) {
	var task models.Task
	err := r.db.
		Preload("Topic").
		Preload("Author").
		Where("id = ?", id).
		First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepository) GetTasksByTopic(topicID string) ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.
		Where("topic_id = ? AND status = ?", topicID, models.TaskStatusPublished).
		Order("created_at DESC").
		Find(&tasks).Error
	return tasks, err
}

func (r *TaskRepository) UpdateTask(task *models.Task) error {
	return r.db.Save(task).Error
}

func (r *TaskRepository) DeleteTask(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.Task{}).Error
}

func (r *TaskRepository) GetTasksByAuthor(authorID string) ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.
		Where("author_id = ?", authorID).
		Order("created_at DESC").
		Find(&tasks).Error
	return tasks, err
}