package goNest

import (
	"encoding/json"
	"fmt"
	"github.com/abema/go-mp4"
	"github.com/sunfish-shogi/bufseekio"
	"io"
	"os"
	"strings"
)

// FormatInfo : Format Structure
type FormatInfo struct {
	TagTypes         []string    `json:"tagTypes"`
	TrackInfo        []TrackInfo `json:"trackInfo"`
	FileSize         int64       `json:"file_size"`
	FileSizeMB       float64     `json:"file_size_mb"`
	TracksCount      int         `json:"tracks_count"`
	DurationRaw      int         `json:"duration_raw"`
	Duration         float64     `json:"duration"`
	NumberOfChannels int         `json:"numberOfChannels"`
	Checksum         string      `json:"checksum"`
}

type AudioInfo struct {
	SamplingFrequency int `json:"samplingFrequency"`
	BitDepth          int `json:"bitDepth"`
	Channels          int `json:"channels"`
}

type TrackInfo struct {
	Duration          int       `json:"duration"`
	DurationInSeconds float64   `json:"duration_in_seconds"`
	Audio             AudioInfo `json:"audio"`
}

type CommonInfo struct {
	EncodedBy string `json:"encodedby"`
}

// ProgressWriter tracks the progress of the download
type ProgressWriter struct {
	Writer        io.Writer
	ContentId     string
	FileExtension string
	ChannelId     string
	ClientId      string
	Total         int64
	Downloaded    int64
	LastProgress  int
}

func ExtractMetaData(filePath string) ([]byte, error) {
	// Open MP4 file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, nil
		{
		}
	}
	defer file.Close()

	// Get file size
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return nil, nil
	}
	fileSize := fileInfo.Size()
	fileSizeMB := float64(fileSize) / (1024 * 1024)

	//Generate checksum for the file
	checksum, _ := GenerateChecksum(file)

	// Get basic mp4 information
	info, err := mp4.Probe(bufseekio.NewReadSeeker(file, 1024, 4))
	if err != nil {
		return nil, nil
	}

	//Apply track count info
	trackCount := len(info.Tracks)

	// Find the duration from the mvhd box

	// Variable to hold duration and timescale
	durationRaw := 0.0
	timeScale := 0.0
	durationInSeconds := 0.0

	var trackInfos []TrackInfo

	i := 0
	for _, box := range info.Tracks {
		//fmt.Println(box.Duration)
		//fmt.Println(box.Timescale)

		//Get main duration from the second track as usually the first track info is empty
		if i == 1 {
			//log.Println(box.Codec)
			durationRaw = float64(box.Duration)
			timeScale = float64(box.Timescale)
			durationInSeconds = durationRaw / timeScale
		}

		//For each track
		trackInfo := TrackInfo{
			Duration:          int(box.Duration),
			DurationInSeconds: float64(box.Duration) / float64(box.Timescale),
			Audio: AudioInfo{
				SamplingFrequency: int(box.Timescale),
				BitDepth:          16,
				Channels:          2,
			},
		}
		trackInfos = append(trackInfos, trackInfo)
		i++
	}

	//Finally, construct as structure
	var tagTypes []string
	tagTypes = append(tagTypes, "iTunes")

	formatInfo := FormatInfo{
		TagTypes:         tagTypes,
		TrackInfo:        trackInfos,
		FileSize:         fileSize,
		FileSizeMB:       fileSizeMB,
		TracksCount:      trackCount,
		DurationRaw:      int(durationRaw),
		Duration:         durationInSeconds,
		NumberOfChannels: 2,
		Checksum:         checksum,
	}

	data := struct {
		Format FormatInfo `json:"format"`
		Common CommonInfo `json:"common"`
	}{
		Format: formatInfo,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
		return nil, nil
	}

	//fmt.Printf("File size (Bytes): %d\n", fileSize)
	//fmt.Printf("File size (MB): %f\n", fileSizeMB)
	//
	//fmt.Printf("Number of tracks: %d\n", trackCount)
	//fmt.Printf("Duration: %d \n", duration)
	//fmt.Printf("TimeScale: %d \n", timeScale)
	//fmt.Printf("Duration (seconds): %f \n", durationInSeconds)

	return jsonData, nil
}

// Windows path starts with \\ and linux on /
func GetDirectoryPath(fullPath string) string {
	// Find the last occurrence of the path separator
	lastSlash := strings.LastIndex(fullPath, "/")
	if lastSlash == -1 {
		// No directory part, return empty
		lastSlash = strings.LastIndex(fullPath, "\\")
		if lastSlash == -1 {
			return ""
		}
	}
	return fullPath[:lastSlash]
}
