package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var GeneralLogger *log.Logger
var ErrorLogger *log.Logger
var DebugLogger *log.Logger
var InfoLogger *log.Logger
var PanicLogger *log.Logger
var FatalLogger *log.Logger

// TODO
// Use builder pattern to create a more descriptive logging functions
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
	GeneralLogger = log.New(generalLog, "General Logger:\t", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(generalLog, "Error:\t", log.Ldate|log.Ltime|log.Lshortfile)
	DebugLogger = log.New(generalLog, "Debug:\t", log.Ldate|log.Ltime|log.Lshortfile)
	InfoLogger = log.New(generalLog, "Info:\t", log.Ldate|log.Ltime|log.Lshortfile)
	PanicLogger = log.New(generalLog, "Panic:\t", log.Ldate|log.Ltime|log.Lshortfile)
	FatalLogger = log.New(generalLog, "Fatal:\t", log.Ldate|log.Ltime|log.Lshortfile)
}
