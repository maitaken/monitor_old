package exec

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gosuri/uilive"
	"github.com/ttacon/chalk"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/maitaken/monitor/option"
)

const POINT = "...................."

var opt *option.Option

type TermInfo struct {
	ClearStr string
	Lines    int
	Cols     int
}

type ShellWriter struct {
	writer *uilive.Writer
	cmd    string
	info   *TermInfo
	cancel context.CancelFunc
}

func New(cmd string) *ShellWriter {
	opt, _ = option.GetOption()
	w := uilive.New()
	w.Start()
	return &ShellWriter{
		writer: w,
		cmd:    cmd,
		info:   getTermInfo(),
	}
}

func (w *ShellWriter) Execute() context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		execCtx, execCancel := context.WithCancel(ctx)

		go w.executePrint(execCtx)

		startDate := time.Now()
		out, e := exec.CommandContext(execCtx, "sh", "-c", w.cmd).CombinedOutput()
		diff := time.Since(startDate)
		execCancel()

		if e != nil {
			w.errorPrint(out)
		} else {
			w.successPrint(out, diff)
		}
	}()
	return cancel
}

// 実行中の表示を行う関数
func (w *ShellWriter) executePrint(ctx context.Context) {
	var index int
	fmt.Fprint(w.writer, w.info.ClearStr)
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
func (w *ShellWriter) errorPrint(out []byte) {
	fmt.Fprint(w.writer, w.info.ClearStr)
	fmt.Fprint(w.writer, chalk.Red.Color("Command Faild"))
	fmt.Fprintln(w.writer, " : "+w.cmd)
	fmt.Fprint(w.writer, string(out))
}

// コマンドの実行結果が成功時の出力
func (w *ShellWriter) successPrint(out []byte, t time.Duration) {
	w.info = getTermInfo()

	outStrs := strings.Split(string(out), "\n")
	fmt.Fprint(w.writer, w.info.ClearStr)
	fmt.Fprint(w.writer, chalk.Green.Color("Command Success"), " : ", w.cmd)

	// outLines := len(outStrs)
	// 出力を短縮する場合
	for index, row := range outStrs {
		if opt.Shortened {
			if index == w.info.Lines-4 {
				break
			}
			outCols := len(row)
			if outCols <= w.info.Cols {
				fmt.Fprintln(w.writer, row)
			} else {
				trinRow := row[:w.info.Cols-3] + "..."
				fmt.Fprintln(w.writer, trinRow)
			}

		} else {
			fmt.Fprintln(w.writer, row)
		}
	}
	for i := len(outStrs); i < w.info.Lines-4; i++ {
		fmt.Fprintln(w.writer, "")
	}

	row := strings.Repeat("-", w.info.Cols)
	fmt.Fprintln(w.writer, "")
	fmt.Fprintln(w.writer, row)
	fmt.Fprintln(w.writer, " Time : ", t)
	fmt.Fprint(w.writer, row)
}

func getTermInfo() *TermInfo {
	clearStr, e := exec.Command("tput", "clear").Output()
	if e != nil {
		os.Exit(1)
	}
	// 標準出力先のサイズ(仮想端末のはず)
	cols, lines, e := terminal.GetSize(1)
	if e != nil {
		os.Exit(1)
	}

	return &TermInfo{
		ClearStr: string(clearStr),
		Lines:    lines,
		Cols:     cols,
	}
}
