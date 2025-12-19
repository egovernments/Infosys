package logger

import (
	"log"
	"os"
)

// InfoLogger logs informational messages to stdout
// WarnLogger logs warning messages to stdout
// ErrorLogger logs error messages to stderr
// DebugLogger logs debug messages to stdout
var (
	InfoLogger  *log.Logger
	WarnLogger  *log.Logger
	ErrorLogger *log.Logger
	DebugLogger *log.Logger
)

// init initializes the loggers with appropriate prefixes and output streams
func init() {
	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarnLogger = log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	DebugLogger = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// // Init initializes the logger (no-op, initialization happens in init())
// func Init() {
// 	// Logger is already initialized in init()
// }

// Info logs an informational message
func Info(v ...interface{}) {
	InfoLogger.Println(v...)
}

// Infof logs a formatted informational message
func Infof(format string, v ...interface{}) {
	InfoLogger.Printf(format, v...)
}

// Warn logs a warning message
func Warn(v ...interface{}) {
	WarnLogger.Println(v...)
}

// Warnf logs a formatted warning message
func Warnf(format string, v ...interface{}) {
	WarnLogger.Printf(format, v...)
}

// Error logs an error message
func Error(v ...interface{}) {
	ErrorLogger.Println(v...)
}

// Errorf logs a formatted error message
func Errorf(format string, v ...interface{}) {
	ErrorLogger.Printf(format, v...)
}

// Debug logs a debug message
func Debug(v ...interface{}) {
	DebugLogger.Println(v...)
}

// Debugf logs a formatted debug message
func Debugf(format string, v ...interface{}) {
	DebugLogger.Printf(format, v...)
}

// Fatal logs a fatal error and exits the application
func Fatal(v ...interface{}) {
	ErrorLogger.Fatal(v...)
}

// Fatalf logs a formatted fatal error and exits the application
func Fatalf(format string, v ...interface{}) {
	ErrorLogger.Fatalf(format, v...)
}
