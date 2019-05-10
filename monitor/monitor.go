package monitor

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/gosuri/uilive"
	"github.com/maitaken/monitor/cprint"
)

const INTERVAL = 1

type Monitor struct {
	Cmds       string
	TargetFile string
	ModifyTime string
	Writer     *uilive.Writer
	Cancel     context.CancelFunc
}

type CommandResult struct {
	out    []byte
	e      error
	cancel context.CancelFunc
}

func New(file string, cmds string) *Monitor {
	return &Monitor{
		Cmds:       cmds,
		TargetFile: file,
	}
}

func (m *Monitor) Start() {
	m.Writer = uilive.New()
	m.Writer.Start()

	monitorChan := make(chan string, 1)
	go m._Monitor(monitorChan)

	for {
		isChanged := <-monitorChan
		if isChanged == "nochange" {
			continue
		}

		if m.Cancel != nil {
			m.Cancel()
			m.Cancel = nil
		}

		c := make(chan context.CancelFunc, 1)

		m.Cancel = m._ExecShell(c)

	}

}

func (m *Monitor) _Monitor(c chan string) {

	for {
		f, e := os.Stat(m.TargetFile)

		if e != nil {
			fmt.Println("Error : ", e)
			os.Exit(1)
		}

		now := f.ModTime().String()

		if now != m.ModifyTime {
			m.ModifyTime = now
			c <- "change"
		} else {
			c <- "nochange"
		}
		time.Sleep(INTERVAL * time.Second)
	}
}

func (m *Monitor) _ExecShell(c chan context.CancelFunc) context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context) {
		printCtx, printCancel := context.WithCancel(ctx)
		go cprint.PrintExecuting(printCtx, m.Writer, m.Cmds)

		out, e := exec.CommandContext(printCtx, "sh", "-c", m.Cmds).CombinedOutput()
		printCancel()
		if e != nil {
			cprint.PrintFaild(m.Writer, m.Cmds, out, e)
		} else {
			cprint.PrintSuccess(m.Writer, m.Cmds, out)
		}
		m.Cancel = nil

	}(ctx)

	return cancel
}
