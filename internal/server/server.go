package server

import (
	"context"

	"github.com/gin-gonic/gin"
)

func Run(ctx context.Context) {

	r := gin.Default()
	r.POST("/video", HandlerGetVideo)

	panic(r.Run(":8090"))
}
