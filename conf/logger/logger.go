package logger

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
)

func ConfigLogger(level logrus.Level) {
	logrus.SetLevel(level)
	logrus.SetFormatter(&logrus.TextFormatter{})
	gin.SetMode(gin.ReleaseMode)
}

func ConfigOutput(logpath string) error {
	if logpath == "" {
		return nil
	}
	_, err := os.Stat(logpath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(logpath, os.ModePerm)
		}
		if err != nil {
			return errors.Errorf("open file [%s] failed: %s", logpath, err)
		}
	}
	fileName := fmt.Sprintf("%s.log", time.Now().Format(time.RFC3339))
	fileName = path.Join(logpath, fileName)
	w, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return errors.Errorf("open file [%s] failed: %s", fileName, err)
	}
	logrus.SetOutput(w)
	Std = w
	return nil
}

var Std = os.Stdout
