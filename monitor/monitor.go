package monitor

import (
	"fmt"
	"os"
	"time"
)

const INTERVAL = 1

func Start(c chan<- string, file string) {
	var modifyTime string

	go func() {
		for {
			f, e := os.Stat(file)

			if e != nil {
				fmt.Println("Error : ", e)
				os.Exit(1)
			}

			now := f.ModTime().String()

			if now != modifyTime {
				modifyTime = now
				c <- "change"
			} else {
				c <- "nochange"
			}
			time.Sleep(INTERVAL * time.Second)
		}
	}()
}
