package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

var (
	ffmpegLogger, _ = os.Create("./ffmpeg.log")
	logger          = log.New(ffmpegLogger, "", 0)
)

func main() {
	data, err := os.ReadFile("./data.json")
	if err != nil {
		log.Fatal(err)
	}

	var videos []Video

	if err := json.Unmarshal(data, &videos); err != nil {
		log.Fatal(err)
	}

	if len(videos) == 0 {
		log.Fatal(errors.New("no video to cut"))
	}

	ctx := context.Background()

	for _, video := range videos {
		for _, clip := range video.Clips {
			// clipName := strings.ToLower(strings.ReplaceAll(clip.Name, " ", "_"))
			// outputVideoFile := fmt.Sprintf("%s_%s.mp4", strings.ToLower(video.Title), clipName)

			if err := execFFMPEG(ctx, video.InputVideoPath, clip); err != nil {
				log.Printf("Failed to execute ffmpeg: %s", err.Error())
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

	// stdoutPipe, _ := cmd.StdoutPipe()
	// stderrPipe, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		log.Printf("Failed to start ffmpeg: %s\n", err.Error())
		return err
	}

	// go streamLinesToLogger(stdoutPipe, logger)
	// go streamLinesToLogger(stderrPipe, logger)

	if err := cmd.Wait(); err != nil {
		log.Printf("ffmpeg exited with error: %s\n", err.Error())
		return err
	}

	return nil
}

func streamLinesToLogger(r io.Reader, logger *log.Logger) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		logger.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		logger.Printf("Scanner error: %s\n", err.Error())
	}
}
