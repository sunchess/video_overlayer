package config

import (
	"log"
	"os"
	"path/filepath"
	"video_updater/src/database"

	"github.com/joho/godotenv"
)

const APP_NAME = "video_overlayer"

type ConfigStruct struct {
	AppName            string
	DbPath             string
	LogsPath           string
	WorkerCount        int
	GroupForProcessing int
	PostsDir           string
	ProcessedVideoName string
	DB                 database.DBConnection
}

func NewConfig() *ConfigStruct {
	loadDotEnv()
	logs_dir := os.Getenv("LOGS_PATH")
	initLogs(logs_dir)

	db_path := os.Getenv("DB_PATH")

	// connect to database
	db, err := database.NewDBConnection(db_path)

	if err != nil {
		log.Fatal(err)
	}

	return &ConfigStruct{
		AppName:            APP_NAME,
		DbPath:             db_path,
		LogsPath:           filepath.Join(logs_dir, APP_NAME+".log"),
		WorkerCount:        3, //depends on the server capacity
		GroupForProcessing: 9,
		PostsDir:           os.Getenv("POSTS_DIR"),
		ProcessedVideoName: "processed.mp4",
		DB:                 *db,
	}
}

func initLogs(logsDir string) {
	if logsDir == "" {
		return
	}

	log_file_path := filepath.Join(logsDir, APP_NAME+".log")
	logFile, err := os.OpenFile(log_file_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	// Set log output to the file
	log.SetOutput(logFile)
}

func loadDotEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
