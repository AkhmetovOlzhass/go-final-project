package service

import (
	"errors"
	"learning-platform/internal/models"
	"learning-platform/internal/repository"
)

type UserService struct {
	users repository.IUserRepository
	s3    *S3Service
}

func NewUserService(users repository.IUserRepository, s3 *S3Service) *UserService {
	return &UserService{users: users, s3: s3}
}

func (s *UserService) GetAllUsers() ([]map[string]interface{}, error) {
    users, err := s.users.GetAll()
    if err != nil {
        return nil, err
    }

    formatted := make([]map[string]interface{}, len(users))

    for i, u := range users {
        formatted[i] = map[string]interface{}{
            "id":           u.ID,
            "email":        u.Email,
            "display_name": u.DisplayName,
            "role":         u.Role,
            "avatar_url":   u.AvatarURL,
            "created_at":   u.CreatedAt,
        }
    }

    return formatted, nil
}


func (s *UserService) FindByID(id string) (*models.User, error) {
	user, err := s.users.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *UserService) Update(id string, email, displayName *string, avatarURL *string) error {
	user, err := s.users.FindByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	if email != nil {
		user.Email = *email
	}
	if displayName != nil {
		user.DisplayName = *displayName
	}
	if avatarURL != nil {
		user.AvatarURL = avatarURL
	}

	return s.users.Update(user)
}
