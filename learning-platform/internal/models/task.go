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
    AnswerTypeFormula AnswerType  = "FORMULA"
)

type Task struct {
    ID              string       `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    Title           string       `gorm:"not null"`
    BodyMD          string       `gorm:"not null"`
    Difficulty      Difficulty   `gorm:"type:difficulty;not null"`
    Status          TaskStatus   `gorm:"type:task_status;not null"`
    CreatedAt       time.Time    `gorm:"autoCreateTime"`
    UpdatedAt       time.Time    `gorm:"autoUpdateTime"`

    TopicID         string       `gorm:"type:uuid;not null"`
    AuthorID        string       `gorm:"type:uuid;not null"`

    OfficialSolution string      
    CorrectAnswer     string      
    AnswerType        AnswerType  `gorm:"type:answer_type;not null"`
    ImageURL          string      

    Topic  *Topic `gorm:"foreignKey:TopicID"`
    Author *User  `gorm:"foreignKey:AuthorID"`
}
