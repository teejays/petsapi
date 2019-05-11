package main

import (
	"github.com/teejays/clog"

	"./server"
)

var listenPort = 8080

func main() {
	// Increase the log level
	clog.LogLevel = 0

	err := server.StartServer("", listenPort)
	if err != nil {
		clog.FatalErr(err)
	}

}
