package service

import (
	"context"

	"learning-platform/internal/models"
	"learning-platform/internal/repository"
	"go.opentelemetry.io/otel"
)

type TopicService struct {
	repo repository.ITopicRepository
}

func NewTopicService(repo repository.ITopicRepository) *TopicService {
	return &TopicService{repo: repo}
}

func (s *TopicService) CreateTopic(ctx context.Context, topic *models.Topic) error {
	ctx, span := otel.Tracer("topic").Start(ctx, "TopicService.CreateTopic")
	defer span.End()

	err := s.repo.Create(ctx, topic)
	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (s *TopicService) GetAllTopics(ctx context.Context) ([]models.Topic, error) {
	ctx, span := otel.Tracer("topic").Start(ctx, "TopicService.GetAllTopics")
	defer span.End()

	topics, err := s.repo.FindAll(ctx)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return topics, nil
}

func (s *TopicService) GetTopicById(ctx context.Context, id string) (*models.Topic, error) {
	ctx, span := otel.Tracer("topic").Start(ctx, "TopicService.GetTopicById")
	defer span.End()

	topic, err := s.repo.FindByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return topic, nil
}

func (s *TopicService) UpdateTopic(ctx context.Context, topic *models.Topic) error {
	ctx, span := otel.Tracer("topic").Start(ctx, "TopicService.UpdateTopic")
	defer span.End()

	err := s.repo.Update(ctx, topic)
	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (s *TopicService) DeleteTopic(ctx context.Context, id string) error {
	ctx, span := otel.Tracer("topic").Start(ctx, "TopicService.DeleteTopic")
	defer span.End()

	err := s.repo.Delete(ctx, id)
	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}
