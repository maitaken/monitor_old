package main

import (
	"fmt"
	"os"

	"github.com/maitaken/monitor/monitor"
	"github.com/urfave/cli"
)

var m *monitor.Monitor

func main() {
	app := cli.NewApp()

	// fmt.Println(m)

	app.Name = "Monitor"
	app.Usage = "Monitor single file changes and run a shell"
	app.Version = "0.0.1"

	app.Action = Handler
	app.Run(os.Args)
}

func Handler(c *cli.Context) (e error) {

	if c.NArg() < 2 {
		fmt.Println("invalid argument")
		os.Exit(1)
	}

	m := &monitor.Monitor{
		TargetFile: c.Args().Get(0),
		Cmds:       c.Args().Get(1),
	}

	m.Start()

	return e
}
