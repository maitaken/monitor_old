package main

import (
	"os"

	"github.com/maitaken/monitor/run"
	"github.com/urfave/cli"
)

var (
	Version string
)

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
	app.Usage = "Monitor files and run a command when modified"
	app.Author = "maitaken"
	app.Version = Version
	app.Flags = appFlag()

	return app
}

func appFlag() []cli.Flag {
	return []cli.Flag{
		cli.StringSliceFlag{
			Name:  "file, f",
			Usage: "specify the target file",
		},
		cli.BoolFlag{
			Name:  "shortened-output, s",
			Usage: "print shortened output",
		},
	}
}
