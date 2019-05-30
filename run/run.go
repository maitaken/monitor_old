package run

import (
	"context"

	"github.com/maitaken/monitor/monitor"
	"github.com/maitaken/monitor/util"
	"github.com/maitaken/monitor/writer"
	"github.com/urfave/cli"
)

type CommandArgs struct {
	files []string
	cmd   string
}

func Run(c *cli.Context) {

	args := new(CommandArgs)

	if files := c.StringSlice("f"); len(files) != 0 {
		args.files = make([]string, len(files))
		copy(args.files, files)
		args.cmd = c.Args().Get(0)
	} else {
		args.files = make([]string, 1)
		args.files[0] = c.Args().Get(0)
		args.cmd = c.Args().Get(1)
	}

	fileChangeChan := make(chan string, len(args.files))
	writer := writer.New()
	writer.Writer.Start()
	var cancelFunc context.CancelFunc

	for _, file := range args.files {
		monitor.Start(fileChangeChan, file)
	}

	for {
		if <-fileChangeChan == "nochange" {
			continue
		}

		if cancelFunc != nil {
			cancelFunc()
		}

		cancelFunc = util.ExecShell(writer, args.cmd)
	}

}
