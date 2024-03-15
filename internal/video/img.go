package video

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"telescope/internal/cli"
)

const LOADING_BAR = "█"

func ExtractThumbnails(config *cli.Config, tempFolder string) error {
	// reference: ffmpeg -i input.mp4 -ss 00:00:03.2 -vframes 1 output.png
	videoPath := config.VideoPath
	thumbPositions := config.ThumbPositions

	// print progress bar
	fmt.Printf("[-] Thumbnail %d/%d ", 0, len(thumbPositions))
	for i, position := range thumbPositions {
		fmt.Printf("%s", LOADING_BAR)
		if i == len(thumbPositions)-1 {
			fmt.Printf("\n")
		}
		thumbPath := filepath.Join(tempFolder, fmt.Sprintf("thumb_%d.png", i+1))
		cmd := exec.Command(
			"ffmpeg",
			"-hide_banner",
			"-loglevel",
			"error",
			"-ss",
			fmt.Sprintf("%f", position),
			"-i",
			videoPath,
			"-vframes",
			"1",
			"-vf",
			"scale=640:-1",
			thumbPath,
		)
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return err
		}

		// print progress bar
		fmt.Printf("\r[-] Thumbnail %d/%d ", i+1, len(thumbPositions))
	}

	fmt.Println()
	return nil
}

func CreateThumbPreview(config *cli.Config, tempFolder string, outputName string) error {
	// iterate through the images in temp folder
	var decodedImages = make([]image.Image, 0, config.ThumbnailCount)

	for i := 1; i <= config.ThumbnailCount; i++ {
		thumbPath := filepath.Join(tempFolder, fmt.Sprintf("thumb_%d.png", i))
		img := openAndDecode(thumbPath)
		if img != nil {
			decodedImages = append(decodedImages, img)
		} else {
			return fmt.Errorf("Failed to decode image: %s", thumbPath)
		}
	}

	// get the size of the first thumbnail
	height, width := decodedImages[0].Bounds().Dy(), decodedImages[0].Bounds().Dx()
	previewHeight := CalculatePreviewHeight(config, height)

	finalImg := image.NewRGBA(image.Rect(0, 0, width*3, previewHeight)) // three images per row

	// fill the final image with dark gray color
	gray := color.RGBA{30, 30, 30, 255} // darker gray color
	draw.Draw(finalImg, finalImg.Bounds(), &image.Uniform{gray}, image.Point{}, draw.Src)

	for i, img := range decodedImages {
		// draw the image on the final image
		x := (i % 3) * width
		y := (i / 3) * height
		draw.Draw(finalImg, image.Rect(x, y, x+width, y+height), img, image.Point{0, 0}, draw.Src)
	}

	// save the final image
	finalImgPath := filepath.Join(tempFolder, outputName)
	file, err := os.Create(finalImgPath)
	if err != nil {
		return err
	}
	defer file.Close()

	if config.CompressImage {
		fmt.Printf(cli.INFO_FMT, "Compressing final image...")
		err = jpeg.Encode(file, finalImg, &jpeg.Options{Quality: 75})
	} else {
		err = png.Encode(file, finalImg)
	}

	if err != nil {
		return err
	}

	return nil
}

func openAndDecode(imgPath string) image.Image {
	file, err := os.Open(imgPath)
	if err != nil {
		return nil
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil
	}

	return img
}
