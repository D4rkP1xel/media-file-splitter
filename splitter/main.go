package splitter

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

func SplitMediaFileByTimedChunks(secondsPerChunk int, inputFilePath string, outputDirectory string, createFolderIfNotExists ...bool) error {
	createOutputFolder := false
	if len(createFolderIfNotExists) > 0 {
		createOutputFolder = createFolderIfNotExists[0]
	}

	if secondsPerChunk <= 0 {
		return fmt.Errorf("Insert a valid number of secondsPerChunk (>0)")
	}

	// check if output folder exists
	_, err := os.Stat(outputDirectory)
	if err != nil {
		//output folder not found
		if !createOutputFolder {
			return fmt.Errorf("Error finding output directory: %s", err.Error())
		}

		// create output folder
		err = os.Mkdir(outputDirectory, 0744)
		if err != nil {
			return fmt.Errorf("Error creating output directory: %s", err.Error())
		}
	}

	_, err = os.Stat(inputFilePath)
	if err != nil {
		return fmt.Errorf("Error finding input file: %s", err.Error())
	}

	cmd := exec.Command("ffmpeg", "-i", inputFilePath, "-f", "null", "-")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error getting file duration: %s", err)
	}

	// Parse the duration from the FFmpeg output
	durationStr := parseFFMPEGDuration(string(output))
	if durationStr == "" {
		return fmt.Errorf("Could not determine file duration")
	}

	durationStrSplit := strings.Split(durationStr, ":")
	hours, err := strconv.ParseInt(durationStrSplit[0], 10, 8)
	minutes, err := strconv.ParseInt(durationStrSplit[1], 10, 0)
	seconds, err := strconv.ParseInt(durationStrSplit[2], 10, 0)

	var duration uint16 = uint16(hours*3600 + minutes*60 + seconds)

	// Calculate number of chunks
	var numChunks uint16 = duration / uint16(secondsPerChunk)
	if duration%uint16(secondsPerChunk) != 0 {
		numChunks++
	}

	fileType := strings.Split(inputFilePath, ".")[1]
	var wg sync.WaitGroup
	errChan := make(chan error, numChunks)
	// Split the file into chunks using FFmpeg
	var i uint16
	for i = 0; i < numChunks; i++ {
		wg.Add(1)
		go generateChunk(i, uint16(secondsPerChunk), outputDirectory, inputFilePath, fileType, errChan, &wg)
	}

	wg.Wait()
	close(errChan)

	// Collect errors from goroutines
	for i = 0; i < numChunks; i++ {
		if err := <-errChan; err != nil {
			return err
		}
	}

	return nil
}

// Helper function to parse the duration from FFmpeg output
func parseFFMPEGDuration(output string) string {
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

func generateChunk(chunkIndex uint16, secondsPerChunk uint16, outputDirectory string, inputFilePath string, fileType string, errChan chan<- error, wg *sync.WaitGroup) {
	startTime := chunkIndex * uint16(secondsPerChunk)
	outputFile := fmt.Sprintf("%s/chunk_%03d.%s", outputDirectory, chunkIndex+1, fileType)

	cmd := exec.Command("ffmpeg", "-i", inputFilePath, "-ss", fmt.Sprintf("%d", startTime), "-t", fmt.Sprintf("%d", secondsPerChunk), "-acodec", "libmp3lame", outputFile)
	output, err := cmd.CombinedOutput()

	if err != nil {
		errChan <- fmt.Errorf("Error splitting file: %s\nOutput: %s", err.Error(), string(output))
		return
	}
	wg.Done()
	errChan <- nil
}
