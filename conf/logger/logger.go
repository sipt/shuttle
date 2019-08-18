package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func ConfigLogger() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.JSONFormatter{})
}
