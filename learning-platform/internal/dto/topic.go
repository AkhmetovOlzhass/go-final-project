package dto

type CreateTopicRequest struct {
    Title       string  `json:"title" binding:"required"`
    Slug        string  `json:"slug" binding:"required"`
    ParentID    *string `json:"parent_id"`
    SchoolClass string  `json:"school_class" binding:"required"`
}

type UpdateTopicRequest struct {
    Title       string  `json:"title" binding:"required"`
    Slug        string  `json:"slug" binding:"required"`
    ParentID    *string `json:"parent_id"`
    SchoolClass string  `json:"school_class" binding:"required"`
}
