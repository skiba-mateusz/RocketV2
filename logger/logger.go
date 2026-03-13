package logger

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
)

type LogLevel int

const (
	SUCCESS LogLevel = iota
	INFO
	WARN
	ERROR
	DEBUG
)

type Logger struct {
	logger *log.Logger
	colorFuncs map[LogLevel]func(a ...interface{}) string
}

func NewLogger() *Logger {
	colorFuncs := map[LogLevel]func(a ...interface{}) string {
		SUCCESS: color.New(color.FgGreen).SprintFunc(),
		INFO: color.New(color.FgBlue).SprintFunc(),
		WARN: color.New(color.FgYellow).SprintFunc(),
		ERROR: color.New(color.FgRed).SprintFunc(),
		DEBUG: color.New(color.FgMagenta).SprintFunc(),
	}

	return &Logger{
		logger: log.New(os.Stdout, "", 0),
		colorFuncs: colorFuncs,
	}
}

func (l *Logger) Success(message string, args ...interface{}) {
	l.logMessage(SUCCESS, "SUCCESS", message, args...)
}

func (l *Logger) Info(message string, args ...interface{}) {
	l.logMessage(INFO, "INFO", message, args...)
}

func (l *Logger) Warn(message string, args ...interface{}) {
	l.logMessage(WARN, "WARN", message, args...)
}

func (l *Logger) Error(message string, args ...interface{}) {
	l.logMessage(ERROR, "ERROR", message, args...)
}

func (l *Logger) Debug(message string, args ...interface{}) {
	l.logMessage(DEBUG, "DEBUG", message, args...)
} 

func (l *Logger) logMessage(level LogLevel, levelStr, message string, args ...interface{}) {
	var formattedMessage string
	if len(args) > 0 {
		formattedMessage = fmt.Sprintf(message	, args...)
	} else {
		formattedMessage = message
	}

	coloredLeveL := l.colorFuncs[level](levelStr)
	l.logger.Printf("[%s] %s", coloredLeveL, formattedMessage)
}