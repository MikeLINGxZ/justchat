package logger

import (
	"path/filepath"

	"github.com/fatih/color"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/utils"
)

var staticLog *Logger

func init() {
	staticLog = &Logger{
		errorLogger: color.New(color.FgRed),
		warmLogger:  color.New(color.FgYellow),
		debugLogger: color.New(color.FgCyan),
		infoLogger:  color.New(color.FgMagenta),
		panicLogger: color.New(color.FgRed),
	}
	staticLog.logger = staticLog.standerLog
}

func NewStaticLogger(name string, options ...Options) error {
	errorLogger := color.New(color.FgRed)
	warmLogger := color.New(color.FgYellow)
	normalLogger := color.New(color.FgCyan)
	debugLogger := color.New(color.FgMagenta)
	panicLogger := color.New(color.FgRed)
	logger := &Logger{
		loggerName:  name,
		errorLogger: errorLogger,
		warmLogger:  warmLogger,
		debugLogger: debugLogger,
		infoLogger:  normalLogger,
		panicLogger: panicLogger,
	}
	logger.logger = logger.standerLog

	dataPath, err := utils.GetDataPath()
	if err != nil {
		return err
	}
	logDir := filepath.Join(dataPath, "logs")
	options = append(options, WithFileLogging(logDir))

	for _, option := range options {
		err := option(logger)
		if err != nil {
			return err
		}
	}
	staticLog = logger
	return nil
}

func Info(any ...any) {
	staticLog.Info(any...)
}

func Warm(any ...any) {
	staticLog.Warm(any...)
}

func Error(any ...any) {
	staticLog.Error(any...)
}

func Debug(any ...any) {
	staticLog.Debug(any...)
}

func Panic(any ...any) {
	staticLog.Panic(any...)
}

func Infof(content string, any ...any) {
	staticLog.Infof(content, any...)
}

func Warmf(content string, any ...any) {
	staticLog.Warmf(content, any...)
}

func Errorf(content string, any ...any) {
	staticLog.Errorf(content, any...)
}

func Debugf(content string, any ...any) {
	staticLog.Debugf(content, any...)
}

func Panicf(content string, any ...any) {
	staticLog.Panicf(content, any...)
}
