package util

import (
	"context"
	"os/exec"

	"github.com/maitaken/monitor/writer"
)

func ExecShell(w *writer.Writer, cmd string) context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		printCtx, printCancel := context.WithCancel(ctx)
		go w.PrintExecuting(printCtx, cmd)

		out, e := exec.CommandContext(printCtx, "sh", "-c", cmd).CombinedOutput()
		printCancel()
		if e != nil {
			w.PrintFaild(cmd, out, e)
		} else {
			w.PrintSuccess(cmd, out)
		}
	}()

	return cancel
}
