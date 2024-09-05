package splitter

import (
	"fmt"

	"strings"
	"sync"

	"github.com/D4rkP1xel/media-file-splitter/utils"
)

func SplitMediaFile(secondsPerChunk int, inputFilePath string, outputDirectoryPath string, createFolderIfNotExists ...bool) ([]string, error) {
	fileData, err := utils.HandleParams(createFolderIfNotExists, secondsPerChunk, outputDirectoryPath, inputFilePath)
	if err != nil {
		return nil, fmt.Errorf("%s\n", err)
	}
	// Calculate number of chunks
	var numChunks uint32 = fileData.Duration / uint32(secondsPerChunk)
	if fileData.Duration%uint32(secondsPerChunk) != 0 {
		numChunks++
	}

	//grab the input file name and type
	fileName := strings.Split(fileData.InputFileInfo.Name(), ".")[0]
	fileType := strings.Split(inputFilePath, ".")[1]

	var wg sync.WaitGroup
	errChan := make(chan error, numChunks)
	// Split the file into chunks using FFmpeg
	var i uint32
	outputFilesPaths := make([]string, numChunks, numChunks)
	for i = 0; i < numChunks; i++ {
		wg.Add(1)
		outputFilePath := fmt.Sprintf("%s/%s_%04d.%s", outputDirectoryPath, fileName, i+1, fileType)
		outputFilesPaths[i] = outputFilePath
		startTime := i * uint32(secondsPerChunk)
		go utils.GenerateChunk(float64(startTime), uint16(secondsPerChunk), inputFilePath, outputFilePath, errChan, &wg)
	}

	err = utils.HandleCloseChannel(errChan, uint16(numChunks), &wg)
	if err != nil {
		return nil, fmt.Errorf("%s\n", err)
	}

	return outputFilesPaths, nil
}

func SplitMediaFileByStartTimePos(secondsPerChunk int, numChunksToSplit int,
	startPosInSec float64, inputFilePath string,
	outputDirectoryPath string, createFolderIfNotExists ...bool) ([]string, error) {

	fileData, err := utils.HandleParams(createFolderIfNotExists, secondsPerChunk, outputDirectoryPath, inputFilePath)
	if err != nil {
		return nil, fmt.Errorf("%s\n", err)
	}

	if startPosInSec > float64(fileData.Duration) {
		return nil, fmt.Errorf("StartPosInSec cannot be bigger than file duration")
	}

	// Calculate number of chunks
	durationLeft := float64(fileData.Duration) - startPosInSec
	var numChunksLeft uint16 = uint16(durationLeft/float64(secondsPerChunk)) + 1

	if uint32(numChunksToSplit) > uint32(numChunksLeft) {
		numChunksToSplit = int(numChunksLeft)
	}

	//grab the input file name and type
	fileName := strings.Split(fileData.InputFileInfo.Name(), ".")[0]
	fileType := strings.Split(inputFilePath, ".")[1]
	var wg sync.WaitGroup

	errChan := make(chan error, numChunksToSplit)
	// Split the file into chunks using FFmpeg
	var i uint32
	outputFilesPaths := make([]string, numChunksToSplit, numChunksToSplit)

	for i = 0; i < uint32(numChunksToSplit); i++ {
		wg.Add(1)
		outputFilePath := fmt.Sprintf("%s/%s_%04d.%s", outputDirectoryPath, fileName, i+1, fileType)
		outputFilesPaths[i] = outputFilePath
		var startTime float64 = float64(i)*float64(secondsPerChunk) + startPosInSec
		go utils.GenerateChunk(startTime, uint16(secondsPerChunk), inputFilePath, outputFilePath, errChan, &wg)
	}

	err = utils.HandleCloseChannel(errChan, uint16(numChunksToSplit), &wg)
	if err != nil {
		return nil, fmt.Errorf("%s\n", err)
	}

	return outputFilesPaths, nil
}

func SplitMediaFileByStartChunkIndex(secondsPerChunk int, numChunksToSplit int,
	startChunkIndex int, inputFilePath string,
	outputDirectoryPath string, createFolderIfNotExists ...bool) ([]string, error) {

	fileData, err := utils.HandleParams(createFolderIfNotExists, secondsPerChunk, outputDirectoryPath, inputFilePath)
	if err != nil {
		return nil, fmt.Errorf("%s\n", err)
	}

	// Calculate number of chunks
	var numChunks uint32 = fileData.Duration / uint32(secondsPerChunk)
	if fileData.Duration%uint32(secondsPerChunk) != 0 {
		numChunks++
	}
	if startChunkIndex < 0 {
		return nil, fmt.Errorf("Starting chunk index cannot be a negative number")
	}
	if (startChunkIndex + 1) > int(numChunks) {
		return nil, fmt.Errorf("Starting chunk cannot be bigger than total num chunks.\nStarting chunk index: %d\nTotal num chunks: %d\n", startChunkIndex, numChunks)
	}

	numChunksLeft := numChunks - (uint32(startChunkIndex))

	if numChunksToSplit > int(numChunksLeft) {
		numChunksToSplit = int(numChunksLeft)
	}

	//grab the input file name and type
	fileName := strings.Split(fileData.InputFileInfo.Name(), ".")[0]
	fileType := strings.Split(inputFilePath, ".")[1]
	var wg sync.WaitGroup

	errChan := make(chan error, numChunksToSplit) // Split the file into chunks using FFmpeg
	var i uint32
	outputFilesPaths := make([]string, numChunksToSplit, numChunksToSplit)

	for i = 0; i < uint32(numChunksToSplit); i++ {
		wg.Add(1)
		outputFilePath := fmt.Sprintf("%s/%s_%04d.%s", outputDirectoryPath, fileName, i+1+uint32(startChunkIndex), fileType)
		outputFilesPaths[i] = outputFilePath
		startTime := (i + uint32(startChunkIndex)) * uint32(secondsPerChunk)
		go utils.GenerateChunk(float64(startTime), uint16(secondsPerChunk), inputFilePath, outputFilePath, errChan, &wg)
	}

	err = utils.HandleCloseChannel(errChan, uint16(numChunksToSplit), &wg)
	if err != nil {
		return nil, fmt.Errorf("%s\n", err)
	}

	return outputFilesPaths, nil
}
