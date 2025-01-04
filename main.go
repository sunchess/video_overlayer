package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
	"video_updater/src/config"
	"video_updater/src/process"
)

var AppConfig = &config.ConfigStruct{}

func main() {
	AppConfig = config.NewConfig()

	defer AppConfig.DB.Close()

	directories, err := postDirectories()
	if err != nil {
		log.Fatal(err)
	}

	if len(directories) == 0 {
		log.Println("No directories to process")
		return
	}

	// split directories into groups
	var directory_groups [][]string
	for i := 0; i < len(directories); i += AppConfig.GroupForProcessing {
		end := i + AppConfig.GroupForProcessing
		if end > len(directories) {
			end = len(directories)
		}
		directory_groups = append(directory_groups, directories[i:end])
	}

	directory_pool := make(chan string, AppConfig.WorkerCount)
	results := make(chan string, AppConfig.WorkerCount)

	for worker := 1; worker <= AppConfig.WorkerCount; worker++ {
		go processingWorker(worker, directory_pool, results)
	}

	defer close(directory_pool)
	defer close(results)

	for idx := range directory_groups {
		work_directories := directory_groups[idx]

		for _, directory := range work_directories {
			directory_pool <- directory
		}

		for i := 0; i < len(work_directories); i++ {
			result := <-results
			log.Println(result)
		}

		time.Sleep(5 * time.Second)
	}
}

// ********* private functions *********

func processingWorker(id int, directory_pool chan string, results chan string) error {
	for directory := range directory_pool {
		pr := process.NewProcess(directory, AppConfig)
		pr.WorkerInvoke(id, results)
	}
	return nil
}

func postDirectories() ([]string, error) {
	// get all files in the directory
	files, err := getPostDirectories()
	if err != nil {
		log.Fatal(err)
	}

	processed_paths, err := AppConfig.DB.GetProcessedPaths()
	if err != nil {
		log.Fatal(err)
	}

	// exclude files that have already been processed
	not_processed_paths := files

	for _, db_path := range processed_paths {
		for idx, path := range not_processed_paths {
			if path == db_path {
				not_processed_paths = append(not_processed_paths[:idx], not_processed_paths[idx+1:]...)
			}
		}
	}

	return not_processed_paths, nil
}

func getPostDirectories() ([]string, error) {
	files, err := os.ReadDir(AppConfig.PostsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read posts directory: %v", err)
	}

	var postDirs []string

	for _, file := range files {
		if file.IsDir() && file.Name() != "." && file.Name() != ".." {
			postDirs = append(postDirs, filepath.Join(AppConfig.PostsDir, file.Name()))
		}
	}

	return postDirs, nil
}
