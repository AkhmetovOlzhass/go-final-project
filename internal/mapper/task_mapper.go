package mapper

import (
    "learning-platform/internal/models"
    "learning-platform/internal/dto"
)

func ToTaskResponse(t *models.Task) dto.TaskResponse {
    return dto.TaskResponse{
        ID:               t.ID,
        Title:            t.Title,
        BodyMD:           t.BodyMD,
        Difficulty:       string(t.Difficulty),
        Status:           string(t.Status),
        TopicID:          t.TopicID,
        AuthorID:         t.AuthorID,
        OfficialSolution: t.OfficialSolution,
        CorrectAnswer:    t.CorrectAnswer,
        AnswerType:       string(t.AnswerType),
        ImageURL:         t.ImageURL,
        CreatedAt:        t.CreatedAt.Format("2006-01-02T15:04:05Z"),
        UpdatedAt:        t.UpdatedAt.Format("2006-01-02T15:04:05Z"),
    }
}

func ToTaskList(tasks []models.Task) []dto.TaskResponse {
    res := make([]dto.TaskResponse, len(tasks))
    for i, t := range tasks {
        res[i] = ToTaskResponse(&t)
    }
    return res
}
