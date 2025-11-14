package models

import (
	"time"
)

type Difficulty string
type TaskStatus string
type AnswerType string

const (
	DifficultyEasy    Difficulty = "EASY"
	DifficultyMedium  Difficulty = "MEDIUM"
	DifficultyHard    Difficulty = "HARD"
	DifficultyExtreme Difficulty = "EXTREME"
)

const (
	TaskStatusDraft     TaskStatus = "DRAFT"
	TaskStatusPublished TaskStatus = "PUBLISHED"
	TaskStatusArchived  TaskStatus = "ARCHIVED"
)

const (
	AnswerTypeText    AnswerType = "TEXT"
	AnswerTypeNumber  AnswerType = "NUMBER"
	AnswerTypeFormula AnswerType = "FORMULA"
)

type Task struct {
	ID             string        `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Title          string        `gorm:"not null" json:"title"`
	BodyMD         string        `gorm:"not null" json:"body_md"`
	Difficulty     Difficulty    `gorm:"type:difficulty;not null" json:"difficulty"`
	Status         TaskStatus    `gorm:"type:task_status;not null" json:"status"`
	CreatedAt      time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time     `gorm:"autoUpdateTime" json:"updated_at"`

	TopicID        string        `gorm:"type:uuid;not null" json:"topic_id"`
	AuthorID       string        `gorm:"type:uuid;not null" json:"author_id"`

	OfficialSolution string      `json:"official_solution,omitempty"`
	CorrectAnswer    string      `json:"correct_answer,omitempty"`
	AnswerType       AnswerType  `gorm:"type:answer_type;not null" json:"answer_type"`
	ImageURL         string      `json:"image_url,omitempty"`

	Topic  *Topic `gorm:"foreignKey:TopicID" json:"topic,omitempty"`
	Author *User  `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
}
