package server

import (
	"net/http"
	"os"
	"strconv"

	"git.my-itclub.ru/utils/VideoSender/internal/telegram"
	"github.com/gin-gonic/gin"
)

type Video struct {
	Path   string `json:"video_path"`
	Camera string `json:"camera"`
}

func HandlerGetVideo(c *gin.Context) {
	var video Video

	if err := c.BindJSON(&video); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tgGroup, err := strconv.Atoi(os.Getenv("TG_GROUP"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	t := telegram.NewBot(os.Getenv("TG_TOKEN"), int64(tgGroup))

	if err := t.SendVideo(video.Camera, video.Path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "request got successfully"})
}
