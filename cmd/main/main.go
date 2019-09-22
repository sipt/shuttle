package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/sipt/shuttle/cmd"
	"github.com/sipt/shuttle/pkg/close"
	"github.com/sirupsen/logrus"
)

func main() {
	flag.Parse()
	err := cmd.Start()
	if err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	err = close.Close(true)
	if err != nil {
		logrus.WithError(err).Error("close failed")
	}
}
