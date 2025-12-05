package handler

import (
	"net/http"

	"learning-platform/internal/dto"
	"learning-platform/internal/models"
	"learning-platform/internal/response"
	"learning-platform/internal/mapper"
	"learning-platform/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type TaskHandler struct {
	taskService *service.TaskService
	s3    *service.S3Service
}

func NewTaskHandler(taskService *service.TaskService, s3 *service.S3Service) *TaskHandler {
	return &TaskHandler{taskService: taskService,  s3: s3}
}

// GetAllTasks godoc
// @Summary Get all published tasks
// @Tags tasks
// @Description Returns list of all published tasks
// @Produce json
// @Success 200 {object} response.SuccessWrapper{data=[]dto.TaskResponse}
// @Failure 500 {object} response.ErrorResponse
// @Router /tasks [get]
func (h *TaskHandler) GetAllTasks(c *gin.Context) {
	ctx := c.Request.Context()

	tasks, err := h.taskService.GetAllTasks(ctx)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to fetch tasks")
		return
	}

	response.Success(c, mapper.ToTaskList(tasks))
}

// GetDraftTasks godoc
// @Summary Get all draft tasks
// @Tags tasks
// @Description Returns tasks with status DRAFT
// @Produce json
// @Success 200 {object} response.SuccessWrapper{data=[]dto.TaskResponse}
// @Failure 500 {object} response.ErrorResponse
// @Security BearerAuth
// @Router /tasks/drafts [get]
func (h *TaskHandler) GetDraftTasks(c *gin.Context) {
	ctx := c.Request.Context()

	tasks, err := h.taskService.GetDraftTasks(ctx)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to fetch draft tasks")
		return
	}

	response.Success(c, mapper.ToTaskList(tasks))
}

// PublishTask godoc
// @Summary Publish a task
// @Tags tasks
// @Description Changes task status to PUBLISHED
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} response.SuccessWrapper{data=dto.TaskResponse}
// @Failure 500 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Security BearerAuth
// @Router /tasks/{id}/publish [patch]
func (h *TaskHandler) PublishTask(c *gin.Context) {
	ctx := c.Request.Context()

    id := c.Param("id")
	t, err := h.taskService.PublishTask(ctx, id);
    if  err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to publish task")
        return
    }

    response.Success(c, mapper.ToTaskResponse(t))
}


// CreateTask godoc
// @Summary Create a new task
// @Tags tasks
// @Description Create a task (multipart form with optional image)
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "Title"
// @Param body_md formData string true "Markdown body"
// @Param difficulty formData string true "Difficulty"
// @Param status formData string true "Task status"
// @Param topicId formData string true "Topic ID"
// @Param officialSolution formData string true "Solution"
// @Param correctAnswer formData string true "Correct answer"
// @Param answerType formData string true "Answer type"
// @Param imageUrl formData file false "Task image"
// @Success 201 {object} response.SuccessWrapper{data=dto.TaskResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Security BearerAuth
// @Router /tasks [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.CreateTaskRequest

	if err := c.ShouldBindWith(&req, binding.FormMultipart); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := c.GetString("userId")

	var imageURL string
	file, err := c.FormFile("imageUrl") 
	if err == nil && file != nil {
		url, uploadErr := h.s3.UploadFile(ctx, file)
		if uploadErr != nil {
			response.Error(c, http.StatusInternalServerError, "failed to upload image")
			return
		}
		imageURL = url
	}

	task := &models.Task{
		Title:            req.Title,
		BodyMD:           req.BodyMD,
		Difficulty:       req.Difficulty,
		Status:           req.Status,
		TopicID:          req.TopicID,
		AuthorID:         userID,
		OfficialSolution: req.OfficialSolution,
		CorrectAnswer:    req.CorrectAnswer,
		AnswerType:       req.AnswerType,
		ImageURL:         imageURL,
	}

	if err := h.taskService.CreateTask(ctx, task); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create task")
		return
	}

	response.SuccessWithStatus(c, http.StatusCreated, mapper.ToTaskResponse(task))
}

// GetTask godoc
// @Summary Get task by ID
// @Tags tasks
// @Description Returns a single task
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} response.SuccessWrapper{data=dto.TaskResponse}
// @Failure 404 {object} response.ErrorResponse
// @Router /tasks/{id} [get]
func (h *TaskHandler) GetTask(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	task, err := h.taskService.GetTaskById(ctx, id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Task not found")
		return
	}

	response.Success(c, mapper.ToTaskResponse(task))
}

