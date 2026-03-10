package main

import (
	"os"

	"github.com/linuxoid69/video_sender/utils/VideoSender/cmd"
)

func main() {
	rootCmd := cmd.RootCmd()
	rootCmd.AddCommand(cmd.RunCmd(nil, nil))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
