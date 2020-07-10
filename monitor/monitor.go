package monitor

import (
	"log"
	"os"
	"time"
)

const INTERVAL = 500

func Start(c chan struct{}, file string) {
	modifyTime := getModTime(file)
	// TODO
	if len(c) == 0 {
		c <- struct{}{}
	}

	go func() {
		for range time.Tick(INTERVAL * time.Millisecond) {
			now := getModTime(file)

			if now != modifyTime {
				modifyTime = now
				c <- struct{}{}
			}
		}
	}()
}

func getModTime(file string) string {
	f, e := os.Stat(file)
	if e != nil {
		log.Fatal(e, file)
	}
	return f.ModTime().String()
}
