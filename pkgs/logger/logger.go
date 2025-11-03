package logger

import (
	"io"
	"log"
	"os"
	"sync"
	"time"
)

type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Context   map[string]interface{} `json:"context,omitempty"`
}

type Logger struct {
	mu          sync.RWMutex
	entries     []LogEntry
	maxSize     int
	broadcastFn func(LogEntry)
}

var globalLogger *Logger

func init() {
	globalLogger = NewLogger(1000)
}

func NewLogger(maxSize int) *Logger {
	return &Logger{
		entries: make([]LogEntry, 0, maxSize),
		maxSize: maxSize,
	}
}

func GetLogger() *Logger {
	return globalLogger
}

func (l *Logger) Log(level LogLevel, message string, context map[string]interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	entry := LogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     string(level),
		Message:   message,
		Context:   context,
	}

	l.entries = append(l.entries, entry)
	if len(l.entries) > l.maxSize {
		l.entries = l.entries[1:]
	}

	log.Printf("[%s] %s", string(level), message)

	if l.broadcastFn != nil {
		go l.broadcastFn(entry)
	}
}

func (l *Logger) Info(message string) {
	l.Log(LogLevelInfo, message, nil)
}

func (l *Logger) Warn(message string) {
	l.Log(LogLevelWarn, message, nil)
}

func (l *Logger) Error(message string) {
	l.Log(LogLevelError, message, nil)
}

func (l *Logger) Debug(message string) {
	l.Log(LogLevelDebug, message, nil)
}

func (l *Logger) InfoWithContext(message string, context map[string]interface{}) {
	l.Log(LogLevelInfo, message, context)
}

func (l *Logger) GetLogs(level string, limit int) []LogEntry {
	l.mu.RLock()
	defer l.mu.RUnlock()

	filtered := make([]LogEntry, 0)
	for i := len(l.entries) - 1; i >= 0; i-- {
		entry := l.entries[i]
		if level == "" || entry.Level == level {
			filtered = append(filtered, entry)
			if limit > 0 && len(filtered) >= limit {
				break
			}
		}
	}

	return filtered
}

func (l *Logger) SetBroadcastFunc(fn func(LogEntry)) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.broadcastFn = fn
}

func Info(message string) {
	globalLogger.Info(message)
}

func Warn(message string) {
	globalLogger.Warn(message)
}

func Error(message string) {
	globalLogger.Error(message)
}

func Debug(message string) {
	globalLogger.Debug(message)
}

func InfoWithContext(message string, context map[string]interface{}) {
	globalLogger.InfoWithContext(message, context)
}

type logWriter struct {
	level  LogLevel
	output io.Writer
}

func (w *logWriter) Write(p []byte) (n int, err error) {
	message := string(p)
	if len(message) > 0 && message[len(message)-1] == '\n' {
		message = message[:len(message)-1]
	}
	globalLogger.Log(w.level, message, nil)
	return w.output.Write(p)
}

func SetupStdLogger() {
	writer := &logWriter{
		level:  LogLevelInfo,
		output: os.Stdout,
	}
	log.SetOutput(writer)
	log.SetFlags(log.LstdFlags)
}
