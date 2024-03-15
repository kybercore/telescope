package video

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
	"telescope/internal/cli"
)

func PropeVideoLength(videoPath string) (float64, error) {
	cmd := exec.Command(
		"ffprobe",
		"-v",
		"error",
		"-show_entries",
		"format=duration",
		"-of",
		"default=noprint_wrappers=1:nokey=1",
		videoPath,
	)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	// trim duration to one decimal behind the comma
	duration, err := strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
	if err != nil {
		return 0, err
	}

	return duration, nil
}

func CalculateThumbPositions(videoLength float64, config *cli.Config) {
	thumbnailCount := config.ThumbnailCount
	thumbPositions := make([]float64, thumbnailCount)

	// calculate thumb positions
	for i := 0; i < thumbnailCount; i++ {
		position := (float64(i) / float64(thumbnailCount)) * videoLength
		// trim to 1 decimal && append to array
		thumbPositions[i] = float64(int(position*10)) / 10
	}

	config.ThumbPositions = thumbPositions
}

func CalculatePreviewHeight(config *cli.Config, height int) int {
	thumbnailCount := config.ThumbnailCount
	numberOfRows := thumbnailCount / 3
	if thumbnailCount%3 != 0 {
		numberOfRows++
	}
	return numberOfRows * height
}
