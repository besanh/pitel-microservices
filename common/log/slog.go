package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gookit/slog/rotatefile"
)

var CustomLogger *slog.Logger

func Info(args ...any) {
	CustomLogger.
		Info(fmt.Sprint(args...))
}

func Warn(args ...any) {
	CustomLogger.
		Warn(fmt.Sprint(args...))
}

func Error(args ...any) {
	CustomLogger.
		Error(fmt.Sprint(args...))
}

func Debug(args ...any) {
	CustomLogger.
		Debug(fmt.Sprint(args...))
}

func Infof(msg string, args ...any) {
	CustomLogger.
		Info(fmt.Sprintf(msg, args...))
}

func Warnf(msg string, args ...any) {
	CustomLogger.
		Warn(fmt.Sprintf(msg, args...))
}

func Errorf(msg string, args ...any) {
	CustomLogger.
		Error(fmt.Sprintf(msg, args...))
}

func Debugf(msg string, args ...any) {
	CustomLogger.
		Debug(fmt.Sprintf(msg, args...))
}

func InfoContext(ctx context.Context, args ...any) {
	CustomLogger.
		InfoContext(ctx, fmt.Sprint(args...))
}

func WarnContext(ctx context.Context, args ...any) {
	CustomLogger.
		WarnContext(ctx, fmt.Sprint(args...))
}

func ErrorContext(ctx context.Context, args ...any) {
	CustomLogger.
		ErrorContext(ctx, fmt.Sprint(args...))
}

func DebugContext(ctx context.Context, args ...any) {
	CustomLogger.
		DebugContext(ctx, fmt.Sprint(args...))
}

func InfofContext(ctx context.Context, msg string, args ...any) {
	CustomLogger.
		InfoContext(ctx, fmt.Sprintf(msg, args...))
}

func WarnfContext(ctx context.Context, msg string, args ...any) {
	CustomLogger.
		WarnContext(ctx, fmt.Sprintf(msg, args...))
}

func ErrorfContext(ctx context.Context, msg string, args ...any) {
	CustomLogger.
		ErrorContext(ctx, fmt.Sprintf(msg, args...))
}

func DebugfContext(ctx context.Context, msg string, args ...any) {
	CustomLogger.
		DebugContext(ctx, fmt.Sprintf(msg, args...))
}

func InitLogger(level string, logFile string) {
	logLevel := slog.LevelDebug
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "error":
		logLevel = slog.LevelError
	case "warn":
		logLevel = slog.LevelWarn
	}
	rotateWriter, err := rotatefile.NewConfig(logFile).Create()
	if err != nil {
		panic(err)
	}
	w := io.MultiWriter(os.Stdout, rotateWriter)
	var handler slog.Handler
	handler = slog.NewTextHandler(w, &slog.HandlerOptions{
		Level: logLevel,
	})
	handler = LoggerHandler{handler}
	CustomLogger = slog.New(handler)
}

type LoggerHandler struct {
	slog.Handler
}

func (h LoggerHandler) Handle(ctx context.Context, r slog.Record) error {
	if logId, ok := ctx.Value("trace_id").(string); ok {
		r.Add("trace_id", logId)
	} else {
		r.Add("trace_id", "unknown")
	}
	_, path, numLine, _ := runtime.Caller(4)
	srcFile := filepath.Base(path)
	r.Add(slog.String("file", fmt.Sprintf("%s:%d", srcFile, numLine)))
	return h.Handler.Handle(ctx, r)
}
