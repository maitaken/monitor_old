package option

import (
	"github.com/urfave/cli"
)

// Singleton

var opt *Option

type Option struct {
	TargetFile []string
	Cmd        string
}

func init() {
	opt = &Option{}
}

func GetOption() *Option {
	return opt
}

func SetOption(c *cli.Context) {
	// 引数とオプションの整理
	if files := c.StringSlice("f"); len(files) != 0 {
		opt.TargetFile = make([]string, len(files))
		copy(opt.TargetFile, files)

	} else {
		opt.TargetFile = make([]string, c.NArg()-1)
		for i := 0; i < c.NArg()-1; i++ {
			opt.TargetFile[i] = c.Args().Get(i)
		}
		opt.TargetFile[0] = c.Args().Get(0)
	}

	opt.Cmd = c.Args().Get(c.NArg() - 1)

}
