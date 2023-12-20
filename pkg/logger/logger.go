package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

func New() *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})

	return logger
}
