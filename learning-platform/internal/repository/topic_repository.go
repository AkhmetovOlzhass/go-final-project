package repository

import (
    "learning-platform/internal/models"
    "gorm.io/gorm"
)

type ITopicRepository interface {
    Create(topic *models.Topic) error
    FindAll() ([]models.Topic, error)
    FindByID(id string) (*models.Topic, error)
    Update(topic *models.Topic) error
    Delete(id string) error
}

type TopicRepository struct {
    db *gorm.DB
}

func NewTopicRepository(db *gorm.DB) *TopicRepository {
    return &TopicRepository{db: db}
}

func (r *TopicRepository) Create(topic *models.Topic) error {
    return r.db.Create(topic).Error
}

func (r *TopicRepository) FindAll() ([]models.Topic, error) {
    var topics []models.Topic
    err := r.db.Find(&topics).Error
    return topics, err
}

func (r *TopicRepository) FindByID(id string) (*models.Topic, error) {
    var topic models.Topic
    err := r.db.First(&topic, "id = ?", id).Error
    return &topic, err
}

func (r *TopicRepository) Update(topic *models.Topic) error {
    return r.db.Save(topic).Error
}

func (r *TopicRepository) Delete(id string) error {
    return r.db.Delete(&models.Topic{}, "id = ?", id).Error
}
