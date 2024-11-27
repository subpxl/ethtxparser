// Package logger provides a simple leveled logging system with file output
package logger

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "runtime"
    "time"
)


// LogLevel represents different severity levels for logging
type LogLevel int

// Logging levels from least to most severe
const (
    DEBUG LogLevel = iota
    INFO
    WARN
    ERROR
    FATAL
)

// Logger contains separate loggers for each severity level
type Logger struct {
    debug *log.Logger
    info  *log.Logger
    warn  *log.Logger
    error *log.Logger
    fatal *log.Logger
}

var (
    defaultLogger *Logger
    logFile      *os.File
)

// Init initializes the logging system with the specified log file path
func Init(logPath string) error {
    if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
        return fmt.Errorf("failed to create log directory: %v", err)
    }

    file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return fmt.Errorf("failed to open log file: %v", err)
    }

    logFile = file
    defaultLogger = &Logger{
        debug: log.New(file, "DEBUG: ", log.Ldate|log.Ltime),
        info:  log.New(file, "INFO:  ", log.Ldate|log.Ltime),
        warn:  log.New(file, "WARN:  ", log.Ldate|log.Ltime),
        error: log.New(file, "ERROR: ", log.Ldate|log.Ltime),
        fatal: log.New(file, "FATAL: ", log.Ldate|log.Ltime),
    }

    return nil
}

// Close properly closes the log file
func Close() {
    if logFile != nil {
        logFile.Close()
    }
}

// getCallerInfo retrieves the file name and line number of the logging call
func getCallerInfo() string {
    _, file, line, ok := runtime.Caller(2)
    if !ok {
        return "unknown"
    }
    return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}

// logMessage formats and writes a log message with the specified level.
func logMessage(level LogLevel, format string, v ...interface{}) {
    if defaultLogger == nil {
        return
    }

    msg := fmt.Sprintf(format, v...)
    caller := getCallerInfo()
    timestamp := time.Now().Format("2006-01-02 15:04:05")

    logLine := fmt.Sprintf("[%s] %s - %s", timestamp, caller, msg)

    switch level {
    case DEBUG:
        defaultLogger.debug.Println(logLine)
    case INFO:
        defaultLogger.info.Println(logLine)
    case WARN:
        defaultLogger.warn.Println(logLine)
    case ERROR:
        defaultLogger.error.Println(logLine)
    case FATAL:
        defaultLogger.fatal.Println(logLine)
        os.Exit(1)
    }
}

// Debug logs debug level messages
func Debug(format string, v ...interface{}) {
    logMessage(DEBUG, format, v...)
}

// Info logs informational messages
func Info(format string, v ...interface{}) {
    logMessage(INFO, format, v...)
}

// Warn logs warning messages
func Warn(format string, v ...interface{}) {
    logMessage(WARN, format, v...)
}

// Error logs error messages
func Error(format string, v ...interface{}) {
    logMessage(ERROR, format, v...)
}

// Fatal logs critical errors and terminates the program
func Fatal(format string, v ...interface{}) {
    logMessage(FATAL, format, v...)
}