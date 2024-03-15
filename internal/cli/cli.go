package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"telescope/internal/util"
)

const (
	PROMPT_FMT  = "[?] %s: "
	ERROR_FMT   = "\033[31m[!] %s\033[0m\n"
	WARN_FMT    = "\033[33m[!] %s\033[0m\n"
	INFO_FMT    = "\033[36m[i] %s\033[0m\n"
	SUCCESS_FMT = "\033[32m[+] %s\033[0m\n"
)

func NewConfig() *Config {
	return &Config{ThumbnailCount: 0, VideoPath: "", CompressImage: false}
}

func GetUserInput(userInput *Config) {
	reader := bufio.NewReader(os.Stdin)

	for {
		var videoPath string
		fmt.Printf(PROMPT_FMT, "Path to video file (try dragging and dropping the file here)")
		videoPath, _ = reader.ReadString('\n')
		videoPath = strings.TrimSpace(videoPath)
		videoPath = strings.Trim(videoPath, `"'`)

		if videoPath != "" && util.PathExists(videoPath) {
			userInput.VideoPath = videoPath
			break
		}
		fmt.Printf(WARN_FMT, "Invalid path to video file")
	}

	for {
		var thumbnailCount int
		fmt.Printf(PROMPT_FMT, "Number of thumbnails")
		fmt.Scanln(&thumbnailCount)

		if thumbnailCount > 0 {
			userInput.ThumbnailCount = thumbnailCount
			break
		}
		fmt.Printf(WARN_FMT, "Number of thumbnails must be greater than 0")
	}

	for {
		var confirm string
		fmt.Printf(PROMPT_FMT, "Optimize final image for online sharing? (y/N)")
		confirm, _ = reader.ReadString('\n')
		confirm = strings.TrimSpace(confirm)
		if confirm == "y" || confirm == "n" || confirm == "" {
			if confirm == "y" {
				userInput.CompressImage = true
			}
			break
		}
		fmt.Printf(WARN_FMT, "Invalid input")
	}
}
