package splitter

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

func SplitMediaFileByTimedChunks(secondsPerChunk int, inputFilePath string, outputDirectoryPath string, createFolderIfNotExists ...bool) ([]string, error) {
	createOutputFolder := false
	if len(createFolderIfNotExists) > 0 {
		createOutputFolder = createFolderIfNotExists[0]
	}

	if secondsPerChunk <= 0 {
		return nil, fmt.Errorf("Insert a valid number of secondsPerChunk (>0)")
	}

	// check if output folder exists
	_, err := os.Stat(outputDirectoryPath)
	if err != nil {
		//output folder not found
		if !createOutputFolder {
			return nil, fmt.Errorf("Error finding output directory: %s", err.Error())
		}

		// create output folder
		err = os.Mkdir(outputDirectoryPath, 0744)
		if err != nil {
			return nil, fmt.Errorf("Error creating output directory: %s", err.Error())
		}
	}

	_, err = os.Stat(inputFilePath)
	if err != nil {
		return nil, fmt.Errorf("Error finding input file: %s", err.Error())
	}

	cmd := exec.Command("ffmpeg", "-i", inputFilePath, "-f", "null", "-")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("Error getting file duration: %s", err)
	}

	// Parse the duration from the FFmpeg output
	durationStr := parseFFMPEGDuration(string(output))
	if durationStr == "" {
		return nil, fmt.Errorf("Could not determine file duration")
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

	//grab the input file type
	fileType := strings.Split(inputFilePath, ".")[1]
	var wg sync.WaitGroup
	errChan := make(chan error, numChunks)
	// Split the file into chunks using FFmpeg
	var i uint16
	outputFilesPaths := make([]string, numChunks, numChunks)
	for i = 0; i < numChunks; i++ {
		wg.Add(1)
		outputFilePath := fmt.Sprintf("%s/chunk_%03d.%s", outputDirectoryPath, i+1, fileType)
		outputFilesPaths[i] = outputFilePath
		go generateChunk(i, uint16(secondsPerChunk), inputFilePath, outputFilePath, errChan, &wg)
	}

	wg.Wait()
	close(errChan)

	// Collect errors from goroutines
	for i = 0; i < numChunks; i++ {
		if err := <-errChan; err != nil {
			return nil, err
		}
	}

	return outputFilesPaths, nil
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

func generateChunk(chunkIndex uint16, secondsPerChunk uint16, inputFilePath string, outputFilePath string, errChan chan<- error, wg *sync.WaitGroup) {
	startTime := chunkIndex * uint16(secondsPerChunk)

	cmd := exec.Command("ffmpeg", "-i", inputFilePath, "-ss", fmt.Sprintf("%d", startTime), "-t", fmt.Sprintf("%d", secondsPerChunk), "-acodec", "libmp3lame", outputFilePath)
	output, err := cmd.CombinedOutput()

	if err != nil {
		errChan <- fmt.Errorf("Error splitting file: %s\nOutput: %s", err.Error(), string(output))
		return
	}
	wg.Done()
	errChan <- nil
}
