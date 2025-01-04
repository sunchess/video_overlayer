package main

import (
	"log"
	"testing"
	"time"
	"video_updater/src/config"
)

func TestProcessingWorker(t *testing.T) {
	// Setup
	AppConfig = config.NewConfig()
	defer AppConfig.DB.Close()

	directory_pool := make(chan string, AppConfig.WorkerCount)
	results := make(chan string, AppConfig.WorkerCount)

	go func() {
		err := processingWorker(1, directory_pool, results)
		if err != nil {
			t.Fatal(err)
		}
	}()

	// Test data
	testDirectory := "test_directory"
	directory_pool <- testDirectory
	close(directory_pool)

	result := <-results
	expected := "No video file found for processing in " + testDirectory + "/media"
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestMainFunction(t *testing.T) {
	// Setup
	AppConfig = config.NewConfig()
	defer AppConfig.DB.Close()

	// Mock data
	directory_groups := [][]string{
		{"test_directory_1", "test_directory_2"},
	}

	directory_pool := make(chan string, AppConfig.WorkerCount)
	results := make(chan string, AppConfig.WorkerCount)

	for worker := 1; worker <= AppConfig.WorkerCount; worker++ {
		go func(worker int) {
			err := processingWorker(worker, directory_pool, results)
			if err != nil {
				t.Fatal(err)
			}
		}(worker)
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
