package dto

type UpdateProfileRequest struct {
    Email       string `form:"email" binding:"required,email"`
    DisplayName string `form:"displayName" binding:"required"`
}