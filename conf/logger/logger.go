package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func ConfigLogger(level logrus.Level) {
	logrus.SetLevel(level)
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.TextFormatter{})
}
