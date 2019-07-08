package run

import (
	"context"
	"fmt"

	"github.com/maitaken/monitor/monitor"
	"github.com/maitaken/monitor/option"
	"github.com/maitaken/monitor/util"
	"github.com/urfave/cli"
)

func Run(c *cli.Context) {

	option.SetOption(c)
	opt := option.GetOption()

	fileChangeChan := make(chan string, len(opt.TargetFile))
	shell := util.New(opt.Cmd)
	var cancelFunc context.CancelFunc

	fmt.Println(opt.TargetFile)
	for _, file := range opt.TargetFile {
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
