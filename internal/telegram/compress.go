package telegram

import (
	"fmt"
	"os"
	"os/exec"
)

func compressForTelegram(inputPath, outputPath string, targetSizeMB int) error {
	targetSize := targetSizeMB * 1024 * 1024

	// Начинаем с нормального качества и увеличиваем сжатие если нужно
	crfValues := []int{35, 40, 50}

	for _, crf := range crfValues {
		fmt.Printf("Пробуем сжать с CRF=%d...\n", crf)

		cmd := exec.Command("ffmpeg",
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

		if err := cmd.Run(); err != nil {
			return err
		}

		// Проверяем размер
		if info, err := os.Stat(outputPath); err == nil {
			// sizeMB := info.Size() / (1024 * 1024)
			// fmt.Printf("Получился файл %dMB\n", sizeMB)

			if info.Size() <= int64(targetSize) {
				// fmt.Printf("Успех! Файл %dMB при CRF=%d\n", sizeMB, crf)
				return nil
			}
		}

		// Удаляем неудачную попытку
		if err := os.Remove(outputPath); err != nil {
			return err
		}
	}

	return fmt.Errorf("не удалось сжать до %dMB", targetSizeMB)
}
