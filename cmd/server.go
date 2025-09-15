package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sipt/shuttle/cmd/api"
	"github.com/sipt/shuttle/conf"
	"github.com/sipt/shuttle/conf/logger"
	"github.com/sipt/shuttle/controller"
	"github.com/sipt/shuttle/events"
	"github.com/sipt/shuttle/handle"
	"github.com/sipt/shuttle/inbound"
	"github.com/sirupsen/logrus"

	closepkg "github.com/sipt/shuttle/pkg/close"
	"github.com/sipt/shuttle/pkg/enhance"

	_ "github.com/sipt/shuttle/conn/stream/include"
	_ "github.com/sipt/shuttle/events/include"
)

var Path = flag.String("c", os.Getenv("CONFIG_PATH"), "config file Path")
var RuntimePath = flag.String("r", os.Getenv("RUNTIME_PATH"), "runtime file Path")
var Encoding = flag.String("e", os.Getenv("ENCODING"), "config file Encoding")
var LogPath = flag.String("logpath", os.Getenv("LOGGER_PATH"), "logger file")

func init() {
	// register func to api
	api.StartFunc = Start
	api.CloseFunc = Close
	api.CheckConfig = CheckConfig
}

func Start() (err error) {
	configPath := *Path
	if configPath == "" {
		return fmt.Errorf("config file path is empty")
	}
	runtimePath := *RuntimePath
	if runtimePath == "" {
		return fmt.Errorf("runtime file path is empty")
	}
	api.Status = api.StatusStarting
	defer func() {
		if err != nil {
			api.Status = api.StatusStopped
		} else {
			api.Status = api.StatusRunning
		}
	}()
	logrus.SetLevel(logrus.DebugLevel)
	err = logger.ConfigOutput(*LogPath)
	if err != nil {
		panic(err)
	}

	configEncoding, err := getEncoding(configPath)
	if err != nil {
		logrus.WithError(err).Error("get config encoding failed")
		return err
	}
	runtimeEncoding, err := getEncoding(runtimePath)
	if err != nil {
		logrus.WithError(err).Error("get runtime encoding failed")
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	params := map[string]string{"path": configPath}
	config, err := conf.LoadConfig(ctx, "file", configEncoding, params, func() {
		fmt.Println("config file change")
	})
	if err != nil {
		logrus.WithError(err).Error("load config failed")
		return err
	}
	params = map[string]string{"path": *RuntimePath}
	runtime, err := conf.LoadRuntime(ctx, "file", runtimeEncoding, params)
	if err != nil {
		logrus.WithError(err).Error("load runtime failed")
		return err
	}
	l, err := logrus.ParseLevel(config.General.LoggerLevel)
	if err != nil {
		l = logrus.DebugLevel
	}
	logger.ConfigLogger(l)
	if err != nil {
		logrus.WithError(err).Error("load config failed")
		return err
	}
	err = conf.ApplyConfig(ctx, config, runtime)
	if err != nil {
		logrus.WithError(err).Error("apply config failed")
		return err
	}
	closer, err := controller.ApplyConfig(config)
	if err != nil {
		logrus.WithError(err).Error("start controller failed")
		return err
	}
	err = inbound.ApplyConfig(ctx, config, handle.Handle())
	if err != nil {
		logrus.WithError(err).Error("start inbound failed")
		return err
	}
	err = events.AutoDial(ctx)
	if err != nil {
		logrus.WithError(err).Error("start events failed")
		return err
	}
	// start enhance mode
	enhanceMode := enhance.NewEnhanceMode()
	err = enhanceMode.Start()
	if err != nil {
		logrus.WithError(err).Error("start enhance mode failed")
		return err
	}

	logrus.Info("server starting...")
	closepkg.AppendCloser(func() error {
		cancel()
		return nil
	})
	closepkg.AppendCloser(func() error {
		enhanceMode.Stop()
		closer()
		return nil
	})
	return nil
}

func CheckConfig() error {
	configPath := *Path
	if configPath == "" {
		return fmt.Errorf("config file path is empty")
	}
	configEncoding, err := getEncoding(configPath)
	if err != nil {
		logrus.WithError(err).Error("get config encoding failed")
		return err
	}
	params := map[string]string{"path": configPath}
	_, err = conf.LoadConfig(context.Background(), "file", configEncoding, params, func() {
		fmt.Println("config file change")
	})
	if err != nil {
		logrus.WithError(err).Error("load config failed")
		return err
	}
	return nil
}

func Close() error {
	_ = closepkg.Close(true)
	api.Status = api.StatusStopped
	return nil
}

func getEncoding(filePath string) (string, error) {
	encoding := *Encoding
	if encoding == "" {
		ext := filepath.Ext(filePath)
		switch ext {
		case ".json":
			encoding = "json"
		case ".yaml", ".yml":
			encoding = "yaml"
		case ".toml":
			encoding = "toml"
		default:
			return "", fmt.Errorf("unsupported config file extension: %s", ext)
		}
	}
	return encoding, nil
}
