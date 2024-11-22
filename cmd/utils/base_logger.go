package utils

import (
	"fmt"
	"kg/procurement/cmd/config"
	"log"
	"os"
	"path/filepath"

	"github.com/newrelic/go-agent/v3/integrations/logcontext-v2/logWriter"
	"github.com/newrelic/go-agent/v3/newrelic"
)

var errorLogger *log.Logger
var debugLogger *log.Logger
var infoLogger *log.Logger
var fatalLogger *log.Logger

func InitLogger(cfg config.NewRelic, nrApp *newrelic.Application) {
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

	if cfg.Enabled {
		writer := logWriter.New(os.Stdout, nrApp)

		errorLogger = log.New(&writer, "", log.Default().Flags())
		debugLogger = log.New(&writer, "", log.Default().Flags())
		infoLogger = log.New(&writer, "", log.Default().Flags())
		fatalLogger = log.New(&writer, "", log.Default().Flags())
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
