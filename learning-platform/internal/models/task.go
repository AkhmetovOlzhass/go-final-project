package models

import (
	"time"
)

type Difficulty string
type TaskStatus string
type SchoolClass string
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
	SchoolClassSeven  SchoolClass = "SEVEN"
	SchoolClassEight  SchoolClass = "EIGHT"
	SchoolClassNine   SchoolClass = "NINE"
	SchoolClassTen    SchoolClass = "TEN"
	SchoolClassEleven SchoolClass = "ELEVEN"
)

const (
	AnswerTypeText   AnswerType = "TEXT"
	AnswerTypeNumber AnswerType = "NUMBER"
	AnswerTypeFormula AnswerType = "FORMULA"
)

type Topic struct {
	ID         string     `json:"id"`
	Title      string     `json:"title"`
	Slug       string     `json:"slug"`
	ParentID   *string    `json:"parent_id,omitempty"`
	SchoolClass SchoolClass `json:"school_class"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type Task struct {
	ID             string      `json:"id"`
	Title          string      `json:"title"`
	BodyMD         string      `json:"body_md"`
	Difficulty     Difficulty  `json:"difficulty"`
	Status         TaskStatus  `json:"status"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
	
	TopicID        string      `json:"topic_id"`
	AuthorID       string      `json:"author_id"`
	
	OfficialSolution string    `json:"official_solution,omitempty"`
	CorrectAnswer    string    `json:"correct_answer,omitempty"`
	AnswerType      AnswerType `json:"answer_type"`
	ImageURL        string     `json:"image_url,omitempty"`
	
	// Relations
	Topic  *Topic `json:"topic,omitempty"`
	Author *User  `json:"author,omitempty"`
}