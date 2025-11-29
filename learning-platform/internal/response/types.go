package response

type ErrorResponse struct {
	Success bool         `json:"success" example:"false"` 
	Error   ErrorMessage `json:"error"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}

type SuccessWrapper struct {
    Success bool        `json:"success" example:"true"`
    Data    interface{} `json:"data"`
}