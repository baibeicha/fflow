package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/lmittmann/tint"
	"gopkg.in/natefinch/lumberjack.v2"
)

type contextKey string

const RequestIDKey contextKey = "request_id"

type Config struct {
	Level      string
	Env        string
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

type Option func(*Config)

func WithLevel(level string) Option {
	return func(c *Config) { c.Level = level }
}

func WithEnv(env string) Option {
	return func(c *Config) { c.Env = env }
}

func WithFileOutput(filename string, maxSizeMB, maxBackups, maxAgeDays int, compress bool) Option {
	return func(c *Config) {
		c.Filename = filename
		c.MaxSize = maxSizeMB
		c.MaxBackups = maxBackups
		c.MaxAge = maxAgeDays
		c.Compress = compress
	}
}

type Logger struct {
	*slog.Logger
	file io.Closer
}

func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

type ContextHandler struct {
	slog.Handler
}

func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok && reqID != "" {
		r.AddAttrs(slog.String("request_id", reqID))
	}
	return h.Handler.Handle(ctx, r)
}

func (h *ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ContextHandler{Handler: h.Handler.WithAttrs(attrs)}
}

func (h *ContextHandler) WithGroup(name string) slog.Handler {
	return &ContextHandler{Handler: h.Handler.WithGroup(name)}
}

type MultiHandler struct {
	handlers []slog.Handler
}

func NewMultiHandler(handlers ...slog.Handler) *MultiHandler {
	return &MultiHandler{handlers: handlers}
}

func (m *MultiHandler) Enabled(ctx context.Context, l slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, l) {
			return true
		}
	}
	return false
}

func (m *MultiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range m.handlers {
		if h.Enabled(ctx, r.Level) {
			if err := h.Handle(ctx, r); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithAttrs(attrs)
	}
	return NewMultiHandler(handlers...)
}

func (m *MultiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithGroup(name)
	}
	return NewMultiHandler(handlers...)
}

func New(opts ...Option) (*Logger, error) {
	cfg := &Config{
		Level: "info",
		Env:   "development",
	}

	for _, opt := range opts {
		opt(cfg)
	}

	var level slog.LevelVar
	switch strings.ToLower(cfg.Level) {
	case "debug":
		level.Set(slog.LevelDebug)
	case "warn", "warning":
		level.Set(slog.LevelWarn)
	case "error", "err":
		level.Set(slog.LevelError)
	default:
		level.Set(slog.LevelInfo)
	}

	var handlers []slog.Handler

	if cfg.Env == "production" || cfg.Env == "prod" {
		handlers = append(handlers, slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     &level,
			AddSource: true,
		}))
	} else {
		handlers = append(handlers, tint.NewHandler(os.Stdout, &tint.Options{
			Level:      &level,
			AddSource:  true,
			TimeFormat: "15:04:05.000",
		}))
	}

	var fileCloser io.Closer
	if cfg.Filename != "" {
		fileWriter := &lumberjack.Logger{
			Filename:   cfg.Filename,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		}
		fileCloser = fileWriter

		handlers = append(handlers, slog.NewJSONHandler(fileWriter, &slog.HandlerOptions{
			Level:     &level,
			AddSource: true,
		}))
	}

	multiHandler := NewMultiHandler(handlers...)
	contextHandler := &ContextHandler{Handler: multiHandler}

	logger := slog.New(contextHandler)
	slog.SetDefault(logger)

	return &Logger{
		Logger: logger,
		file:   fileCloser,
	}, nil
}

func WithRequestID(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, reqID)
}

func GetRequestID(ctx context.Context) string {
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		return reqID
	}
	return ""
}
