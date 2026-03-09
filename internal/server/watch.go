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

const (
	DefaultCompressSizeMB int    = 9
	DefaultTMPDir         string = "/tmp/"
)

func watchJobs(ctx context.Context, cfg vars.Config, s Storage) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			keys, err := s.Keys(ctx, "*")
			if err != nil {
				slog.Error("failed to get all keys", "error", err)
				continue
			}

			slices.Sort(keys)

			for _, key := range keys {
				var tmpVideoFile string
				res, err := s.Get(ctx, key)
				if err != nil {
					slog.Error("failed to get key", "error", err, "key", key)

					continue
				}

				var vd VideoData

				if err := json.Unmarshal([]byte(res), &vd); err != nil {
					slog.Error("failed to unmarshal key", "error", err, "key", key)

					continue
				}

				outFile := vd.VideoFile

				if vd.FileSize > video.AllowVideoSize {
					_, fileName := path.Split(vd.VideoFile)

					if err := video.VideoCompress(ctx, vd.VideoFile, DefaultTMPDir+fileName, DefaultCompressSizeMB); err != nil {
						slog.Warn("failed to compress file", "file", vd.VideoFile, "error", err)

						continue
					}

					slog.Info("Finish compress file", "file", vd.VideoFile)

					outFile = DefaultTMPDir + fileName
					tmpVideoFile = outFile

				}

				slog.Info("Start send file", "file", outFile)

				if err = telegram.NewBot(cfg.TelegramToken, cfg.TelegramGroup).
					SendVideo(ctx, vd.CameraName, outFile); err != nil {
					slog.Error("failed to send video file", "file", vd.VideoFile, "error", err)

					continue
				}

				slog.Info("File was sent successfuly", "file", outFile)

				if tmpVideoFile != "" {
					if err = os.Remove(tmpVideoFile); err != nil {
						slog.Error("failed to remove temp video file", "file", tmpVideoFile, "error", err)

						continue
					}

					slog.Info("Temp file was delete successfuly", "file", tmpVideoFile)
				}

				if err = s.Delete(ctx, key); err != nil {
					slog.Error("failed to delete key", "key", key, "error", err)

					continue
				}

				slog.Info("Job was delete successfuly", "job", key)
			}
		}
	}
}
