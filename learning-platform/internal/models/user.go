package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string
type UserStatus string

const (
	UserRoleAdmin   UserRole = "Admin"
	UserRoleTeacher UserRole = "Teacher"
	UserRoleStudent UserRole = "Student"
)

const (
	UserStatusPending UserStatus = "PENDING"
	UserStatusActive  UserStatus = "ACTIVE"
)

type User struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	DisplayName  string         `gorm:"type:varchar(255)"`
	Email        string         `gorm:"uniqueIndex;not null"`
	PasswordHash string         `gorm:"type:text;not null"`
	AvatarURL    *string        `gorm:"type:text" json:"avatar,omitempty"`
	Role         UserRole       `gorm:"type:userRole;default:'Student';not null"`
	Status       UserStatus     `gorm:"type:userStatus;default:'PENDING';not null"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
