package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

const (
	_appName  = "simple-http-server"
	_appUsage = "simple-http-server implement simple response from request"
	_port     = "port"
)

var (
	Version = "0.0.0"
	BuildAt = ""
)

var (
	fListenPort = &cli.IntFlag{
		Name:    _port,
		Aliases: []string{"p"},
		EnvVars: []string{"PORT"},
		Usage:   "HTTP port",
		Value:   8080,
	}
)

func main() {
	app := &cli.App{
		Name:    _appName,
		Usage:   fmt.Sprintf("%s\nbuild at: %s", _appUsage, BuildAt),
		Version: Version,
		Flags: []cli.Flag{
			fListenPort,
		},
		Action: action,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal("simple-http-server failed:", err)
	}
}
