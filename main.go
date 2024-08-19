package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/vanderkilu/github-service/cmd"
)

func main() {
	app := &cli.App{
		Name:    "github-service",
		Usage:   "root command of the github-service service",
		Version: os.Getenv("COMMIT"),
		Commands: []*cli.Command{
			cmd.CronMonitorRepo(),
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
	