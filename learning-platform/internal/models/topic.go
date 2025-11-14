package models

import "time"

type Topic struct {
    ID          string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    Title       string    `gorm:"not null"`
    Slug        string    `gorm:"unique;not null"`
    ParentID    *string   `gorm:"type:uuid"`
    SchoolClass string    `gorm:"type:school_class;not null"`
    CreatedAt   time.Time `gorm:"autoCreateTime"`
    UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
