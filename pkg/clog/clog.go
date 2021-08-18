package clog

import (
	"log"
	"os"
)

const (
	// log line format settings.
	logSettings = log.Ldate | log.Ltime | log.Lmsgprefix
	// Prefix tag for the info logs.
	logInfoTag = `[INFO]: `
	// Prefix tag for the error logs.
	logErrorTag = `[ERROR]: `
)

// Info and Error loggers for the console and log file.
var (
	infoLogger   *log.Logger
	errorLogger  *log.Logger
	infoFLogger  *log.Logger
	errorFLogger *log.Logger
)

// logFile is the log file reference. It will be initialized
// when the init function runs.
var logFile *os.File

// init will be executed when the application starts, it will initilize
// the loggers and open the log file to write the logs.
func init() {
	infoLogger = log.New(os.Stdout, logInfoTag, logSettings)
	errorLogger = log.New(os.Stderr, logErrorTag, logSettings)
	var err error

	logFile, err = os.CreateTemp(os.TempDir(), "catchmyfile-*.log")
	if err != nil {
		Error(err)
		os.Exit(1)
	}

	infoFLogger = log.New(logFile, logInfoTag, logSettings)
	errorFLogger = log.New(logFile, logErrorTag, logSettings)
}

// Info will log an info message to the stdout and log file.
func Info(msg string, a ...interface{}) {
	infoLogger.Printf(msg, a...)
	infoFLogger.Printf(msg, a...)
}

// Error will log an error to the stderr and log file.
func Error(err error) {
	errorLogger.Println(err)
	errorFLogger.Println(err)
}

// LogFile returns the log file path and name.
func LogFile() string {
	return logFile.Name()
}

// Close will close the log file io.
func Close() error {
	return logFile.Close()
}
