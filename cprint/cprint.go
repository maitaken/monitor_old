package cprint

import (
	"context"
	"os/exec"
	"time"

	"github.com/gosuri/uilive"
	"github.com/ttacon/chalk"
)

const POINT = "...................."

// const

func PrintSuccess(w *uilive.Writer, cmds string, out []byte) {
	Clear(w)
	w.Write([]byte(chalk.Green.Color("Command Success")))
	w.Write([]byte(" : " + cmds + "\n"))
	w.Write(out)
	return
}

func PrintFaild(w *uilive.Writer, cmds string, out []byte, e error) {
	Clear(w)
	w.Write([]byte(chalk.Red.Color("Command Faild")))
	w.Write([]byte(" : " + cmds + "\n"))
	w.Write([]byte(e.Error() + "\n"))
	w.Write(out)
	return
}

func PrintExecuting(ctx context.Context, w *uilive.Writer, cmds string) {
	Clear(w)

	for {
		for i := 0; i < len(POINT); i++ {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Millisecond * 500):
				w.Write([]byte(chalk.Yellow.Color("Command Executing")))
				w.Write([]byte(" : " + cmds + "\n"))
				w.Write([]byte(string(POINT[:i]) + "\n"))
			}
		}
	}
}

func Clear(w *uilive.Writer) {
	out, _ := exec.Command("tput", "clear").Output()
	w.Write(out)
}
