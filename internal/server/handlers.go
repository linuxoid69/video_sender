package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/linuxoid69/video_sender/utils/VideoSender/internal/queue"
	"github.com/linuxoid69/video_sender/utils/VideoSender/internal/redis"
)

type Video struct {
	Path   string `json:"video_path"`
	Camera string `json:"camera"`
}

type Handler struct {
	redisClient *redis.Client
}

func NewHandler(rdb *redis.Client) *Handler {
	return &Handler{redisClient: rdb}
}

func (h *Handler) AddJob(c *gin.Context) {
	var q queue.Query
	if err := c.BindJSON(&q); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if q.TTL != 0 {
		h.redisClient.KeyTTL = time.Duration(q.TTL)
	}

	data, err := json.Marshal(q.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.redisClient.CreateJob(c.Request.Context(), q.Key, string(data)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "request got successfully"})
}
