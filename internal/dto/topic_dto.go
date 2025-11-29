package dto

type CreateTopicRequest struct {
    Title       string  `json:"title" binding:"required"`
    Slug        string  `json:"slug" binding:"required"`
    ParentID    *string `json:"parentId"`
    SchoolClass string  `json:"schoolClass" binding:"required"`
}

type UpdateTopicRequest struct {
    Title       string  `json:"title" binding:"required"`
    Slug        string  `json:"slug" binding:"required"`
    ParentID    *string `json:"parentId"`
    SchoolClass string  `json:"schoolClass" binding:"required"`
}
