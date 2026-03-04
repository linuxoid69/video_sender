package server

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Request struct {
	Key   string    `json:"key"`
	TTL   int       `json:"ttl,omitempty"`
	Value VideoData `json:"value"`
}

type VideoData struct {
	FileSize   int64  `json:"file_size"`
	VideoFile  string `json:"video_file"`
	CameraName string `json:"camera_name"`
}

type Handler struct {
	storage Storage
}

func NewHandler(storage Storage) *Handler {
	return &Handler{storage: storage}
}

func (h *Handler) AddJob(c *gin.Context) {
	var (
		req Request
		err error
	)

	if err = c.BindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var data []byte

	data, err = json.Marshal(req.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err = h.storage.Create(c.Request.Context(), req.Key, string(data)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "request got successfully"})
}
