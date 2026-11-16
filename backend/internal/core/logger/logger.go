package logger

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

func Init(logDir string) {
	logWriter := &lumberjack.Logger{
		Filename: filepath.Join(logDir, "app.log"),
		MaxSize:  50,
		Compress: true,
	}

	multiWriter := io.MultiWriter(os.Stdout, logWriter)

	handler := slog.NewTextHandler(multiWriter, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	slog.SetDefault(slog.New(handler))
}
