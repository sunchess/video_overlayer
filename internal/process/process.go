package process

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"video_updater/internal/config"

	"golang.org/x/exp/rand"
)

type Process struct {
	inputPostPath      string
	mediaDirectory     string
	outroVideoPath     string
	inputVideoPath     string
	processedVideoPath string
	appConfig          *config.ConfigStruct
}

func NewProcess(inputPostPath string, AppConfig *config.ConfigStruct) *Process {
	mediaDirectory := filepath.Join(inputPostPath, "media")

	return &Process{
		inputPostPath:      inputPostPath,
		mediaDirectory:     mediaDirectory,
		outroVideoPath:     getoutroVideoPath(),
		processedVideoPath: filepath.Join(mediaDirectory, AppConfig.ProcessedVideoName),
		appConfig:          AppConfig,
	}
}

func (p *Process) WorkerInvoke(id int, results chan string) {
	log.Printf("Worker %d started", id)
	file, err := getVideoFile(p.mediaDirectory)
	if err != nil {
		log.Printf("Media directory error %e", err)
	}

	p.inputVideoPath = file
	is_processed := false

	_, err = os.Stat(filepath.Join(p.mediaDirectory, p.appConfig.ProcessedVideoName))

	if err == nil {
		is_processed = true
	}

	if file != "" {
		if is_processed {
			log.Printf("File %s already processed", file)
		} else {
			err = p.processVideoUpdate()
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	var output string
	if p.inputVideoPath != "" {
		output = p.processedVideoPath
	} else {
		output = "No video file found for processing in " + p.mediaDirectory
	}

	err = p.appConfig.DB.UpdateProcessed(p.inputPostPath)
	if err != nil {
		log.Fatal(err)
	}

	results <- output
}

func (p *Process) processVideoUpdate() error {
	log.Printf("processing video %s", p.inputVideoPath)
	// add video to the end of the video file by ffmpeg
	args := []string{
		"-i", p.inputVideoPath, "-i", p.outroVideoPath,
		"-filter_complex",
		"[0:v]scale=720:1280,setsar=1[in0];[1:v]scale=720:1280,setsar=1[in1];[in0][0:a][in1][1:a]concat=n=2:v=1:a=1[outv][outa]",
		"-map", "[outv]", "-map", "[outa]",
		p.processedVideoPath,
	}
	cmd := exec.Command("ffmpeg", args...)

	err := cmd.Run()

	if err != nil {
		return err
	}

	return nil
}

func getoutroVideoPath() string {
	outroVideoPath := os.Getenv("OUTRO_VIDEOS_PATH")
	if outroVideoPath == "" {
		log.Fatal("OUTRO_VIDEOS_PATH is not set")
	}

	files, err := os.ReadDir(outroVideoPath)

	if err != nil {
		log.Fatal(err)
	}

	allFiles := []string{}

	for _, file := range files {
		if isVideoFile(file.Name()) {
			allFiles = append(allFiles, file.Name())
		}
	}

	randomInt := rand.Intn(len(allFiles))

	return filepath.Join(outroVideoPath, allFiles[randomInt])
}

func getVideoFile(mediaDirectory string) (string, error) {
	files, err := os.ReadDir(mediaDirectory)
	if err != nil {
		return "", err
	}

	var videoFile string

	for _, file := range files {
		if isVideoFile(file.Name()) {
			videoFile = filepath.Join(mediaDirectory, file.Name())
		}
	}

	return videoFile, nil
}

func isVideoFile(file string) bool {
	file = strings.ToLower(file)
	return filepath.Ext(file) == ".mp4"
}
