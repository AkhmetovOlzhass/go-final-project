package dto

type TaskResponse struct {
    ID               string  `json:"id"`
    Title            string  `json:"title"`
    BodyMD           string  `json:"bodyMd"`
    Difficulty       string  `json:"difficulty"`
    Status           string  `json:"status"`
    TopicID          string  `json:"topicId"`
    AuthorID         string  `json:"authorId"`
    AnswerType       string  `json:"answerType"`
    ImageURL         string  `json:"imageUrl,omitempty"`
    OfficialSolution string  `json:"officialSolution,omitempty"`
    CorrectAnswer    string  `json:"correctAnswer,omitempty"`
    CreatedAt        string  `json:"createdAt"`
    UpdatedAt        string  `json:"updatedAt"`
}
