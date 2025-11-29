package dto

type AuthTokensResponse struct {
    AccessToken  string `json:"accessToken"`
    RefreshToken string `json:"refreshToken"`
}

type MeResponse struct {
    ID          string  `json:"id"`
    Email       string  `json:"email"`
    DisplayName string  `json:"displayName"`
    Role        string  `json:"role"`
    AvatarURL   *string `json:"avatarUrl,omitempty"`
}

type RegisterResponse struct {
    Message string `json:"message"`
}