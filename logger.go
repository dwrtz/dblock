package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

func SetupLogging() {
	logDir := filepath.Join(GetHomeDir(), ".dblock", "logs")
	if err := EnsureDir(logDir); err != nil {
		fmt.Println("Error creating log directory:", err)
		os.Exit(1)
	}
	logPath := filepath.Join(logDir, "dblock.log")
	log.SetOutput(&lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    10, // megabytes
		MaxBackups: 5,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	})

	// Change ownership of the log file
	uid, gid, err := getOriginalUserIDs()
	if err == nil {
		os.Chown(logPath, uid, gid)
	}
}

func logError(err error) {
	log.Println("ERROR:", err)
}
