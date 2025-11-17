package initializers

import (
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func InitLogger() {
	Log = logrus.New()

	if level, err := logrus.ParseLevel("info"); err == nil {
		Log.SetLevel(level)
	}
	Log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	logPath := filepath.Join(GetEnv("LOG_DIR", "logs"), GetEnv("LOG_FILE", "app.log"))
	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err == nil {
		if file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err == nil {
			Log.SetOutput(io.MultiWriter(os.Stdout, file))
			Log.WithField("log_file", logPath).Info("Logger initialized")
			return
		}
	}

	Log.SetOutput(os.Stdout)
	Log.Info("Logger initialized")
}
