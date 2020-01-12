package option

import (
	"os"
	"path/filepath"

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
	if paths := c.StringSlice("f"); len(paths) != 0 {

		for _, path := range paths {
			p, e := filepath.Glob(path)
			if e != nil {
				os.Exit(1)
			}
			opt.TargetFile = append(opt.TargetFile, p...)
		}

	} else {
		for i := 0; i < c.NArg()-1; i++ {
			path := c.Args().Get(i)

			p, e := filepath.Glob(path)
			if e != nil {
				os.Exit(1)
			}
			opt.TargetFile = append(opt.TargetFile, p...)
		}
	}

	opt.Cmd = c.Args().Get(c.NArg() - 1)
}
