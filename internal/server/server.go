package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/gin-gonic/gin"
	"github.com/linuxoid69/video_sender/utils/VideoSender/internal/queue"
	"github.com/linuxoid69/video_sender/utils/VideoSender/internal/redis"
	"github.com/linuxoid69/video_sender/utils/VideoSender/internal/vars"
)

func Run(ctx context.Context) {
	var queueClient queue.Queuer

	var cfg vars.Config
	err := env.Parse(&cfg)
	if err != nil {
		fmt.Println("Need set env variables")
		os.Exit(1)
	}

	cfg, err = env.ParseAs[vars.Config]()

	shutdownCtx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	queueClient = redis.NewClient(cfg.RedisHost, cfg.RedisPassword, 0, 0)

	handler := NewHandler(queueClient.(*redis.Client))
	defer handler.redisClient.RedisClient.Close()

	r := gin.Default()

	r.POST("/addjob", handler.AddJob)

	srv := &http.Server{
		Addr:         ":8090",
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		fmt.Println("Server starting on :8090")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server failed: %v\n", err)
			stop()
		}
	}()

	go watchJobs(ctx, cfg, queueClient)

	<-shutdownCtx.Done()
	fmt.Println("Server is shutting down...")

	shutdownTimeoutCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownTimeoutCtx); err != nil {
		fmt.Printf("Server shutdown failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Server stopped gracefully")
}
