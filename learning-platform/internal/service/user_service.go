package service

import (
    "context"
	"errors"
    "strings"

	"learning-platform/internal/models"
	"learning-platform/internal/repository"
    "go.opentelemetry.io/otel"
)

type UserService struct {
	users repository.IUserRepository
}

func NewUserService(users repository.IUserRepository) *UserService {
	return &UserService{users: users}
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	ctx, span := otel.Tracer("user").Start(ctx, "UserService.GetAllUsers")
	defer span.End()

	users, err := s.users.GetAll(ctx)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return users, nil
}

func (s *UserService) FindByID(ctx context.Context, id string) (*models.User, error) {
	ctx, span := otel.Tracer("user").Start(ctx, "UserService.FindByID")
	defer span.End()

	user, err := s.users.FindByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return user, nil
}

func (s *UserService) Update(ctx context.Context, id string, email, displayName *string, avatarURL *string) (*models.User, error) {
	ctx, span := otel.Tracer("user").Start(ctx, "UserService.Update")
	defer span.End()

	user, err := s.users.FindByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	if user == nil {
		err := errors.New("user not found")
		span.RecordError(err)
		return nil, err
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

	_, spanUpdate := otel.Tracer("user").Start(ctx, "DB.UpdateUser")
	err = s.users.Update(ctx, id, updates)
	spanUpdate.End()

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			e := errors.New("email already taken")
			span.RecordError(e)
			return nil, e
		}
		span.RecordError(err)
		return nil, err
	}

	updated, err := s.users.FindByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return updated, nil
}