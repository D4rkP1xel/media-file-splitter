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
## 1. Import the package

<br/>

`import "github.com/D4rkP1xel/media-file-splitter/splitter"`

<br/>


## 2. Functions

### 2.1 SplitMediaFile

Splits the media file into **all** possible segments of `secondsPerChunk` size.

**Example:** A 300 second audio file split in 50 second segments will result in 6 chunks.

<br/>
   
`SplitMediaFileByTimedChunks(secondsPerChunk, inputFilePath, outputDirectoryPath, ...createFolderIfNotExists) ([]string, error)`

<br/>

**secondsPerChunk** \<int>: How much time (in seconds) each chunk should have.

**inputFilePath** \<string>: Path to the input media file.

**outputDirectoryPath** \<string>: Path to directory where the chunks will be stored.

**createFolderIfNotExist** \<bool>(optional): Whether to create the output directory if it does not exist. Default is false.

<br/>

**returns**:

**outputFilePaths** \<[]string>: Array with the paths to the newly created chunks.

**fileData** \<FileData>: Struct with file info like file duration.

**error** \<error>: Error return in case something goes wrong.

<br/>

### Example

`chunkPaths, err := splitter.SplitMediaFileByTimedChunks(30, "/path/to/input/folder/input.mp3", "/path/to/output/folder", true)`

<br/>

### 2.2 SplitMediaFileByStartChunkIndex

Splits the media file into `numChunksToSplit` segments of `secondsPerChunk` size starting at the chunk at index `startChunkIndex`.
The media file will only be split until there are no more chunks to split. In case you want all chunks starting from `startChunkIndex` index to the last chunk, make `numChunksToSplit` a big number (ex: 9999).

**Example:** A 300 second audio file, starting at index 2, with three 50 second chunks to split, will split the file at chunk index 2, 3, and 4.

<br/>

### Example
   
`SplitMediaFileByStartChunkIndex(secondsPerChunk, numChunksToSplit, startChunkIndex, inputFilePath, outputDirectoryPath, ...createFolderIfNotExists) ([]string, error)`

<br/>

### 2.3 SplitMediaFileByStartTimePos

Split the media file into `numChunksToSplit` segments of `secondsPerChunk` size starting at timestamp `startPosInSec`.
The media file will only be split until there are no more chunks to split. In case you want all chunks starting from `startChunkIndex` index to the last chunk, make `numChunksToSplit` a big number (ex: 9999).

**Example:** A 300 second audio file, starting at second 25, with three 60 second chunks to split, will split the file at chunk 00:25-01:25, 01:25-02:25 and 02:25-03:25.

<br/>

### Example
   
`SplitMediaFileByStartTimePos(secondsPerChunk, numChunksToSplit, startPosInSec, inputFilePath, outputDirectoryPath, ...createFolderIfNotExists) ([]string, error)`

<br/>


## Additional Info

1. Only mp3 and mp4 files were tested.
2. Feel free to use, modify and/or contribute to this repository.
