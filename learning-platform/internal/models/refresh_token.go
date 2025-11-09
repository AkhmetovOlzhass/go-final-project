package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshToken struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID     uuid.UUID `gorm:"not null;index"`
	TokenHash  string    `gorm:"type:text;uniqueIndex;not null"`
	Revoked    bool      `gorm:"default:false"`
	ExpiresAt  time.Time `gorm:"not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}

func (r *RefreshToken) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ExpiresAt.IsZero() {
		r.ExpiresAt = time.Now().Add(30 * 24 * time.Hour)
	}
	return
}
