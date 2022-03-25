package log

import (
	formatter "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

var Logger *log.Logger
var ChangesLogger *log.Logger

func CreateLoggers() {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	fileChanges, err := os.OpenFile("changes.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	stdAndFile := io.MultiWriter(file, os.Stdout)

	stdAndFileChanges := io.MultiWriter(fileChanges, os.Stdout)

	Logger = log.New()
	Logger.SetOutput(stdAndFile)
	Logger.SetFormatter(&formatter.Formatter{})
	ChangesLogger = log.New()
	ChangesLogger.SetOutput(stdAndFileChanges)
	ChangesLogger.SetFormatter(&formatter.Formatter{})

}
