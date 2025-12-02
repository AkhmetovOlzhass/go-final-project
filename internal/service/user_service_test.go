package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"learning-platform/internal/models"
)

func strPtr(s string) *string { return &s }

func TestUserService_GetAllUsers(t *testing.T) {
	ctx := context.Background()

	repo := newFakeUserRepo()

	u1 := &models.User{
		ID:    uuid.New(),
		Email: "a@test.com",
	}
	u2 := &models.User{
		ID:    uuid.New(),
		Email: "b@test.com",
	}

	require.NoError(t, repo.Create(ctx, u1))
	require.NoError(t, repo.Create(ctx, u2))

	svc := NewUserService(repo)

	users, err := svc.GetAllUsers(ctx)
	require.NoError(t, err)
	require.Len(t, users, 2)
}

func TestUserService_FindByID_Success(t *testing.T) {
	ctx := context.Background()
	repo := newFakeUserRepo()

	u := &models.User{
		ID:    uuid.New(),
		Email: "test@example.com",
	}

	require.NoError(t, repo.Create(ctx, u))

	svc := NewUserService(repo)

	found, err := svc.FindByID(ctx, u.ID.String())
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, "test@example.com", found.Email)
}

func TestUserService_FindByID_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := newFakeUserRepo()
	svc := NewUserService(repo)

	found, err := svc.FindByID(ctx, uuid.New().String())
	require.NoError(t, err)
	assert.Nil(t, found)
}

func TestUserService_Update_Success(t *testing.T) {
	ctx := context.Background()
	repo := newFakeUserRepo()

	u := &models.User{
		ID:          uuid.New(),
		Email:       "old@mail.com",
		DisplayName: "Old Name",
		AvatarURL:   strPtr("old.png"), 
	}

	require.NoError(t, repo.Create(ctx, u))

	svc := NewUserService(repo)

	newEmail := "new@mail.com"
	newName := "New Name"
	newAvatar := "new.png"

	updated, err := svc.Update(ctx, u.ID.String(), &newEmail, &newName, &newAvatar)
	require.NoError(t, err)
	require.NotNil(t, updated)

	assert.Equal(t, "new@mail.com", updated.Email)
	assert.Equal(t, "New Name", updated.DisplayName)
	require.NotNil(t, updated.AvatarURL)
	assert.Equal(t, "new.png", *updated.AvatarURL)

}
