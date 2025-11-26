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

func (h *TaskHandler) GetAllTasks(c *gin.Context) {
	tasks, err := h.taskService.GetAllTasks()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to fetch tasks")
		return
	}

	response.Success(c, mapper.ToTaskList(tasks))
}

func (h *TaskHandler) GetDraftTasks(c *gin.Context) {
	tasks, err := h.taskService.GetDraftTasks()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to fetch draft tasks")
		return
	}

	response.Success(c, mapper.ToTaskList(tasks))
}

func (h *TaskHandler) PublishTask(c *gin.Context) {
    id := c.Param("id")

    existingTask, err := h.taskService.GetTask(id)
    if err != nil {
        response.Error(c, http.StatusNotFound, "Task not found")
        return
    }

    userID, exists := c.Get("userId")
    if !exists || existingTask.AuthorID != userID.(string) {
		response.Error(c, http.StatusForbidden, "Not authorized to publish this task")
        return
    }

    if err := h.taskService.PublishTask(id); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to publish task")
        return
    }

    updated, _ := h.taskService.GetTask(id)
    response.Success(c, mapper.ToTaskResponse(updated))
}



func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req dto.CreateTaskRequest

	if err := c.ShouldBindWith(&req, binding.FormMultipart); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	userID, exists := c.Get("userId")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var imageURL string
	file, err := c.FormFile("imageUrl") 
	if err == nil && file != nil {
		url, uploadErr := h.s3.UploadFile(file)
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
		AuthorID:         userID.(string),
		OfficialSolution: req.OfficialSolution,
		CorrectAnswer:    req.CorrectAnswer,
		AnswerType:       req.AnswerType,
		ImageURL:         imageURL,
	}

	if err := h.taskService.CreateTask(task); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create task")
		return
	}

	response.SuccessWithStatus(c, http.StatusCreated, mapper.ToTaskResponse(task))
}


func (h *TaskHandler) GetTask(c *gin.Context) {
	id := c.Param("id")

	task, err := h.taskService.GetTask(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Task not found")
		return
	}

	response.Success(c, mapper.ToTaskResponse(task))
}

func (h *TaskHandler) GetTasksByTopic(c *gin.Context) {
	topicID := c.Param("topicId")

	tasks, err := h.taskService.GetTasksByTopic(topicID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch tasks")
		return
	}

	response.Success(c, mapper.ToTaskList(tasks))
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
    id := c.Param("id")

    existingTask, err := h.taskService.GetTask(id)
    if err != nil {
        response.Error(c, http.StatusNotFound, "Task not found")
        return
    }

    userID, exists := c.Get("userId")
    if !exists || existingTask.AuthorID != userID.(string) {
        response.Error(c, http.StatusForbidden, "Not authorized to update this task")
        return
    }

    var req dto.UpdateTaskRequest
    if err := c.ShouldBindWith(&req, binding.FormMultipart); err != nil {
        response.Error(c, http.StatusBadRequest, err.Error())
        return
    }

    imageURL := existingTask.ImageURL

    file, err := c.FormFile("imageUrl")
    if err == nil && file != nil {
        url, uploadErr := h.s3.UploadFile(file)
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
        AuthorID:         existingTask.AuthorID,
        OfficialSolution: req.OfficialSolution,
        CorrectAnswer:    req.CorrectAnswer,
        AnswerType:       req.AnswerType,
        ImageURL:         imageURL,
    }

    if err := h.taskService.UpdateTask(updated); err != nil {
        response.Error(c, http.StatusInternalServerError, "Failed to update task")
        return
    }

    finalTask, _ := h.taskService.GetTask(id)
    response.Success(c, mapper.ToTaskResponse(finalTask))
}



func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id := c.Param("id")

	existingTask, err := h.taskService.GetTask(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Task not found")
		return
	}

	userID, exists := c.Get("userId")
	if !exists || existingTask.AuthorID != userID.(string) {
		response.Error(c, http.StatusForbidden, "Not authorized to delete this task")
		return
	}

	if err := h.taskService.DeleteTask(id); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete task")
		return
	}

	response.Success(c, "Task deleted successfully")
}

func (h *TaskHandler) GetMyTasks(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	tasks, err := h.taskService.GetTasksByAuthor(userID.(string))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch tasks")
		return
	}

	response.Success(c, mapper.ToTaskList(tasks))
}