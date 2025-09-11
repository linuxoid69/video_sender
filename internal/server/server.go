package server

import (
	"context"
	"fmt"
	"os"

	"git.my-itclub.ru/utils/VideoSender/internal/telegram"
	"github.com/gin-gonic/gin"
)

func Run(ctx context.Context) {
	if err := telegram.CheckEnvVars(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	r := gin.Default()
	r.POST("/video", HandlerGetVideo)

	panic(r.Run(":8090"))
}
