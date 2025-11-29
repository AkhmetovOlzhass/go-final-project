package dto

type TopicResponse struct {
    ID          string  `json:"id"`
    Title       string  `json:"title"`
    Slug        string  `json:"slug"`
    ParentID    *string `json:"parentId,omitempty"`
    SchoolClass string  `json:"schoolClass"`
}