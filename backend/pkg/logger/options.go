package logger

import (
	"os"

	"github.com/fatih/color"
)

type Options func(logger *Logger) error

func WithColor(errorLogger, warmLogger, debugLogger, infoLogger *color.Color) Options {
	return func(logger *Logger) error {
		logger.errorLogger = errorLogger
		logger.warmLogger = warmLogger
		logger.debugLogger = debugLogger
		logger.infoLogger = infoLogger
		return nil
	}
}

func WithDebug() Options {
	return func(logger *Logger) error {
		logger.enableDebug = true
		return nil
	}
}

func WithFileLogging(logDir string) Options {
	return func(logger *Logger) error {
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return err
		}
		logger.logDir = logDir
		return nil
	}
}
