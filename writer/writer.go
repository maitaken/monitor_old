package writer

import (
	"context"
	"os/exec"
	"time"

	"github.com/gosuri/uilive"
	"github.com/ttacon/chalk"
)

const POINT = "...................."

type Writer struct {
	Writer *uilive.Writer
}

func New() *Writer {
	return &Writer{
		Writer: uilive.New(),
	}
}

func (w *Writer) PrintSuccess(cmd string, out []byte) {
	w.Clear()
	w.Writer.Write([]byte(chalk.Green.Color("Command Success")))
	w.Writer.Write([]byte(" : " + cmd + "\n"))
	w.Writer.Write(out)
	return
}

func (w *Writer) PrintFaild(cmd string, out []byte, e error) {
	w.Clear()
	w.Writer.Write([]byte(chalk.Red.Color("Command Faild")))
	w.Writer.Write([]byte(" : " + cmd + "\n"))
	w.Writer.Write([]byte(e.Error() + "\n"))
	w.Writer.Write(out)
	return
}

func (w *Writer) PrintExecuting(ctx context.Context, cmd string) {
	w.Clear()

	for {
		for i := 0; i < len(POINT); i++ {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Millisecond * 500):
				w.Writer.Write([]byte(chalk.Yellow.Color("Command Executing")))
				w.Writer.Write([]byte(" : " + cmd + "\n"))
				w.Writer.Write([]byte(string(POINT[:i]) + "\n"))
			}
		}
	}
}

func (w *Writer) Clear() {
	out, _ := exec.Command("tput", "clear").Output()
	w.Writer.Write(out)
}
