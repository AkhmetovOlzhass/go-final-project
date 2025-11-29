package handler

import (
	"net/http"

	"learning-platform/internal/models"
	"learning-platform/internal/response"
	"learning-platform/internal/mapper"
	"learning-platform/internal/service"
	"learning-platform/internal/dto"

	"github.com/gin-gonic/gin"
)

type TopicHandler struct {
	topicService *service.TopicService
}

func NewTopicHandler(topicService *service.TopicService) *TopicHandler {
	return &TopicHandler{topicService: topicService}
}

// Create godoc
// @Summary Create topic
// @Tags topics
// @Description Creates a new topic
// @Accept json
// @Produce json
// @Param request body dto.CreateTopicRequest true "Topic payload"
// @Success 201 {object} response.SuccessWrapper{data=dto.TopicResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Security BearerAuth
// @Router /topics [post]
func (h *TopicHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.CreateTopicRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request")
		return
	}

	topic := &models.Topic{
		Title:       req.Title,
		Slug:        req.Slug,
		ParentID:    req.ParentID,
		SchoolClass: req.SchoolClass,
	}

	if err := h.topicService.CreateTopic(ctx, topic); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithStatus(c, http.StatusCreated, mapper.ToTopicResponse(topic))
}

// GetAll godoc
// @Summary Get all topics
// @Tags topics
// @Description Returns all topics
// @Produce json
// @Success 200 {object} response.SuccessWrapper{data=[]dto.TopicResponse}
// @Failure 500 {object} response.ErrorResponse
// @Router /topics [get]
func (h *TopicHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()

	topics, err := h.topicService.GetAllTopics(ctx)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, mapper.ToTopicList(topics))
}

// GetByID godoc
// @Summary Get topic by ID
// @Tags topics
// @Description Returns topic by ID
// @Produce json
// @Param id path string true "Topic ID"
// @Success 200 {object} response.SuccessWrapper{data=dto.TopicResponse}
// @Failure 404 {object} response.ErrorResponse
// @Router /topics/{id} [get]
func (h *TopicHandler) GetByID(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")
	topic, err := h.topicService.GetTopicById(ctx, id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Not found")
		return
	}

	response.Success(c, mapper.ToTopicResponse(topic))
}

// Update godoc
// @Summary Update topic
// @Tags topics
// @Description Updates topic by ID
// @Accept json
// @Produce json
// @Param id path string true "Topic ID"
// @Param request body dto.UpdateTopicRequest true "Update payload"
// @Success 200 {object} response.SuccessWrapper{data=dto.TopicResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Security BearerAuth
// @Router /topics/{id} [put]
func (h *TopicHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	var req dto.UpdateTopicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	topic := models.Topic{
		ID:          id,
		Title:       req.Title,
		Slug:        req.Slug,
		ParentID:    req.ParentID,
		SchoolClass: req.SchoolClass,
	}

	if err := h.topicService.UpdateTopic(ctx, &topic); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	finalTopic, err := h.topicService.GetTopicById(ctx, id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to fetch updated topic")
		return
	}

	response.Success(c, mapper.ToTopicResponse(finalTopic))
}

// Delete godoc
// @Summary Delete topic
// @Tags topics
// @Description Deletes topic by ID
// @Produce json
// @Param id path string true "Topic ID"
// @Success 200 {object} response.SuccessWrapper{data=string}
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Security BearerAuth
// @Router /topics/{id} [delete]
func (h *TopicHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	_, err := h.topicService.GetTopicById(ctx, id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Topic not found")
		return
	}

	if err := h.topicService.DeleteTopic(ctx, id); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, "Deleted")
}
