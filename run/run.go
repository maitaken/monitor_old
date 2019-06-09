package run

import (
	"context"

	"github.com/maitaken/monitor/monitor"
	"github.com/maitaken/monitor/util"
	"github.com/urfave/cli"
)

type CommandArgs struct {
	Files []string
	Cmd   string
}

func Run(c *cli.Context) {

	args := new(CommandArgs)

	// 引数とオプションの整理
	if files := c.StringSlice("f"); len(files) != 0 {
		args.Files = make([]string, len(files))
		copy(args.Files, files)
		args.Cmd = c.Args().Get(0)
	} else {
		args.Files = make([]string, 1)
		args.Files[0] = c.Args().Get(0)
		args.Cmd = c.Args().Get(1)
	}

	fileChangeChan := make(chan string, len(args.Files))
	shell := util.New(args.Cmd)
	var cancelFunc context.CancelFunc

	for _, file := range args.Files {
		monitor.Start(fileChangeChan, file)
	}

	for {
		if <-fileChangeChan == "nochange" {
			continue
		}

		if cancelFunc != nil {
			cancelFunc()
		}

		cancelFunc = shell.Execute()
	}

}
