package mapper

import (
    "learning-platform/internal/models"
    "learning-platform/internal/dto"
)

func ToMeResponse(u *models.User) dto.MeResponse {
    return dto.MeResponse{
        ID:          u.ID.String(),
        Email:       u.Email,
        DisplayName: u.DisplayName,
        Role:        string(u.Role),
        AvatarURL:   u.AvatarURL,
    }
}