package service

import (
	"context"

	"learning-platform/internal/models"
	"learning-platform/internal/repository"
	"go.opentelemetry.io/otel"
	"github.com/redis/go-redis/v9"
	"encoding/json"
	"time"
)

type TopicService struct {
	repo repository.ITopicRepository
	redis *redis.Client
}

func NewTopicService(repo repository.ITopicRepository, rdb *redis.Client) *TopicService {
  return &TopicService{
    repo:  repo,
    redis: rdb,
  }
}

func (s *TopicService) CreateTopic(ctx context.Context, topic *models.Topic) error {
	ctx, span := otel.Tracer("topic").Start(ctx, "TopicService.CreateTopic")
	defer span.End()

	err := s.repo.Create(ctx, topic)
	if err != nil {
		span.RecordError(err)
		return err
	}

	s.redis.Del(context.Background(), "topics:all")

	return nil
}

func (s *TopicService) GetAllTopics(ctx context.Context) ([]models.Topic, error) {
	ctx, span := otel.Tracer("topic").Start(ctx, "TopicService.GetAllTopics")
	defer span.End()

	cacheKey := "topics:all"
	if cached, err := s.redis.Get(ctx, cacheKey).Result(); err == nil {
		var topics []models.Topic
		if err := json.Unmarshal([]byte(cached), &topics); err == nil {
		return topics, nil
		}
	}

	topics, err := s.repo.FindAll(ctx)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	data, _ := json.Marshal(topics)
  	s.redis.Set(ctx, cacheKey, data, 10*time.Minute)

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
	s.redis.Del(context.Background(), "topics:all")
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
	s.redis.Del(context.Background(), "topics:all")
	return nil
}
