package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Video struct {
	Path string `json:"video_path"`
}

func HandlerGetVideo(c *gin.Context) {
	var video Video

	if err := c.BindJSON(&video); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "video path got successfully",
	})
}
