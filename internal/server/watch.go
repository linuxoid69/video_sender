package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"path"
	"slices"
	"time"

	"github.com/linuxoid69/video_sender/utils/VideoSender/internal/queue"
	"github.com/linuxoid69/video_sender/utils/VideoSender/internal/telegram"
	"github.com/linuxoid69/video_sender/utils/VideoSender/internal/vars"
	"github.com/linuxoid69/video_sender/utils/VideoSender/internal/video"
)

var (
	DefaultCompressSizeMB int    = 9
	DefaultTMPDir         string = "/tmp/"
)

func watchJobs(ctx context.Context, cfg vars.Config, q queue.Queuer) {
	var tmpVideoFile string
	for {
		keys, err := q.GetKeys(ctx, "*")
		if err != nil {
			slog.Error("can't get all keys", "error", err)
		}

		slices.Sort(keys)

		for _, key := range keys {
			res, err := q.GetJob(ctx, key)
			if err != nil {
				slog.Error("can't get key", "error", err, "key", key)
			}

			var r queue.VideoData
			if err := json.Unmarshal([]byte(res), &r); err != nil {
				slog.Error("can't unmarshal key", "error", err, "key", key)
			}

			outFile := r.VideoFile

			if r.FileSize > video.AllowVideoSize {
				_, f := path.Split(r.VideoFile)
				if err := video.VideoCompress(r.VideoFile, DefaultTMPDir+f, DefaultCompressSizeMB); err != nil {
					slog.Error("can't compress file", "file", r.VideoFile, "error", err)
				}

				slog.Info("Finish compress file", "file", r.VideoFile)

				outFile = DefaultTMPDir + f
				tmpVideoFile = outFile
			}

			slog.Info("Start send file", "file", outFile)

			if err := telegram.NewBot(cfg.TelegramToken, cfg.TelegramGroup).
				SendVideo(r.CameraName, outFile); err != nil {
				slog.Error("can't send video file", "file", r.VideoFile, "error", err)
			}

			slog.Info("File was sent successfuly", "file", outFile)

			if tmpVideoFile != "" {
				if err := os.Remove(tmpVideoFile); err != nil {
					slog.Error("Can't remove temp video file", "file", tmpVideoFile, "error", err)
				}
				slog.Info("Temp file was delete successfuly", "file", tmpVideoFile)
			}

			if err := q.DeleteJob(ctx, key); err != nil {
				slog.Error("Can't delete key", "key", key, "error", err)
			}

			slog.Info("Job was delete successfuly", "job", key)
		}

		time.Sleep(time.Duration(1 * time.Second))
	}
}
