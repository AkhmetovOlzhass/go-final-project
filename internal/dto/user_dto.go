package dto

import "time"

type UpdateProfileRequest struct {
    Email       *string `form:"email" binding:"required,email"`
    DisplayName *string `form:"displayName" binding:"required"`
}

type BanProfileRequest struct {
    BannedReason *string     `json:"bannedReason,omitempty"`              
    BannedUntil  *time.Time `json:"bannedUntil,omitempty"`
}
