package handler

import (
	"net/http"

	"learning-platform/internal/models"
	"learning-platform/internal/service"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	taskService *service.TaskService
}

func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

type CreateTaskRequest struct {
	Title           string                `json:"title" binding:"required"`
	BodyMD          string                `json:"body_md" binding:"required"`
	Difficulty      models.Difficulty      `json:"difficulty" binding:"required,oneof=EASY MEDIUM HARD EXTREME"`
	Status          models.TaskStatus      `json:"status" binding:"required,oneof=DRAFT PUBLISHED ARCHIVED"`
	TopicID         string                `json:"topic_id" binding:"required,uuid"`
	OfficialSolution string               `json:"official_solution,omitempty"`
	CorrectAnswer   string                `json:"correct_answer,omitempty"`
	AnswerType      models.AnswerType      `json:"answer_type" binding:"required,oneof=TEXT NUMBER FORMULA"`
	ImageURL        string                `json:"image_url,omitempty"`
}

type UpdateTaskRequest struct {
	Title           string                `json:"title" binding:"required"`
	BodyMD          string                `json:"body_md" binding:"required"`
	Difficulty      models.Difficulty      `json:"difficulty" binding:"required,oneof=EASY MEDIUM HARD EXTREME"`
	Status          models.TaskStatus      `json:"status" binding:"required,oneof=DRAFT PUBLISHED ARCHIVED"`
	TopicID         string                `json:"topic_id" binding:"required,uuid"`
	OfficialSolution string               `json:"official_solution,omitempty"`
	CorrectAnswer   string                `json:"correct_answer,omitempty"`
	AnswerType      models.AnswerType      `json:"answer_type" binding:"required,oneof=TEXT NUMBER FORMULA"`
	ImageURL        string                `json:"image_url,omitempty"`
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	task := &models.Task{
		Title:           req.Title,
		BodyMD:          req.BodyMD,
		Difficulty:      req.Difficulty,
		Status:          req.Status,
		TopicID:         req.TopicID,
		AuthorID:        userID.(string),
		OfficialSolution: req.OfficialSolution,
		CorrectAnswer:   req.CorrectAnswer,
		AnswerType:      req.AnswerType,
		ImageURL:        req.ImageURL,
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

	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify task exists and user is author
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

	task := &models.Task{
		ID:              id,
		Title:           req.Title,
		BodyMD:          req.BodyMD,
		Difficulty:      req.Difficulty,
		Status:          req.Status,
		TopicID:         req.TopicID,
		AuthorID:        existingTask.AuthorID,
		OfficialSolution: req.OfficialSolution,
		CorrectAnswer:   req.CorrectAnswer,
		AnswerType:      req.AnswerType,
		ImageURL:        req.ImageURL,
	}

	if err := h.taskService.UpdateTask(task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	// Get updated task with relations
	updatedTask, _ := h.taskService.GetTask(id)
	c.JSON(http.StatusOK, updatedTask)
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