# What is it?
**media-file-splitter** is a simple Go package designed to split audio/video files into chunks.

Useful for data streaming services where media files have to be delivered in small segments.

<br/>

# Dependencies
Ensure that ffmpeg is installed on your system. Installation commands may vary depending on your operating system.

<br/>

`sudo apt-get install ffmpeg`

<br/>

# How to use
### 1. Import the package

<br/>

`import "github.com/D4rkP1xel/media-file-splitter/splitter"`

<br/>


### 2. Call the function

<br/>
   
`splitter.SplitMediaFileByTimedChunks(secondsPerChunk, inputFilePath, outputDirectoryPath, ...createFolderIfNotExists) ([]string, error)`

<br/>

**secondsPerChunk** \<int>: How much time (in seconds) each chunk should have, except the last one.

**inputFilePath** \<string>: Path to the input media file.

**outputDirectoryPath** \<string>: Path to directory where the chunks will be stored.

**createFolderIfNotExist** \<bool>(optional): Whether to create the output directory if it does not exist. Default is false.

<br/>

**returns**:

**outputFilePaths** \<[]string>: Array with the paths to the newly created chunks.

**error** \<error>: Error return in case something goes wrong.

<br/>

# Example

`chunkPaths, err := splitter.SplitMediaFileByTimedChunks(30, "/path/to/input/folder/input.mp3", "/path/to/output/folder", true)`

<br/>

<br/>

## Additional Info

1. Only mp3 and mp4 files were tested.
2. Feel free to use, modify and/or contribute to this repository.
