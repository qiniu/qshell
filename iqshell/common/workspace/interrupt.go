package workspace

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"os"
	"os/signal"
	"time"

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
		log.Alert("")
		log.DebugF("Got signal:%s", si)
		isCmdInterrupt = true
		Cancel()
		time.Sleep(time.Millisecond * 500)
		os.Exit(data.StatusUserCancel)
	}()
}
