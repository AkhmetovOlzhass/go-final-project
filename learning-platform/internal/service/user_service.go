package service

import (
	"errors"

	"learning-platform/internal/models"
	"learning-platform/internal/repository"
)

type UserService struct {
	users repository.IUserRepository
}

func NewUserService(users repository.IUserRepository) *UserService {
	return &UserService{users: users}
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.users.GetAll()
}

func (s *UserService) FindByID(id string) (*models.User, error) {
	return s.users.FindByID(id)
}

func (s *UserService) Update(id string, email, displayName *string, avatarURL *string) (*models.User, error) {
    user, err := s.users.FindByID(id)
    if err != nil {
        return nil, err
    }
    if user == nil {
        return nil, errors.New("user not found")
    }

    updates := map[string]interface{}{}

    if email != nil {
        updates["email"] = *email
    }
    if displayName != nil {
        updates["displayName"] = *displayName
    }
    if avatarURL != nil {
        updates["avatarUrl"] = *avatarURL
    }

    if len(updates) == 0 {
        return user, nil
    }

    if err := s.users.Update(id, updates); err != nil {
        return nil, err
    }

    return s.users.FindByID(id)
}

