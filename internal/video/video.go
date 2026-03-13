package video

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
)

const (
	AllowVideoSize    int64 = 10 * 1024 * 1024
	AllowMaxVideoSize int64 = 50 * 1024 * 1024
)

func VideoCompress(ctx context.Context, inputPath, outputPath string, targetSizeMB int) (err error) {
	targetSize := targetSizeMB * 1024 * 1024

	// Начинаем с нормального качества и увеличиваем сжатие если нужно
	crfValues := []int{35, 40, 50}

	for _, crf := range crfValues {
		slog.Info("Try compress", "CRF", crf, "file", inputPath)

		cmd := exec.CommandContext(ctx, "ffmpeg",
			"-i", inputPath,
			"-c:v", "libx264",
			"-crf", fmt.Sprintf("%d", crf),
			"-preset", "medium",
			"-pix_fmt", "yuv420p",
			"-c:a", "aac",
			"-b:a", "96k",
			"-movflags", "+faststart",
			"-vf", "scale=trunc(iw/2)*2:trunc(ih/2)*2:flags=lanczos",
			outputPath,
		)

		if err = cmd.Run(); err != nil {
			return err
		}

		var fileInfo os.FileInfo

		if fileInfo, err = os.Stat(outputPath); err == nil {
			if fileInfo.Size() <= int64(targetSize) {
				return nil
			}
		}
	}

	return fmt.Errorf("failed to compress file %s for %dMB", inputPath, targetSizeMB)
}
