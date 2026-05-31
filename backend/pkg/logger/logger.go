package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	debugLevel = iota
	infoLevel
	warnLevel
	errorLevel
)

// Logger writes application log lines into date-split files inside the active data directory.
type Logger struct {
	dataDir  string
	level    int
	now      func() time.Time
	fileMode os.FileMode
	dirMode  os.FileMode
	mu       sync.Mutex
}

// Options configures a Logger instance.
type Options struct {
	DataDir  string
	Level    string
	Now      func() time.Time
	FileMode os.FileMode
	DirMode  os.FileMode
}

// New creates a logger that writes into {dataDir}/logs/YYYY-MM-DD.log.
func New(options Options) (*Logger, error) {
	level, err := parseLevel(options.Level)
	if err != nil {
		return nil, err
	}
	if options.Now == nil {
		options.Now = time.Now
	}
	if options.FileMode == 0 {
		options.FileMode = 0o644
	}
	if options.DirMode == 0 {
		options.DirMode = 0o755
	}

	return &Logger{
		dataDir:  options.DataDir,
		level:    level,
		now:      options.Now,
		fileMode: options.FileMode,
		dirMode:  options.DirMode,
	}, nil
}

// Debug writes a debug log line when the configured level allows it.
func (l *Logger) Debug(message string) {
	_ = l.write("DEBUG", debugLevel, message)
}

// Info writes an info log line when the configured level allows it.
func (l *Logger) Info(message string) {
	_ = l.write("INFO", infoLevel, message)
}

// Warn writes a warning log line when the configured level allows it.
func (l *Logger) Warn(message string) {
	_ = l.write("WARN", warnLevel, message)
}

// Error writes an error log line when the configured level allows it.
func (l *Logger) Error(message string) {
	_ = l.write("ERROR", errorLevel, message)
}

// write appends a single line to the current daily log file when the message passes filtering.
func (l *Logger) write(label string, level int, message string) error {
	if level < l.level {
		return nil
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	logDir := filepath.Join(l.dataDir, "logs")
	if err := os.MkdirAll(logDir, l.dirMode); err != nil {
		return err
	}

	now := l.now()
	logPath := filepath.Join(logDir, fmt.Sprintf("%s.log", now.Format("2006-01-02")))
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, l.fileMode)
	if err != nil {
		return err
	}
	defer file.Close()

	line := fmt.Sprintf("%s [%s] %s\n", now.Format(time.RFC3339), label, message)
	_, err = file.WriteString(line)
	return err
}

// parseLevel converts a log-level string into its internal numeric representation.
func parseLevel(level string) (int, error) {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "", "debug":
		return debugLevel, nil
	case "info":
		return infoLevel, nil
	case "warn":
		return warnLevel, nil
	case "error":
		return errorLevel, nil
	default:
		return 0, fmt.Errorf("unsupported log level: %s", level)
	}
}
