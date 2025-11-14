package handler

import (
	"net/http"

	"learning-platform/internal/models"
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

type CreateTaskRequest struct {
    Title            string                `form:"title" binding:"required"`
    BodyMD           string                `form:"body_md" binding:"required"`
    Difficulty       models.Difficulty     `form:"difficulty" binding:"required"`
    Status           models.TaskStatus     `form:"status" binding:"required"`
    TopicID          string                `form:"topic_id" binding:"required"`
    OfficialSolution string                `form:"official_solution"`
    CorrectAnswer    string                `form:"correct_answer"`
    AnswerType       models.AnswerType     `form:"answer_type" binding:"required"`
}

type UpdateTaskRequest struct {
    Title            string            `form:"title" binding:"required"`
    BodyMD           string            `form:"body_md" binding:"required"`
    Difficulty       models.Difficulty `form:"difficulty" binding:"required"`
    Status           models.TaskStatus `form:"status" binding:"required"`
    TopicID          string            `form:"topic_id" binding:"required"`
    OfficialSolution string            `form:"official_solution"`
    CorrectAnswer    string            `form:"correct_answer"`
    AnswerType       models.AnswerType `form:"answer_type" binding:"required"`
}


func (h *TaskHandler) GetAllTasks(c *gin.Context) {
	tasks, err := h.taskService.GetAllTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) GetDraftTasks(c *gin.Context) {
	tasks, err := h.taskService.GetDraftTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch draft tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) PublishTask(c *gin.Context) {
    id := c.Param("id")

    existingTask, err := h.taskService.GetTask(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
        return
    }

    userID, exists := c.Get("user_id")
    if !exists || existingTask.AuthorID != userID.(string) {
        c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to publish this task"})
        return
    }

    if err := h.taskService.PublishTask(id); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish task"})
        return
    }

    updated, _ := h.taskService.GetTask(id)
    c.JSON(http.StatusOK, updated)
}



func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req CreateTaskRequest

	if err := c.ShouldBindWith(&req, binding.FormMultipart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var imageURL string
	file, err := c.FormFile("image_url") 
	if err == nil && file != nil {
		url, uploadErr := h.s3.UploadFile(file)
		if uploadErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload image"})
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
		AuthorID:         userID.(string),
		OfficialSolution: req.OfficialSolution,
		CorrectAnswer:    req.CorrectAnswer,
		AnswerType:       req.AnswerType,
		ImageURL:         imageURL,
	}

	if err := h.taskService.CreateTask(task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, task)
}


func (h *TaskHandler) GetTask(c *gin.Context) {
	id := c.Param("id")

	task, err := h.taskService.GetTask(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) GetTasksByTopic(c *gin.Context) {
	topicID := c.Param("topic_id")

	tasks, err := h.taskService.GetTasksByTopic(topicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
    id := c.Param("id")

    existingTask, err := h.taskService.GetTask(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
        return
    }

    userID, exists := c.Get("user_id")
    if !exists || existingTask.AuthorID != userID.(string) {
        c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to update this task"})
        return
    }

    var req UpdateTaskRequest
    if err := c.ShouldBindWith(&req, binding.FormMultipart); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    imageURL := existingTask.ImageURL

    file, err := c.FormFile("image_url")
    if err == nil && file != nil {
        url, uploadErr := h.s3.UploadFile(file)
        if uploadErr != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload image"})
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
        AuthorID:         existingTask.AuthorID,
        OfficialSolution: req.OfficialSolution,
        CorrectAnswer:    req.CorrectAnswer,
        AnswerType:       req.AnswerType,
        ImageURL:         imageURL,
    }

    if err := h.taskService.UpdateTask(updated); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
        return
    }

    finalTask, _ := h.taskService.GetTask(id)
    c.JSON(http.StatusOK, finalTask)
}



func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id := c.Param("id")

	// Verify task exists and user is author
	existingTask, err := h.taskService.GetTask(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists || existingTask.AuthorID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to delete this task"})
		return
	}

	if err := h.taskService.DeleteTask(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

func (h *TaskHandler) GetMyTasks(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	tasks, err := h.taskService.GetTasksByAuthor(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}