package cmd

import (
	"context"
	"fmt"

	"git.my-itclub.ru/utils/VideoSender/internal/server"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use: "run",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Run server")

		ctx := context.WithoutCancel(context.Background())

		server.Run(ctx)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
