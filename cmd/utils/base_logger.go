package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var errorLogger *log.Logger
var debugLogger *log.Logger
var infoLogger *log.Logger
var fatalLogger *log.Logger

func init() {
	projectRoot, err := filepath.Abs("./")
	if err != nil {
		fmt.Println("Error finding project root path:", err)
		return
	}

	logFilePath := filepath.Join(projectRoot, "log", "general-log.log")

	logDir := filepath.Dir(logFilePath)
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		fmt.Println("Failed to create log directory:", err)
		return
	}

	generalLog, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("error opening file: ", err)
		return
	}
	errorLogger = log.New(generalLog, "Error:\t", log.Ldate|log.Ltime|log.Lshortfile)
	debugLogger = log.New(generalLog, "Debug:\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLogger = log.New(generalLog, "Info:\t", log.Ldate|log.Ltime|log.Lshortfile)
	fatalLogger = log.New(generalLog, "Fatal:\t", log.Ldate|log.Ltime|log.Lshortfile)
}
