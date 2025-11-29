package mapper

import (
    "learning-platform/internal/models"
    "learning-platform/internal/dto"
)

func ToTopicResponse(t *models.Topic) dto.TopicResponse {
    return dto.TopicResponse{
        ID:          t.ID,
        Title:       t.Title,
        Slug:        t.Slug,
        ParentID:    t.ParentID,
        SchoolClass: t.SchoolClass,
    }
}

func ToTopicList(topics []models.Topic) []dto.TopicResponse {
    res := make([]dto.TopicResponse, len(topics))
    for i, t := range topics {
        res[i] = ToTopicResponse(&t)
    }
    return res
}
