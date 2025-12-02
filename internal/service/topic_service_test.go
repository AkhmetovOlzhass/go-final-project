package service

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"learning-platform/internal/models"
)

type fakeTopicRepo struct {
	topics       []models.Topic
	findAllCalls int
}

func newFakeTopicRepo() *fakeTopicRepo {
	return &fakeTopicRepo{
		topics: make([]models.Topic, 0),
	}
}

func (f *fakeTopicRepo) Create(ctx context.Context, topic *models.Topic) error {
	if topic.ID == "" {
		topic.ID = "generated-" + time.Now().Format("150405.000")
	}
	f.topics = append(f.topics, *topic)
	return nil
}

func (f *fakeTopicRepo) FindAll(ctx context.Context) ([]models.Topic, error) {
	f.findAllCalls++
	return f.topics, nil
}

func (f *fakeTopicRepo) FindByID(ctx context.Context, id string) (*models.Topic, error) {
	for i := range f.topics {
		if f.topics[i].ID == id {
			return &f.topics[i], nil
		}
	}
	return nil, nil
}

func (f *fakeTopicRepo) Update(ctx context.Context, topic *models.Topic) error {
	for i := range f.topics {
		if f.topics[i].ID == topic.ID {
			f.topics[i] = *topic
			return nil
		}
	}
	return nil
}

func (f *fakeTopicRepo) Delete(ctx context.Context, id string) error {
	out := make([]models.Topic, 0, len(f.topics))
	for _, t := range f.topics {
		if t.ID != id {
			out = append(out, t)
		}
	}
	f.topics = out
	return nil
}

func TestTopicService_GetAllTopics_UsesCache(t *testing.T) {
	ctx := context.Background()
	redisClient := newTestRedis(t)

	repo := newFakeTopicRepo()
	repo.topics = []models.Topic{
		{ID: "topic-1", Title: "Topic 1"},
		{ID: "topic-2", Title: "Topic 2"},
	}

	svc := NewTopicService(repo, redisClient)

	topics1, err := svc.GetAllTopics(ctx)
	require.NoError(t, err)
	require.Len(t, topics1, 2)
	assert.Equal(t, 1, repo.findAllCalls, "первый вызов должен ходить в репозиторий")

	raw, err := redisClient.Get(ctx, "topics:all").Result()
	require.NoError(t, err)

	var cached []models.Topic
	require.NoError(t, json.Unmarshal([]byte(raw), &cached))
	require.Len(t, cached, 2)
	topics2, err := svc.GetAllTopics(ctx)
	require.NoError(t, err)
	require.Len(t, topics2, 2)
	assert.Equal(t, 1, repo.findAllCalls, "второй вызов должен брать из кеша, repo.FindAll не вызывается")
}

func TestTopicService_Create_Update_Delete_InvalidatesCache(t *testing.T) {
	ctx := context.Background()
	redisClient := newTestRedis(t)
	repo := newFakeTopicRepo()
	svc := NewTopicService(repo, redisClient)

	topic := &models.Topic{
		ID:    "topic-1",
		Title: "Initial",
	}

	err := svc.CreateTopic(ctx, topic)
	require.NoError(t, err)

	_, err = redisClient.Get(ctx, "topics:all").Result()
	assert.Error(t, err, "после CreateTopic кеш topics:all должен быть очищен")

	data, _ := json.Marshal(repo.topics)
	require.NoError(t, redisClient.Set(ctx, "topics:all", data, 10*time.Minute).Err())

	topic.Title = "Updated"
	err = svc.UpdateTopic(ctx, topic)
	require.NoError(t, err)

	_, err = redisClient.Get(ctx, "topics:all").Result()
	assert.Error(t, err, "после UpdateTopic кеш topics:all должен быть очищен")

	require.NoError(t, redisClient.Set(ctx, "topics:all", data, 10*time.Minute).Err())

	err = svc.DeleteTopic(ctx, topic.ID)
	require.NoError(t, err)

	_, err = redisClient.Get(ctx, "topics:all").Result()
	assert.Error(t, err, "после DeleteTopic кеш topics:all должен быть очищен")
}
