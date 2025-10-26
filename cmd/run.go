package cmd

import (
	"context"
	"fmt"

	"github.com/linuxoid69/video_sender/utils/VideoSender/internal/logger"
	"github.com/linuxoid69/video_sender/utils/VideoSender/internal/server"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:              "run",
	PersistentPreRun: func(cmd *cobra.Command, args []string) { logger.InitLogger() },
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Run server")

		ctx := context.WithoutCancel(context.Background())

		server.Run(ctx)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
