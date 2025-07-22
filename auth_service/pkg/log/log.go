package log

import (
	"io"
	"log/slog"
	"os"
	"sync"
)

const (
	defaultLevel       = slog.LevelWarn
	defaultServiceName = "auth_service"
)

var (
	defaultIO   io.Writer = os.Stdout
	logger      *slog.Logger
	mtx         sync.RWMutex
	serviceName string

	logFilePath = "logger/auth_service.log"
)

func init() {
	f, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		// если файл открыть не удалось, логгер будет писать только в stdout
		defaultIO = os.Stdout
	} else {
		defaultIO = io.MultiWriter(os.Stdout, f)
	}

	logger = slog.New(
		slog.NewJSONHandler(
			defaultIO,
			&slog.HandlerOptions{Level: defaultLevel},
		),
	)
}

// Configure sets up the global logger with the provided writer, options, and format.
// If w is nil, it defaults to os.Stdout. If opts is nil, it defaults to LevelWarn.
// The json parameter determines whether to use JSON or text format for log output.
func Configure(w io.Writer, opts *slog.HandlerOptions, json bool) {
	mtx.Lock()
	defer mtx.Unlock()

	if w == nil {
		w = defaultIO
	}

	if opts == nil {
		opts = &slog.HandlerOptions{Level: defaultLevel}
	}

	serviceName = defaultServiceName

	if json {
		logger = slog.New(slog.NewJSONHandler(w, opts))
	} else {
		logger = slog.New(slog.NewTextHandler(w, opts))
	}
}

// Copy returns a new logger instance that uses the same handler as the global logger.
// This allows for independent logger instances that maintain the same configuration.
func Copy() *slog.Logger {
	mtx.RLock()
	defer mtx.RUnlock()

	return slog.New(logger.Handler()).With("service", serviceName)
}

// Debug logs a message at debug level with the given key-value pairs as attributes.
func Debug(msg string, args ...any) {
	mtx.RLock()
	defer mtx.RUnlock()

	logger.With("service", serviceName).Debug(msg, args...)
}

// Info logs a message at infoUser level with the given key-value pairs as attributes.
func Info(msg string, args ...any) {
	mtx.RLock()
	defer mtx.RUnlock()

	logger.With("service", serviceName).Info(msg, args...)
}

// Warn logs a message at warn level with the given key-value pairs as attributes.
func Warn(msg string, args ...any) {
	mtx.RLock()
	defer mtx.RUnlock()

	logger.With("service", serviceName).Warn(msg, args...)
}

// Error logs a message at error level with the given key-value pairs as attributes.
func Error(msg string, args ...any) {
	mtx.RLock()
	defer mtx.RUnlock()

	logger.With("service", serviceName).Error(msg, args...)
}
