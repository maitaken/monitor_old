package main

import (
	"os"

	"github.com/maitaken/monitor/run"
	"github.com/urfave/cli"
)

const version = "v0.0.4"

func main() {
	app := newApp()

	app.Action = run.Run

	if e := app.Run(os.Args); e != nil {
		exitCode := 1
		os.Exit(exitCode)
	}
}

func newApp() *cli.App {
	app := cli.NewApp()
	app.Name = "Monitor"
	app.Usage = "Monitor files changed and run a command"
	app.Author = "maitaken"
	app.Version = version
	app.Flags = appFlag()

	return app
}

func appFlag() []cli.Flag {
	return []cli.Flag{
		cli.StringSliceFlag{
			Name:  "file, f",
			Usage: "Monitoring file name",
		},
	}
}
