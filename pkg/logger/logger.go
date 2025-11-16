package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
	FATAL
)

var levelNames = map[LogLevel]string{
	DEBUG:   "DEBUG",
	INFO:    "INFO",
	WARNING: "WARNING",
	ERROR:   "ERROR",
	FATAL:   "FATAL",
}

var levelColors = map[LogLevel]string{
	DEBUG:   "\033[36m",
	INFO:    "\033[32m",
	WARNING: "\033[33m",
	ERROR:   "\033[31m",
	FATAL:   "\033[35m",
}

const colorReset = "\033[0m"

type Logger struct {
	debugLogger   *log.Logger
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
	fatalLogger   *log.Logger
	level         LogLevel
	useColors     bool
	logFile       *os.File
}

var defaultLogger *Logger

func init() {
	defaultLogger = New(os.Stdout, INFO, true)
}

func New(output io.Writer, level LogLevel, useColors bool) *Logger {
	return &Logger{
		debugLogger:   log.New(output, "", 0),
		infoLogger:    log.New(output, "", 0),
		warningLogger: log.New(output, "", 0),
		errorLogger:   log.New(output, "", 0),
		fatalLogger:   log.New(output, "", 0),
		level:         level,
		useColors:     useColors,
	}
}

func NewWithFile(logFilePath string, level LogLevel, useColors bool) (*Logger, error) {
	dir := filepath.Dir(logFilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	multiWriter := io.MultiWriter(os.Stdout, file)

	logger := New(multiWriter, level, useColors)
	logger.logFile = file

	return logger, nil
}

func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

func (l *Logger) GetLevel() LogLevel {
	return l.level
}

func (l *Logger) Close() error {
	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}

func (l *Logger) formatMessage(level LogLevel, message string) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	levelStr := levelNames[level]

	_, file, line, ok := runtime.Caller(3)
	caller := "unknown"
	if ok {
		caller = fmt.Sprintf("%s:%d", filepath.Base(file), line)
	}

	if l.useColors {
		color := levelColors[level]
		return fmt.Sprintf("%s [%s%s%s] [%s] %s",
			timestamp, color, levelStr, colorReset, caller, message)
	}

	return fmt.Sprintf("%s [%s] [%s] %s", timestamp, levelStr, caller, message)
}

func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	message := fmt.Sprintf(format, args...)
	formattedMessage := l.formatMessage(level, message)

	switch level {
	case DEBUG:
		l.debugLogger.Println(formattedMessage)
	case INFO:
		l.infoLogger.Println(formattedMessage)
	case WARNING:
		l.warningLogger.Println(formattedMessage)
	case ERROR:
		l.errorLogger.Println(formattedMessage)
	case FATAL:
		l.fatalLogger.Println(formattedMessage)
		os.Exit(1)
	}
}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

func (l *Logger) Warning(format string, args ...interface{}) {
	l.log(WARNING, format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(FATAL, format, args...)
}

func (l *Logger) LogHTTPRequest(method, path string, statusCode int, duration time.Duration) {
	statusColor := "\033[32m"
	if statusCode >= 400 {
		statusColor = "\033[31m"
	} else if statusCode >= 300 {
		statusColor = "\033[33m"
	}

	if l.useColors {
		l.Info("%s %s - %s%d%s - %v", method, path, statusColor, statusCode, colorReset, duration)
	} else {
		l.Info("%s %s - %d - %v", method, path, statusCode, duration)
	}
}

func SetDefaultLogger(logger *Logger) {
	defaultLogger = logger
}

func GetDefaultLogger() *Logger {
	return defaultLogger
}

func SetLevel(level LogLevel) {
	defaultLogger.SetLevel(level)
}

func Debug(format string, args ...interface{}) {
	defaultLogger.Debug(format, args...)
}

func Info(format string, args ...interface{}) {
	defaultLogger.Info(format, args...)
}

func Warning(format string, args ...interface{}) {
	defaultLogger.Warning(format, args...)
}

func Error(format string, args ...interface{}) {
	defaultLogger.Error(format, args...)
}

func Fatal(format string, args ...interface{}) {
	defaultLogger.Fatal(format, args...)
}

func LogHTTPRequest(method, path string, statusCode int, duration time.Duration) {
	defaultLogger.LogHTTPRequest(method, path, statusCode, duration)
}

func ParseLevel(levelStr string) LogLevel {
	switch strings.ToUpper(levelStr) {
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "WARNING", "WARN":
		return WARNING
	case "ERROR":
		return ERROR
	case "FATAL":
		return FATAL
	default:
		return INFO
	}
}
