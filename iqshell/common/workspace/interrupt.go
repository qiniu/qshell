package workspace

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"os"
	"os/signal"
	"sync/atomic"
	"time"

	"github.com/qiniu/qshell/v2/iqshell/common/log"
)

var (
	// 程序是否退出
	isCmdInterrupt uint32 = 0
)

func IsCmdInterrupt() bool {
	return atomic.LoadUint32(&isCmdInterrupt) > 0
}

func observerCmdInterrupt() {
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt, os.Kill)
	go func() {
		si := <-s
		log.Alert("")
		log.DebugF("Got signal:%s", si)
		atomic.StoreUint32(&isCmdInterrupt, 1)
		Cancel()
		time.Sleep(time.Millisecond * 500)
		os.Exit(data.StatusUserCancel)
	}()
}
