package monitor

import (
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

// func (m *Monitor) New() *Monitor {
// 	return &Monitor{
// 		Cmds:
// 	}
// }

func (m *Monitor) Start() {
	m.Writer = uilive.New()
	m.Writer.Start()

	for {
		f, e := os.Stat(m.TargetFile)

		if e != nil {
			fmt.Println("FileNotFoundError : ", e)
			os.Exit(1)
		}

		now := f.ModTime().String()

		if now == m.ModifyTime {
			time.Sleep(INTERVAL * time.Second)
			continue
		}

		m.ModifyTime = now

		c := cprint.PrintExecuting(m.Writer, m.Cmds)

		out, e := m.ExecShell(m.Cmds)

		c <- true
		if e != nil {
			cprint.PrintFaild(m.Writer, m.Cmds, out, e)
		} else {
			cprint.PrintSuccess(m.Writer, m.Cmds, out)
		}

	}
}

func (m *Monitor) ExecShell(cmds string) ([]byte, error) {
	out, e := exec.Command("sh", "-c", m.Cmds).CombinedOutput()
	return out, e
}
