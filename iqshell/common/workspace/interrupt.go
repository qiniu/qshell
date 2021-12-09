package workspace

import (
	"os"
	"os/signal"

	"github.com/qiniu/qshell/v2/iqshell/common/log"
)

var (
	// 程序是否终端
	isCmdInterrupt = false
)

func IsCmdInterrupt() bool {
	return isCmdInterrupt
}

func observerCmdInterrupt() {
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt, os.Kill)
	go func() {
		si := <-s
		log.ErrorF("Got signal:", si)
		isCmdInterrupt = true
		Cancel()
	}()
}
