package dto

import "time"

type UserResponse struct {
	ID          string  `json:"id"`
	Email       string  `json:"email"`
	DisplayName string  `json:"displayName"`
	Role        string  `json:"role"`
	AvatarURL   *string `json:"avatarUrl"`
	IsBanned    string  `json:"isBanned"`
}

type BanProfileResponse struct {
	UserID       string     `json:"userId"`
	IsBanned     string     `json:"isBanned"`
	BannedAt     *time.Time `json:"bannedAt"`
	BannedUntil  *time.Time `json:"bannedUntil,omitempty"`
	BannedReason *string    `json:"bannedReason,omitempty"`
}
