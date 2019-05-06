package cprint

import (
	"os/exec"

	"github.com/gosuri/uilive"
	"github.com/ttacon/chalk"
)

func PrintSuccess(w *uilive.Writer, cmds string, out []byte) {
	Clear(w)
	w.Write([]byte(chalk.Green.Color("Command Success")))
	w.Write([]byte(" : " + cmds + "\n"))
	w.Write(out)
	return
}

func PrintFaild(w *uilive.Writer, cmds string, e error) {
	Clear(w)
	w.Write([]byte(chalk.Red.Color("Command Faild")))
	w.Write([]byte(" : " + cmds + "\n"))
	w.Write([]byte(e.Error()))
	return
}

func Clear(w *uilive.Writer) {
	out, _ := exec.Command("tput", "clear").Output()
	w.Write(out)
}
