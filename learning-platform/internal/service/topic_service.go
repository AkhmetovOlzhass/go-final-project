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

func (s *TopicService) CreateTopic(topic *models.Topic) error {
    return s.repo.Create(topic)
}

func (s *TopicService) GetAllTopics() ([]models.Topic, error) {
    return s.repo.FindAll()
}

func (s *TopicService) GetTopicById(id string) (*models.Topic, error) {
    return s.repo.FindByID(id)
}

func (s *TopicService) UpdateTopic(topic *models.Topic) error {
    return s.repo.Update(topic)
}

func (s *TopicService) DeleteTopic(id string) error {
    return s.repo.Delete(id)
}
