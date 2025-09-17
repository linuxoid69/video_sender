package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.my-itclub.ru/utils/VideoSender/internal/telegram"
	"github.com/gin-gonic/gin"
)

func Run(ctx context.Context) {
	if err := telegram.CheckEnvVars(); err != nil {
		fmt.Printf("Environment validation failed: %v\n", err)
		os.Exit(1)
	}

	shutdownCtx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	r := gin.Default()
	r.POST("/video", HandlerGetVideo)

	srv := &http.Server{
		Addr:    ":8090",
		Handler: r,
		ReadTimeout:  30 * time.Second,    // Увеличил для загрузки
		WriteTimeout: 30 * time.Second,    // Увеличил для выгрузки
		IdleTimeout:  120 * time.Second,   // Увеличил для долгих операций
	}

	go func() {
		fmt.Println("Server starting on :8090")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server failed: %v\n", err)
			stop()
		}
	}()

	<-shutdownCtx.Done()
	fmt.Println("Server is shutting down...")

	shutdownTimeoutCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // 30 секунд для больших видео
	defer cancel()

	if err := srv.Shutdown(shutdownTimeoutCtx); err != nil {
		fmt.Printf("Server shutdown failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Server stopped gracefully")
}
