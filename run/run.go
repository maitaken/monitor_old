package run

import (
	"context"
	"fmt"
	"os"

	"github.com/maitaken/monitor/exec"
	"github.com/maitaken/monitor/monitor"
	"github.com/maitaken/monitor/option"
	"github.com/urfave/cli"
)

func Run(c *cli.Context) {

	option.SetOption(c)
	opt, e := option.GetOption()
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}

	fileChangeChan := make(chan string, len(opt.TargetFile))
	shell := exec.New(opt.Cmd)
	var cancelFunc context.CancelFunc

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
