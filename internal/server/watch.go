package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"path"
	"slices"
	"time"

	"github.com/linuxoid69/video_sender/utils/VideoSender/internal/telegram"
	"github.com/linuxoid69/video_sender/utils/VideoSender/internal/vars"
	"github.com/linuxoid69/video_sender/utils/VideoSender/internal/video"
)

var (
	DefaultCompressSizeMB int    = 9
	DefaultTMPDir         string = "/tmp/"
)

func watchJobs(ctx context.Context, cfg vars.Config, s Storage) {
	var tmpVideoFile string

	ticker := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			keys, err := s.Keys(ctx, "*")
			if err != nil {
				slog.Error("failed to get all keys", "error", err)
			}

			slices.Sort(keys)

			for _, key := range keys {
				res, err := s.Get(ctx, key)
				if err != nil {
					slog.Error("failed to get key", "error", err, "key", key)
				}

				var vd VideoData

				if err := json.Unmarshal([]byte(res), &vd); err != nil {
					slog.Error("failed to unmarshal key", "error", err, "key", key)
				}

				outFile := vd.VideoFile

				if vd.FileSize > video.AllowVideoSize {
					_, f := path.Split(vd.VideoFile)
					if err := video.VideoCompress(vd.VideoFile, DefaultTMPDir+f, DefaultCompressSizeMB); err != nil {
						slog.Error("failed to compress file", "file", vd.VideoFile, "error", err)
					}

					slog.Info("Finish compress file", "file", vd.VideoFile)

					outFile = DefaultTMPDir + f
					tmpVideoFile = outFile
				}

				slog.Info("Start send file", "file", outFile)

				if err := telegram.NewBot(cfg.TelegramToken, cfg.TelegramGroup).
					SendVideo(vd.CameraName, outFile); err != nil {
					slog.Error("failed to send video file", "file", vd.VideoFile, "error", err)
				}

				slog.Info("File was sent successfuly", "file", outFile)

				if tmpVideoFile != "" {
					if err := os.Remove(tmpVideoFile); err != nil {
						slog.Error("failed to remove temp video file", "file", tmpVideoFile, "error", err)
					}
					slog.Info("Temp file was delete successfuly", "file", tmpVideoFile)
				}

				if err := s.Delete(ctx, key); err != nil {
					slog.Error("failed to delete key", "key", key, "error", err)
				}

				slog.Info("Job was delete successfuly", "job", key)
			}
		}
	}
}
