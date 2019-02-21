package main

import "github.com/sipt/shuttle/cmd"

//export Run
func Run(logMode, logPath, configPath string) {
	go cmd.Run(logMode, logPath, configPath)
}

func main() {
}
