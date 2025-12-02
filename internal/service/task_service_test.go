package service

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"learning-platform/internal/models"
)

type fakeTaskRepo struct {
	all       []models.Task
	drafts    []models.Task
	byID      map[string]*models.Task
	byTopic   map[string][]models.Task
	byAuthor  map[string][]models.Task
	statusSet map[string]models.TaskStatus

	getAllCalls int
}

func newFakeTaskRepo() *fakeTaskRepo {
	return &fakeTaskRepo{
		byID:      make(map[string]*models.Task),
		byTopic:   make(map[string][]models.Task),
		byAuthor:  make(map[string][]models.Task),
		statusSet: make(map[string]models.TaskStatus),
	}
}

func (f *fakeTaskRepo) GetAll(ctx context.Context) ([]models.Task, error) {
	f.getAllCalls++
	return f.all, nil
}

func (f *fakeTaskRepo) GetDraft(ctx context.Context) ([]models.Task, error) {
	return f.drafts, nil
}

func (f *fakeTaskRepo) UpdateStatus(ctx context.Context, id string, status models.TaskStatus) error {
	f.statusSet[id] = status
	if t, ok := f.byID[id]; ok && t != nil {
		t.Status = status
	}
	return nil
}

func (f *fakeTaskRepo) Create(ctx context.Context, task *models.Task) error {
	if task.ID == "" {
		task.ID = "generated-" + time.Now().Format("150405.000")
	}
	f.all = append(f.all, *task)
	f.byID[task.ID] = task
	f.byTopic[task.TopicID] = append(f.byTopic[task.TopicID], *task)
	f.byAuthor[task.AuthorID] = append(f.byAuthor[task.AuthorID], *task)
	return nil
}

func (f *fakeTaskRepo) GetByID(ctx context.Context, id string) (*models.Task, error) {
	return f.byID[id], nil
}

func (f *fakeTaskRepo) GetByTopic(ctx context.Context, topicID string) ([]models.Task, error) {
	return f.byTopic[topicID], nil
}

func (f *fakeTaskRepo) Update(ctx context.Context, task *models.Task) error {
	if existing, ok := f.byID[task.ID]; ok && existing != nil {
		*existing = *task
	}
	return nil
}

func (f *fakeTaskRepo) Delete(ctx context.Context, id string) error {
	delete(f.byID, id)
	return nil
}

func (f *fakeTaskRepo) GetByAuthor(ctx context.Context, authorID string) ([]models.Task, error) {
	return f.byAuthor[authorID], nil
}

func newTestRedis(t *testing.T) *redis.Client {
	mr, err := miniredis.Run()
	require.NoError(t, err)
	t.Cleanup(mr.Close)

	return redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
}

func TestTaskService_GetAllTasks_UsesCache(t *testing.T) {
	ctx := context.Background()
	rdb := newTestRedis(t)

	repo := newFakeTaskRepo()
	task1 := models.Task{
		ID:       "task-1",
		Title:    "Task 1",
		Status:   models.TaskStatusDraft,
		TopicID:  "topic-1",
		AuthorID: "author-1",
	}
	task2 := models.Task{
		ID:       "task-2",
		Title:    "Task 2",
		Status:   models.TaskStatusDraft,
		TopicID:  "topic-1",
		AuthorID: "author-1",
	}
	repo.all = []models.Task{task1, task2}

	svc := NewTaskService(repo, rdb)

	tasks1, err := svc.GetAllTasks(ctx)
	require.NoError(t, err)
	require.Len(t, tasks1, 2)
	assert.Equal(t, 1, repo.getAllCalls)

	raw, err := rdb.Get(ctx, "tasks:all").Result()
	require.NoError(t, err)

	var cached []models.Task
	require.NoError(t, json.Unmarshal([]byte(raw), &cached))
	require.Len(t, cached, 2)

	tasks2, err := svc.GetAllTasks(ctx)
	require.NoError(t, err)
	require.Len(t, tasks2, 2)
	assert.Equal(t, 1, repo.getAllCalls, "repo.GetAll не должен вызываться повторно при хите в кеш")
}

func TestTaskService_PublishTask_InvalidatesCache(t *testing.T) {
	ctx := context.Background()
	rdb := newTestRedis(t)

	repo := newFakeTaskRepo()
	taskID := "task-1"
	task := &models.Task{
		ID:       taskID,
		Title:    "Draft task",
		Status:   models.TaskStatusDraft,
		TopicID:  "topic-1",
		AuthorID: "author-1",
	}
	repo.byID[taskID] = task

	svc := NewTaskService(repo, rdb)

	data, _ := json.Marshal([]models.Task{*task})
	require.NoError(t, rdb.Set(ctx, "tasks:all", data, 10*time.Minute).Err())

	published, err := svc.PublishTask(ctx, taskID)
	require.NoError(t, err)

	assert.Equal(t, models.TaskStatusPublished, repo.statusSet[taskID])
	assert.Equal(t, models.TaskStatusPublished, published.Status)

	_, err = rdb.Get(ctx, "tasks:all").Result()
	assert.Error(t, err, "после PublishTask кеш должен быть удалён")
}

func TestTaskService_Create_Update_Delete_InvalidatesCache(t *testing.T) {
	ctx := context.Background()
	rdb := newTestRedis(t)
	repo := newFakeTaskRepo()
	svc := NewTaskService(repo, rdb)

	task := &models.Task{
		ID:       "task-1",
		Title:    "New task",
		Status:   models.TaskStatusDraft,
		TopicID:  "topic-1",
		AuthorID: "author-1",
	}

	require.NoError(t, svc.CreateTask(ctx, task))
	_, err := rdb.Get(ctx, "tasks:all").Result()
	assert.Error(t, err, "после CreateTask кеш должен быть удалён")

	data, _ := json.Marshal(repo.all)
	require.NoError(t, rdb.Set(ctx, "tasks:all", data, 10*time.Minute).Err())

	require.NoError(t, svc.UpdateTask(ctx, task))
	_, err = rdb.Get(ctx, "tasks:all").Result()
	assert.Error(t, err, "после UpdateTask кеш должен быть удалён")

	require.NoError(t, rdb.Set(ctx, "tasks:all", data, 10*time.Minute).Err())

	require.NoError(t, svc.DeleteTask(ctx, task.ID))
	_, err = rdb.Get(ctx, "tasks:all").Result()
	assert.Error(t, err, "после DeleteTask кеш должен быть удалён")
}
