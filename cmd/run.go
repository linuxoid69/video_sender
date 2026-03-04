package cmd

import (
	"context"
	"fmt"

	"github.com/linuxoid69/video_sender/utils/VideoSender/internal/logger"
	"github.com/linuxoid69/video_sender/utils/VideoSender/internal/server"
	"github.com/spf13/cobra"
)

func RunCmd(_ *cobra.Command, _ []string) *cobra.Command {
	return &cobra.Command{
		Use:              "run",
		PersistentPreRun: func(_ *cobra.Command, _ []string) { logger.InitLogger() },
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println("Run server")

			ctx := context.WithoutCancel(context.Background())

			server.Run(ctx)
		},
	}
}
