package service

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"learning-platform/internal/models"
)

type userRepoMockForUserService struct {
	mock.Mock
}

func (m *userRepoMockForUserService) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *userRepoMockForUserService) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if u, ok := args.Get(0).(*models.User); ok {
		return u, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *userRepoMockForUserService) FindByID(id string) (*models.User, error) {
	args := m.Called(id)
	if u, ok := args.Get(0).(*models.User); ok {
		return u, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *userRepoMockForUserService) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *userRepoMockForUserService) GetAll() ([]models.User, error) {
	args := m.Called()
	if users, ok := args.Get(0).([]models.User); ok {
		return users, args.Error(1)
	}
	return nil, args.Error(1)
}


func TestUserService_FindByID_Success(t *testing.T) {
	repo := new(userRepoMockForUserService)
	svc := NewUserService(repo, nil)

	idUUID := uuid.New()
	id := idUUID.String()

	expected := &models.User{
		ID:          idUUID,
		Email:       "user@example.com",
		DisplayName: "User",
	}

	repo.
		On("FindByID", id).
		Return(expected, nil).
		Once()

	user, err := svc.FindByID(id)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expected.ID, user.ID)
	assert.Equal(t, expected.Email, user.Email)
	assert.Equal(t, expected.DisplayName, user.DisplayName)

	repo.AssertExpectations(t)
}

func TestUserService_FindByID_RepoError(t *testing.T) {
	repo := new(userRepoMockForUserService)
	svc := NewUserService(repo, nil)

	id := uuid.New().String()
	expErr := errors.New("db error")

	repo.
		On("FindByID", id).
		Return((*models.User)(nil), expErr).
		Once()

	user, err := svc.FindByID(id)

	assert.Error(t, err)
	assert.Equal(t, expErr, err)
	assert.Nil(t, user)

	repo.AssertExpectations(t)
}

func TestUserService_FindByID_UserNotFound(t *testing.T) {
	repo := new(userRepoMockForUserService)
	svc := NewUserService(repo, nil)

	id := uuid.New().String()

	repo.
		On("FindByID", id).
		Return((*models.User)(nil), nil).
		Once()

	user, err := svc.FindByID(id)

	assert.Error(t, err)
	assert.EqualError(t, err, "user not found")
	assert.Nil(t, user)

	repo.AssertExpectations(t)
}

func TestUserService_Update_Success_AllFields(t *testing.T) {
	repo := new(userRepoMockForUserService)
	svc := NewUserService(repo, nil)

	idUUID := uuid.New()
	id := idUUID.String()

	oldAvatar := "http://old.example.com/avatar.png"

	existing := &models.User{
		ID:          idUUID,
		Email:       "old@example.com",
		DisplayName: "Old Name",
		AvatarURL:   &oldAvatar,
	}

	newEmail := "new@example.com"
	newName := "New Name"
	newAvatar := "http://new.example.com/avatar.png"

	repo.
		On("FindByID", id).
		Return(existing, nil).
		Once()

	repo.
		On("Update", mock.MatchedBy(func(u *models.User) bool {
			return u.ID == idUUID &&
				u.Email == newEmail &&
				u.DisplayName == newName &&
				u.AvatarURL != nil &&
				*u.AvatarURL == newAvatar
		})).
		Return(nil).
		Once()

	err := svc.Update(id, &newEmail, &newName, &newAvatar)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestUserService_Update_Partial_UpdateEmailOnly(t *testing.T) {
	repo := new(userRepoMockForUserService)
	svc := NewUserService(repo, nil)

	idUUID := uuid.New()
	id := idUUID.String()

	oldAvatar := "http://old.example.com/avatar.png"

	existing := &models.User{
		ID:          idUUID,
		Email:       "old@example.com",
		DisplayName: "Old Name",
		AvatarURL:   &oldAvatar,
	}

	newEmail := "updated@example.com"

	repo.
		On("FindByID", id).
		Return(existing, nil).
		Once()

	repo.
		On("Update", mock.MatchedBy(func(u *models.User) bool {
			return u.ID == idUUID &&
				u.Email == newEmail &&
				u.DisplayName == existing.DisplayName &&
				u.AvatarURL == existing.AvatarURL
		})).
		Return(nil).
		Once()

	err := svc.Update(id, &newEmail, nil, nil)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestUserService_Update_UserNotFound(t *testing.T) {
	repo := new(userRepoMockForUserService)
	svc := NewUserService(repo, nil)

	id := uuid.New().String()

	repo.
		On("FindByID", id).
		Return((*models.User)(nil), nil).
		Once()

	err := svc.Update(id, nil, nil, nil)

	assert.Error(t, err)
	assert.EqualError(t, err, "user not found")

	repo.AssertExpectations(t)
}

func TestUserService_Update_FindError(t *testing.T) {
	repo := new(userRepoMockForUserService)
	svc := NewUserService(repo, nil)

	id := uuid.New().String()
	expErr := errors.New("db error")

	repo.
		On("FindByID", id).
		Return((*models.User)(nil), expErr).
		Once()

	err := svc.Update(id, nil, nil, nil)

	assert.Error(t, err)
	assert.Equal(t, expErr, err)

	repo.AssertExpectations(t)
}

func TestUserService_Update_UpdateError(t *testing.T) {
	repo := new(userRepoMockForUserService)
	svc := NewUserService(repo, nil)

	idUUID := uuid.New()
	id := idUUID.String()

	existing := &models.User{
		ID:          idUUID,
		Email:       "old@example.com",
		DisplayName: "Old Name",
	}

	newEmail := "new@example.com"
	expErr := errors.New("update failed")

	repo.
		On("FindByID", id).
		Return(existing, nil).
		Once()

	repo.
		On("Update", mock.AnythingOfType("*models.User")).
		Return(expErr).
		Once()

	err := svc.Update(id, &newEmail, nil, nil)

	assert.Error(t, err)
	assert.Equal(t, expErr, err)

	repo.AssertExpectations(t)
}
