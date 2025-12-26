package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
)

type Clip struct {
	Name      string `json:"name"`
	StartTime string `json:"startTime"` // formato "HH:MM:SS"
	EndTime   string `json:"endTime"`   // formato "HH:MM:SS"
}

type Video struct {
	Title          string `json:"title"`
	InputVideoPath string `json:"inputVideoPath"`
	Clips          []Clip `json:"clips"`
}

var logger = log.New(os.Stderr, "LOG: ", 0)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <input_file>\n", os.Args[0])
		return
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		logger.Fatal(err)
	}

	var videos []Video

	if err := json.Unmarshal(data, &videos); err != nil {
		logger.Fatal(err)
	}

	ctx := context.Background()

	for _, video := range videos {
		for _, clip := range video.Clips {
			// clipName := strings.ToLower(strings.ReplaceAll(clip.Name, " ", "_"))
			// outputVideoFile := fmt.Sprintf("%s_%s.mp4", strings.ToLower(video.Title), clipName)

			if err := execFFMPEG(ctx, video.InputVideoPath, clip); err != nil {
				logger.Printf("Failed to execute ffmpeg: %s", err.Error())
				break
			}
		}
	}
}

func execFFMPEG(ctx context.Context, inputVideo string, clip Clip) error {
	// ffmpeg -i .\IMG_9032.MOV -ss 00:09:06 -to 00:09:18 -c copy output.mp4

	outputVideoFile := fmt.Sprintf("%s.mp4", clip.Name)

	cmd := exec.CommandContext(ctx, "/usr/bin/ffmpeg", "-hide_banner", "-loglevel", "error", "-y",
		"-i", inputVideo,
		"-ss", clip.StartTime, "-to", clip.EndTime,
		"-c", "copy", outputVideoFile)

	if err := cmd.Start(); err != nil {
		logger.Printf("Failed to start ffmpeg: %s\n", err.Error())
		return err
	}

	if err := cmd.Wait(); err != nil {
		logger.Printf("ffmpeg exited with error: %s\n", err.Error())
		return err
	}

	return nil
}
