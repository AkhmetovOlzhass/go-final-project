package dto

type UserResponse struct {
    ID          string  `json:"id"`
    Email       string  `json:"email"`
    DisplayName string  `json:"displayName"`
    Role        string  `json:"role"`
    AvatarURL   *string `json:"avatarUrl"`
}