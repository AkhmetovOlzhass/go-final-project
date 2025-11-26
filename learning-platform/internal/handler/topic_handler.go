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
    topics *service.TopicService
}

func NewTopicHandler(topics *service.TopicService) *TopicHandler {
    return &TopicHandler{topics: topics}
}

func (h *TopicHandler) Create(c *gin.Context) {
    var req dto.CreateTopicRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, http.StatusBadRequest, "Invalid request")
        return
    }

    topic := models.Topic{
        Title:       req.Title,
        Slug:        req.Slug,
        ParentID:    req.ParentID,
        SchoolClass: req.SchoolClass,
    }

    if err := h.topics.Create(&topic); err != nil {
        response.Error(c, http.StatusInternalServerError, err.Error())
        return
    }

    response.SuccessWithStatus(c, http.StatusCreated, mapper.ToTopicResponse(&topic))
}

func (h *TopicHandler) GetAll(c *gin.Context) {
    topics, err := h.topics.GetAll()
    if err != nil {
        response.Error(c, http.StatusInternalServerError, err.Error())
        return
    }
    response.Success(c, mapper.ToTopicList(topics))
}

func (h *TopicHandler) GetByID(c *gin.Context) {
    id := c.Param("id")
    topic, err := h.topics.GetByID(id)
    if err != nil {
        response.Error(c, http.StatusNotFound, "Not found")
        return
    }
    response.Success(c, mapper.ToTopicResponse(topic))
}

func (h *TopicHandler) Update(c *gin.Context) {
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

    if err := h.topics.Update(&topic); err != nil {
        response.Error(c, http.StatusInternalServerError, err.Error())
        return
    }
    response.Success(c, mapper.ToTopicResponse(&topic))
}

func (h *TopicHandler) Delete(c *gin.Context) {
    id := c.Param("id")
    if err := h.topics.Delete(id); err != nil {
        response.Error(c, http.StatusInternalServerError, err.Error())
        return
    }
    response.Success(c, "Deleted")
}
