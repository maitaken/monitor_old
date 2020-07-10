package exec

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"syscall"
	"time"

	"github.com/maitaken/monitor/option"
	"github.com/maitaken/monitor/print"
)

type ExecErrorType int

const (
	SUCCESS ExecErrorType = iota
	FAILED
	CANCELED
	TIMEOUT
)

const KILL_INTERVAL = 500

var opt *option.Option

type ExecError struct {
	Type    ExecErrorType
	message string
}

func NewExecError(t ExecErrorType, m string) *ExecError {
	return &ExecError{
		Type:    t,
		message: m,
	}
}

func (e *ExecError) Error() string {
	return e.message
}

func init() {
	opt = option.GetOption()
}

func Execute() context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		go print.Execute(ctx, opt.Cmd)

		startDate := time.Now()
		out, e := run(ctx, opt.Cmd)
		diff := time.Since(startDate)
		defer cancel()

		switch e.Type {
		case SUCCESS:
			print.Success(out, opt.Cmd, diff)
		case FAILED, TIMEOUT:
			print.Error(out, opt.Cmd, e)
		}
	}()
	return cancel
}

func run(ctx context.Context, command string) ([]byte, *ExecError) {
	var buf bytes.Buffer
	errChan := make(chan *ExecError, 1)
	done := make(chan struct{})

	cmd := exec.Command("sh", "-c", command)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	cmd.Stdout = &buf
	cmd.Stderr = &buf

	cmd.Start()
	pgid, _ := syscall.Getpgid(cmd.Process.Pid)

	go func() {
		timeout := generateTimer()
		select {
		case <-timeout:
			errChan <- NewExecError(TIMEOUT, fmt.Sprintf("Command Timeout (%ds)", opt.Timeout))
			syscall.Kill(-pgid, 15)

		case <-ctx.Done():
			errChan <- NewExecError(CANCELED, "Command Canceled")
			syscall.Kill(-pgid, 15)
		case <-done:
		}
	}()

	e := cmd.Wait()
	close(done)

	if e != nil {
		return buf.Bytes(), NewExecError(FAILED, "Command Failed")
	}
	if len(errChan) != 0 {
		return buf.Bytes(), <-errChan
	}
	return buf.Bytes(), NewExecError(SUCCESS, "")
}

func generateTimer() chan struct{} {
	c := make(chan struct{})
	go func() {
		if opt.Timeout != 0 {
			time.Sleep(time.Duration(opt.Timeout) * time.Second)
			c <- struct{}{}
		}
	}()
	return c
}
