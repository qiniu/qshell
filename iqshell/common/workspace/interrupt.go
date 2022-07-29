package workspace

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
)

var (
	// 程序是否退出
	isCmdInterrupt  uint32 = 0
	locker          sync.Mutex
	cancelObservers = make([]func(s os.Signal), 0)
)

func AddCancelObserver(observer func(s os.Signal)) {
	if observer == nil {
		return
	}

	locker.Lock()
	cancelObservers = append(cancelObservers, observer)
	locker.Unlock()
}

func notifyCancelSignalToObservers(s os.Signal) {
	locker.Lock()
	for _, observer := range cancelObservers {
		observer(s)
	}
	locker.Unlock()
}

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
		data.SetCmdStatusUserCancel()
		Cancel()
		notifyCancelSignalToObservers(si)
	}()
}
