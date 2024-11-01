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
	absPath, err := filepath.Abs(`\ppl\kompas-gramedia\be\log`)
	if err != nil {
		fmt.Println("Error reading abs path: ", err)
	}
	generalLog, err := os.OpenFile(absPath+"/general-log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("error opening file: ", err)
		os.Exit(1)
	}
	GeneralLogger = log.New(generalLog, "General Logger:\t", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(generalLog, "Error:\t", log.Ldate|log.Ltime|log.Lshortfile)
	DebugLogger = log.New(generalLog, "Debug:\t", log.Ldate|log.Ltime|log.Lshortfile)
	InfoLogger = log.New(generalLog, "Info:\t", log.Ldate|log.Ltime|log.Lshortfile)
	PanicLogger = log.New(generalLog, "Panic:\t", log.Ldate|log.Ltime|log.Lshortfile)
	FatalLogger = log.New(generalLog, "Fatal:\t", log.Ldate|log.Ltime|log.Lshortfile)
}
