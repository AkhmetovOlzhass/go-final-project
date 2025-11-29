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
	ctx := c.Request.Context()
	
	tasks, err := h.taskService.GetAllTasks(ctx)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to fetch tasks")
		return
	}

	response.Success(c, mapper.ToTaskList(tasks))
}

func (h *TaskHandler) GetDraftTasks(c *gin.Context) {
	ctx := c.Request.Context()

	tasks, err := h.taskService.GetDraftTasks(ctx)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to fetch draft tasks")
		return
	}

	response.Success(c, mapper.ToTaskList(tasks))
}

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