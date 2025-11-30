package shared

import (
	"log"
	"os"
)

var (
	// InfoLogger logs info messages
	InfoLogger *log.Logger
	// ErrorLogger logs error messages
	ErrorLogger *log.Logger
)

func init() {
	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// LogInfo logs an info message
func LogInfo(format string, v ...interface{}) {
	InfoLogger.Printf(format, v...)
}

// LogError logs an error message
func LogError(format string, v ...interface{}) {
	ErrorLogger.Printf(format, v...)

}
