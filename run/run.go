package run

import (
	"context"

	"github.com/maitaken/monitor/exec"
	"github.com/maitaken/monitor/monitor"
	"github.com/maitaken/monitor/option"
	"github.com/urfave/cli"
)

func Run(c *cli.Context) {
	option.SetOption(c)
	opt := option.GetOption()

	fileChangeChan := make(chan struct{}, len(opt.TargetFile))
	var cancel context.CancelFunc

	for _, file := range opt.TargetFile {
		monitor.Start(fileChangeChan, file)
	}

	for {
		<-fileChangeChan

		if cancel != nil {
			cancel()
		}

		cancel = exec.Execute()
	}
}
