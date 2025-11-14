package models

import "time"

type Topic struct {
    ID          string     `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    Title       string     `json:"title" gorm:"not null"`
    Slug        string     `json:"slug" gorm:"unique;not null"`
    ParentID    *string    `json:"parent_id" gorm:"type:uuid"`
    SchoolClass string     `json:"school_class" gorm:"type:school_class;not null"`
    CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}
