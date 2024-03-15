package main

import (
	"fmt"
	"telescope/internal/cli"
	"telescope/internal/util"
	"telescope/internal/video"
)

func main() {
	config := cli.NewConfig()
	cli.GetUserInput(config)

	tempFolder, err := util.CreateTempFolder()
	if err != nil {
		fmt.Printf(cli.ERROR_FMT, "Failed to create temp folder: "+err.Error())
		util.CleanExit(tempFolder)
	}

	videoLength, err := video.PropeVideoLength(config.VideoPath)
	if err != nil {
		fmt.Printf(cli.ERROR_FMT, "Failed to get video duration with ffmpeg: "+err.Error())
		util.CleanExit(tempFolder)
	}

	video.CalculateThumbPositions(videoLength, config)
	fmt.Printf(cli.INFO_FMT, "Extracting thumbnails (this can take some time)...")
	err = video.ExtractThumbnails(config, tempFolder)
	if err != nil {
		fmt.Printf(cli.ERROR_FMT, "Failed to extract thumbnails with ffmpeg: "+err.Error())
		util.CleanExit(tempFolder)
	}

	fmt.Printf(cli.SUCCESS_FMT, fmt.Sprintf("%v", config.ThumbnailCount)+" thumbnails extracted successfully")

	outputName := util.CreateOutputName(config.VideoPath, config.CompressImage)

	fmt.Printf(cli.INFO_FMT, "Creating preview image...")
	err = video.CreateThumbPreview(config, tempFolder, outputName)
	if err != nil {
		fmt.Printf(cli.ERROR_FMT, "Failed to create preview image: "+err.Error())
		util.CleanExit(tempFolder)
	}

	err = util.MoveFile(tempFolder, outputName)
	if err != nil {
		fmt.Printf(cli.ERROR_FMT, "Failed to move file: "+err.Error())
		util.CleanExit(tempFolder)
	}

	fmt.Printf(cli.SUCCESS_FMT, "Preview image created successfully")
	util.CleanExit(tempFolder)
}
