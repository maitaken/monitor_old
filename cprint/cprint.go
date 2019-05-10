package cprint

import (
	"os/exec"
	"time"

	"github.com/gosuri/uilive"
	"github.com/ttacon/chalk"
)

const POINT = "...................."

func PrintSuccess(w *uilive.Writer, cmds string, out []byte) {
	Clear(w)
	w.Write([]byte(chalk.Green.Color("Command Success")))
	w.Write([]byte(" : " + cmds + "\n"))
	w.Write(out)
	return
}

func PrintFaild(w *uilive.Writer, cmds string, out []byte, e error) {
	w.Write([]byte(chalk.Red.Color("Command Faild")))
	w.Write([]byte(" : " + cmds + "\n"))
	w.Write([]byte(e.Error() + "\n"))
	w.Write(out)
	return
}

func PrintExecuting(w *uilive.Writer, cmds string) chan bool {
	c := make(chan bool, 1)

	Clear(w)

	go func() {
		for {
			for i := 0; i < len(POINT); i++ {
				if len(c) != 0 {
					return
				}

				w.Write([]byte(chalk.Yellow.Color("Command Executing")))
				w.Write([]byte(" : " + cmds + "\n"))
				w.Write([]byte(string(POINT[:i]) + "\n"))
				w.Flush()
				time.Sleep(time.Millisecond * 500)
			}
		}
	}()

	return c
}

func Clear(w *uilive.Writer) {
	out, _ := exec.Command("tput", "clear").Output()
	w.Write(out)
}
