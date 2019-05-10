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
}

type CommandResult struct {
	out []byte
	e   error
}

// func (m *Monitor) New() *Monitor {
// 	return &Monitor{
// 		Cmds:
// 	}
// }

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

		c := cprint.PrintExecuting(m.Writer, m.Cmds)

		out, e, _ := m.ExecShell(m.Cmds)

		c <- true

		if e != nil {
			cprint.PrintFaild(m.Writer, m.Cmds, out, e)
		} else {
			cprint.PrintSuccess(m.Writer, m.Cmds, out)
		}

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

func (m *Monitor) ExecShell(cmds string) ([]byte, error, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	out, e := exec.CommandContext(ctx, "sh", "-c", m.Cmds).CombinedOutput()
	return out, e, cancel
}
