package log

import (
	"io"
	"log/slog"
	"os"
	"sync"
)

const (
	defaultLevel = slog.LevelWarn
)

var (
	defaultIO = os.Stdout
	logger    *slog.Logger
	mtx       sync.RWMutex
)

func init() {
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

	return slog.New(logger.Handler())
}

// Debug logs a message at debug level with the given key-value pairs as attributes.
func Debug(msg string, args ...any) {
	mtx.RLock()
	defer mtx.RUnlock()

	logger.Debug(msg, args...)
}

// Info logs a message at infoUser level with the given key-value pairs as attributes.
func Info(msg string, args ...any) {
	mtx.RLock()
	defer mtx.RUnlock()

	logger.Info(msg, args...)
}

// Warn logs a message at warn level with the given key-value pairs as attributes.
func Warn(msg string, args ...any) {
	mtx.RLock()
	defer mtx.RUnlock()

	logger.Warn(msg, args...)
}

// Error logs a message at error level with the given key-value pairs as attributes.
func Error(msg string, args ...any) {
	mtx.RLock()
	defer mtx.RUnlock()

	logger.Error(msg, args...)
}
