package mapper

import (
    "learning-platform/internal/dto"
    "learning-platform/internal/models"
)

func ToUserResponse(u *models.User) dto.UserResponse {
	id := u.ID.String() 
	role := string(u.Role)
    return dto.UserResponse{
        ID:          id,
        Email:       u.Email,
        DisplayName: u.DisplayName,
        Role:        role,
        AvatarURL:   u.AvatarURL,
    }
}

func ToUserList(users []models.User) []dto.UserResponse {
    result := make([]dto.UserResponse, 0, len(users))
    for _, u := range users {
        result = append(result, ToUserResponse(&u))
    }
    return result
}

func ToBanProfileResponse(u *models.User) dto.BanProfilePayload {
  return dto.BanProfilePayload{
    UserID:      u.ID.String(),
    IsBanned:      string(u.IsBanned),
    BannedAt:    u.BannedAt,
    BannedUntil: u.BannedUntil,
    BannedReason: u.BanReason,
  }
}