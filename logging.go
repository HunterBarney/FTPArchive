package main

import (
	"io"
	"log"
	"os"
	"time"
)

func initLogging(config *Config) *os.File {
	logFolder := config.LogDirectory
	if _, err := os.Stat(logFolder); os.IsNotExist(err) {
		err = os.MkdirAll(logFolder, 0755)
		if err != nil {
			log.Fatal("Error creating log folder: ", err)
		}
	}

	logPath := logFolder + "/" + time.Now().Format("01-02-2006_15-04-05_MST") + ".txt"
	logFile, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error opening log file: ", err)
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	return logFile
}
