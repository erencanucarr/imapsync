package app

import (
	"fmt"
	"io"
	"os"
	"time"
)

// LogLevel represents the logging level
type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Logger provides a simple logging interface
type Logger struct {
	level LogLevel
	out   io.Writer
}

// NewLogger creates a new logger instance
func NewLogger() *Logger {
	return &Logger{
		level: LevelInfo,
		out:   os.Stdout,
	}
}

// SetLevel sets the logging level
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// SetOutput sets the output writer
func (l *Logger) SetOutput(out io.Writer) {
	l.out = out
}

// log writes a log message if the level is sufficient
func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)

	logLine := fmt.Sprintf("[%s] %s: %s\n", timestamp, level.String(), message)
	l.out.Write([]byte(logLine))
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(LevelDebug, format, args...)
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(LevelInfo, format, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(LevelWarn, format, args...)
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(LevelError, format, args...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(LevelFatal, format, args...)
	os.Exit(1)
}

// Debugf logs a debug message (alias for Debug)
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Debug(format, args...)
}

// Infof logs an info message (alias for Info)
func (l *Logger) Infof(format string, args ...interface{}) {
	l.Info(format, args...)
}

// Warnf logs a warning message (alias for Warn)
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Warn(format, args...)
}

// Errorf logs an error message (alias for Error)
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Error(format, args...)
}

// Fatalf logs a fatal message and exits (alias for Fatal)
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.Fatal(format, args...)
}
