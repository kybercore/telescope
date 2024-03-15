package util

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func CreateOutputName(videoPath string, compress bool) string {
	var outputName string
	switch compress {
	case true:
		outputName = filepath.Base(videoPath) + "_preview.jpg"
	case false:
		outputName = filepath.Base(videoPath) + "_preview.png"
	}
	return outputName
}

func CreateTempFolder() (string, error) {
	tempDir := os.TempDir() + "/telescope"
	err := os.Mkdir(tempDir, 0755)
	if err != nil {
		if os.IsExist(err) {
			os.RemoveAll(tempDir)
			err = os.Mkdir(tempDir, 0755)
		}
	}

	return tempDir, err
}

func CleanExit(tempFolder string) {
	os.RemoveAll(tempFolder)
	fmt.Println("\nPress any key to exit...")
	fmt.Scanln()
	os.Exit(0)
}

func MoveFile(tempFolder, outputName string) error {
	workingDir, _ := os.Getwd()
	destination := filepath.Join(workingDir, outputName)
	source := filepath.Join(tempFolder, outputName)
	destinationFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}
