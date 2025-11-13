package service

import (
    "learning-platform/internal/models"
    "learning-platform/internal/repository"
)

type TopicService struct {
    repo repository.ITopicRepository
}

func NewTopicService(repo repository.ITopicRepository) *TopicService {
    return &TopicService{repo: repo}
}

func (s *TopicService) Create(topic *models.Topic) error {
    return s.repo.Create(topic)
}

func (s *TopicService) GetAll() ([]models.Topic, error) {
    return s.repo.FindAll()
}

func (s *TopicService) GetByID(id string) (*models.Topic, error) {
    return s.repo.FindByID(id)
}

func (s *TopicService) Update(topic *models.Topic) error {
    return s.repo.Update(topic)
}

func (s *TopicService) Delete(id string) error {
    return s.repo.Delete(id)
}
