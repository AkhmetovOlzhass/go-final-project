package dto

import "learning-platform/internal/models"

type CreateTaskRequest struct {
    Title            string            `form:"title" binding:"required"`
    BodyMD           string            `form:"bodyMd" binding:"required"`
    Difficulty       models.Difficulty `form:"difficulty" binding:"required"`
    Status           models.TaskStatus `form:"status" binding:"required"`
    TopicID          string            `form:"topicId" binding:"required"`
    OfficialSolution string            `form:"officialSolution"`
    CorrectAnswer    string            `form:"correctAnswer"`
    AnswerType       models.AnswerType `form:"answerType" binding:"required"`
}

type UpdateTaskRequest struct {
    Title            string            `form:"title" binding:"required"`
    BodyMD           string            `form:"bodyMd" binding:"required"`
    Difficulty       models.Difficulty `form:"difficulty" binding:"required"`
    Status           models.TaskStatus `form:"status" binding:"required"`
    TopicID          string            `form:"topicId" binding:"required"`
    OfficialSolution string            `form:"officialSolution"`
    CorrectAnswer    string            `form:"correctAnswer"`
    AnswerType       models.AnswerType `form:"answerType" binding:"required"`
}


// TaskSubmitRequest - запрос на проверку ответа
type TaskSubmitRequest struct {
    Answer string `json:"answer" binding:"required"`
}

// TaskSubmitResponse - ответ на проверку ответа
type TaskSubmitResponse struct {
    Correct bool `json:"correct"`
}