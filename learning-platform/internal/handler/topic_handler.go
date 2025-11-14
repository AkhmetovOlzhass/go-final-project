package handler

import (
    "net/http"
    "learning-platform/internal/models"
    "learning-platform/internal/service"
    "github.com/gin-gonic/gin"
)

type TopicHandler struct {
    topics *service.TopicService
}

func NewTopicHandler(topics *service.TopicService) *TopicHandler {
    return &TopicHandler{topics: topics}
}

func (h *TopicHandler) Create(c *gin.Context) {
    var topic models.Topic
    if err := c.ShouldBindJSON(&topic); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }
    if err := h.topics.Create(&topic); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, topic)
}

func (h *TopicHandler) GetAll(c *gin.Context) {
    topics, err := h.topics.GetAll()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, topics)
}

func (h *TopicHandler) GetByID(c *gin.Context) {
    id := c.Param("id")
    topic, err := h.topics.GetByID(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
        return
    }
    c.JSON(http.StatusOK, topic)
}

func (h *TopicHandler) Update(c *gin.Context) {
    id := c.Param("id")
    var topic models.Topic
    if err := c.ShouldBindJSON(&topic); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }
    topic.ID = id
    if err := h.topics.Update(&topic); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, topic)
}

func (h *TopicHandler) Delete(c *gin.Context) {
    id := c.Param("id")
    if err := h.topics.Delete(id); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
