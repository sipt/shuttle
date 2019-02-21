package main

import (
	"flag"

	"github.com/sipt/shuttle/cmd"
)

func main() {
	configPath := flag.String("c", "shuttle.yaml", "configuration file path")
	logMode := flag.String("l", "file", "logMode: off | console | file")
	logPath := flag.String("lp", "logs", "logs path")
	flag.Parse()
	cmd.Run(*logMode, *logPath, *configPath)
}
