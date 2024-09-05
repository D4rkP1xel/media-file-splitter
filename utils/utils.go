package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

// Helper function to parse the duration from FFmpeg output
func ParseFFMPEGDuration(output string) string {
	// Example FFmpeg output line: "Duration: 00:02:30.15"
	const durationPrefix = "Duration: "
	startIndex := strings.Index(output, durationPrefix)

	if startIndex == -1 {
		return ""
	}

	startIndex += len(durationPrefix)
	endIndex := startIndex + 8 // Duration string is 8 characters long: HH:MM:SS
	return output[startIndex:endIndex]
}

func GenerateChunk(startTime float64, secondsPerChunk uint16,
	inputFilePath string, outputFilePath string,
	errChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	cmd := exec.Command("ffmpeg", "-i", inputFilePath, "-ss", fmt.Sprintf("%f", startTime), "-t", fmt.Sprintf("%d", secondsPerChunk), "-acodec", "libmp3lame", outputFilePath)
	output, err := cmd.CombinedOutput()

	if err != nil {
		errChan <- fmt.Errorf("Error splitting file: %s\nOutput: %s", err.Error(), string(output))
		return
	}

	errChan <- nil
}

type FileData struct {
	InputFileInfo os.FileInfo
	Duration      uint32
}

func HandleParams(createFolderIfNotExists []bool, secondsPerChunk int, outputDirectoryPath string, inputFilePath string) (FileData, error) {
	createOutputFolder := false
	if len(createFolderIfNotExists) > 0 {
		createOutputFolder = createFolderIfNotExists[0]
	}

	if secondsPerChunk <= 0 {
		return FileData{}, fmt.Errorf("Insert a valid number of secondsPerChunk (>0)")
	}

	// check if output folder exists
	_, err := os.Stat(outputDirectoryPath)
	if err != nil {
		//output folder not found
		if !createOutputFolder {
			return FileData{}, fmt.Errorf("Error finding output directory: %s", err.Error())
		}

		// create output folder
		err = os.Mkdir(outputDirectoryPath, 0744)
		if err != nil {
			return FileData{}, fmt.Errorf("Error creating output directory: %s", err.Error())
		}
	}

	// check if input file exists
	fileInfo, err := os.Stat(inputFilePath)
	if err != nil {
		return FileData{}, fmt.Errorf("Error finding input file: %s", err.Error())
	}

	cmd := exec.Command("ffmpeg", "-i", inputFilePath, "-f", "null", "-")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return FileData{}, fmt.Errorf("Error getting file duration: %s", err)
	}

	// Parse the duration from the FFmpeg output
	durationStr := ParseFFMPEGDuration(string(output))
	if durationStr == "" {
		return FileData{}, fmt.Errorf("Could not determine file duration")
	}

	durationStrSplit := strings.Split(durationStr, ":")
	hours, err := strconv.ParseInt(durationStrSplit[0], 10, 8)
	minutes, err := strconv.ParseInt(durationStrSplit[1], 10, 0)
	seconds, err := strconv.ParseInt(durationStrSplit[2], 10, 0)

	var duration uint32 = uint32(hours*3600 + minutes*60 + seconds)
	return FileData{
		InputFileInfo: fileInfo,
		Duration:      duration,
	}, nil
}

func HandleCloseChannel(errChan chan error, numberOfChunks uint16, wg *sync.WaitGroup) error {
	wg.Wait()
	close(errChan)

	var i uint16
	// Collect errors from goroutines
	for i = 0; i < uint16(numberOfChunks); i++ {
		if err := <-errChan; err != nil {
			return err
		}
	}

	return nil
}
