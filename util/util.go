package util

import (
	"context"
	"fmt"
	"os/exec"
	"sync"
	"time"

	"github.com/ttacon/chalk"

	"github.com/gosuri/uilive"
)

const POINT = "...................."

type ShellWriter struct {
	writer *uilive.Writer
	cmd    string
	mu     sync.Mutex
	cancel context.CancelFunc
}

func New(cmd string) *ShellWriter {
	w := uilive.New()
	w.Start()
	return &ShellWriter{
		cmd:    cmd,
		writer: w,
	}
}

func (w *ShellWriter) Execute() context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		execCtx, execCancel := context.WithCancel(ctx)

		go w.executePrint(execCtx)

		out, e := exec.CommandContext(execCtx, "sh", "-c", w.cmd).CombinedOutput()
		execCancel()

		if e != nil {
			w.errorPrint(e)
		} else {
			w.successPrint(out)
		}
	}()
	return cancel
}

// 実行中の表示を行う関数
func (w *ShellWriter) executePrint(ctx context.Context) {
	var index int
	fmt.Fprint(w.writer, string(clear()))
	for range time.Tick(100 * time.Millisecond) {
		if len(POINT) == index {
			index = 0
		}
		select {
		case <-ctx.Done():
			return
		default:
			fmt.Fprint(w.writer, chalk.Yellow.Color("Command Executing"))
			fmt.Fprintln(w.writer, " : "+w.cmd)
			fmt.Fprintln(w.writer, POINT[:index])
		}
		index += 1
	}
}

// コマンドの実行結果がエラー時の出力
func (w *ShellWriter) errorPrint(e error) {
	fmt.Fprint(w.writer, string(clear()))
	fmt.Fprint(w.writer, chalk.Red.Color("Command Faild"))
	fmt.Fprintln(w.writer, " : "+w.cmd)
	fmt.Fprint(w.writer, e.Error())
}

// コマンドの実行結果が成功時の出力
func (w *ShellWriter) successPrint(out []byte) {
	fmt.Fprint(w.writer, string(clear()))
	fmt.Fprint(w.writer, chalk.Green.Color("Command Success"))
	fmt.Fprintln(w.writer, " : "+w.cmd)
	fmt.Fprint(w.writer, string(out))
}

func clear() []byte {
	out, _ := exec.Command("tput", "clear").Output()
	return out
}