// GetTasksByTopic godoc
// @Summary Get tasks by topic
// @Tags tasks
// @Description Returns all tasks belonging to topic
// @Produce json
// @Param topicId path string true "Topic ID"
// @Success 200 {object} response.SuccessWrapper{data=[]dto.TaskResponse}
// @Failure 500 {object} response.ErrorResponse
// @Router /topics/{topicId}/tasks [get]
func (h *TaskHandler) GetTasksByTopic(c *gin.Context) {
	ctx := c.Request.Context()

	topicID := c.Param("topicId")

	tasks, err := h.taskService.GetTasksByTopic(ctx, topicID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch tasks")
		return
	}

	response.Success(c, mapper.ToTaskList(tasks))
}

// UpdateTask godoc
// @Summary Update a task
// @Tags tasks
// @Description Update task metadata + (optional) image
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Task ID"
// @Param title formData string true "Title"
// @Param body_md formData string true "Body"
// @Param difficulty formData string true "Difficulty"
// @Param status formData string true "Status"
// @Param topicId formData string true "Topic ID"
// @Param officialSolution formData string true "Solution"
// @Param correctAnswer formData string true "Correct answer"
// @Param answerType formData string true "Answer type"
// @Param imageUrl formData file false "Task image"
// @Success 200 {object} response.SuccessWrapper{data=dto.TaskResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Security BearerAuth
// @Router /tasks/{id} [put]
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	ctx := c.Request.Context()

    id := c.Param("id")

    var req dto.UpdateTaskRequest
    if err := c.ShouldBindWith(&req, binding.FormMultipart); err != nil {
        response.Error(c, http.StatusBadRequest, err.Error())
        return
    }

	userID := c.GetString("userId")

	existing, err := h.taskService.GetTaskById(ctx, id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Task not found")
		return
	}
	
	imageURL := existing.ImageURL

    file, err := c.FormFile("imageUrl")
    if err == nil && file != nil {
        url, uploadErr := h.s3.UploadFile(ctx, file)
        if uploadErr != nil {
            response.Error(c, http.StatusInternalServerError, "failed to upload image")
            return
        }
        imageURL = url
    }

    updated := &models.Task{
        ID:               id,
        Title:            req.Title,
        BodyMD:           req.BodyMD,
        Difficulty:       req.Difficulty,
        Status:           req.Status,
        TopicID:          req.TopicID,
        AuthorID:         userID,
        OfficialSolution: req.OfficialSolution,
        CorrectAnswer:    req.CorrectAnswer,
        AnswerType:       req.AnswerType,
        ImageURL:         imageURL,
    }

    if err := h.taskService.UpdateTask(ctx, updated); err != nil {
        response.Error(c, http.StatusInternalServerError, "Failed to update task")
        return
    }

    finalTask, err := h.taskService.GetTaskById(ctx, id)

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to fetch updated task")
		return
	}
	
    response.Success(c, mapper.ToTaskResponse(finalTask))
}


// DeleteTask godoc
// @Summary Delete a task
// @Tags tasks
// @Description Remove task by ID
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} response.SuccessWrapper{data=string}
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Security BearerAuth
// @Router /tasks/{id} [delete]
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	_, err := h.taskService.GetTaskById(ctx, id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Task not found")
		return
	}

	if err := h.taskService.DeleteTask(ctx, id); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete task")
		return
	}

	response.Success(c, "Task deleted successfully")
}

// GetMyTasks godoc
// @Summary Get tasks created by current user
// @Tags tasks
// @Description Returns all tasks where author == current user
// @Produce json
// @Success 200 {object} response.SuccessWrapper{data=[]dto.TaskResponse}
// @Failure 500 {object} response.ErrorResponse
// @Security BearerAuth
// @Router /tasks/my [get]
func (h *TaskHandler) GetMyTasks(c *gin.Context) {
	ctx := c.Request.Context()

	userID := c.GetString("userId")

	tasks, err := h.taskService.GetTasksByAuthor(ctx, userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch tasks")
		return
	}

	response.Success(c, mapper.ToTaskList(tasks))
}

// SubmitTaskAnswer godoc
// @Summary Submit answer for a task
// @Tags tasks
// @Description Check if user's answer is correct
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Param request body dto.TaskSubmitRequest true "User's answer"
// @Success 200 {object} response.SuccessWrapper{data=dto.TaskSubmitResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Security BearerAuth
// @Router /tasks/{id}/submit [post]
func (h *TaskHandler) SubmitTaskAnswer(c *gin.Context) {
    ctx := c.Request.Context()

    id := c.Param("id")

    var req dto.TaskSubmitRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, http.StatusBadRequest, "Invalid request body")
        return
    }

    isCorrect, err := h.taskService.SubmitAnswer(ctx, id, req.Answer)
    if err != nil {
        if err.Error() == "record not found" {
            response.Error(c, http.StatusNotFound, "Task not found")
        } else {
            response.Error(c, http.StatusInternalServerError, "Failed to check answer")
        }
        return
    }

    response.Success(c, dto.TaskSubmitResponse{
        Correct: isCorrect,
    })
}