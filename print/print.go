package print

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gosuri/uilive"
	"github.com/maitaken/monitor/option"
	"github.com/ttacon/chalk"
	"golang.org/x/crypto/ssh/terminal"
)

type TermInfo struct {
	ClearStr string
	Lines    int
	Cols     int
}

type ShellWriter struct {
	w    *uilive.Writer
	info *TermInfo
}

const POINT = "...................."
const TIME_LINES_NUMBER = 4

var opt *option.Option
var writer *ShellWriter

func init() {
	opt = option.GetOption()
	w := uilive.New()
	w.Start()
	writer = &ShellWriter{
		w:    w,
		info: getTermInfo(),
	}
}

func Execute(ctx context.Context, cmd string) {
	var index int
	fmt.Fprint(writer.w, writer.info.ClearStr)
	for range time.Tick(100 * time.Millisecond) {
		if len(POINT) == index {
			index = 0
		}
		select {
		case <-ctx.Done():
			return
		default:
			fmt.Fprint(writer.w, chalk.Yellow.Color("Command Executing"))
			fmt.Fprintln(writer.w, " : "+cmd)
			fmt.Fprintln(writer.w, POINT[:index])
		}
		index++
	}
}

// コマンドの実行結果が成功時の出力
func Success(out []byte, cmd string, t time.Duration) {
	writer.info = getTermInfo()

	outStrs := strings.Split(string(out), "\n")
	fmt.Fprint(writer.w, writer.info.ClearStr)
	fmt.Fprintln(writer.w, chalk.Green.Color("Command Success"), " : ", cmd)

	// outLines := len(outStrs)
	// 出力を短縮する場合
	for index, row := range outStrs {
		if opt.Shortened {
			if index == writer.info.Lines-TIME_LINES_NUMBER {
				break
			}
			outCols := len(row)
			if outCols <= writer.info.Cols {
				fmt.Fprintln(writer.w, row)
			} else {
				trinRow := row[:writer.info.Cols-3] + "..."
				fmt.Fprintln(writer.w, trinRow)
			}

		} else {
			fmt.Fprintln(writer.w, row)
		}
	}
	if opt.Shortened {
		for i := len(outStrs); i < writer.info.Lines-TIME_LINES_NUMBER; i++ {
			fmt.Fprintln(writer.w, "")
		}
	}

	row := strings.Repeat("-", writer.info.Cols)
	fmt.Fprintln(writer.w, row)
	fmt.Fprintln(writer.w, " Time : ", t)
	fmt.Fprint(writer.w, row)
}

// コマンドの実行結果がエラー時の出力
func Error(out []byte, cmd string, e error) {
	fmt.Fprint(writer.w, writer.info.ClearStr)
	fmt.Fprint(writer.w, chalk.Red.Color(e.Error()))
	fmt.Fprintln(writer.w, " : "+cmd)
	fmt.Fprint(writer.w, string(out))
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
