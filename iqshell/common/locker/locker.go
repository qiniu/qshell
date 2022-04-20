package locker

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"os"
	"path/filepath"
	"sync"
)

var lockerPath string
var locker sync.Mutex

func SetLockerPath(path string)  {
	if len(path) == 0 {
		return
	}

	locker.Lock()
	defer locker.Unlock()

	lockerPath = filepath.Join(path, ".lock")
}

func IsLock() bool {
	locker.Lock()
	defer locker.Unlock()

	if len(lockerPath) == 0 {
		return false
	}

	if exist, e := utils.ExistFile(lockerPath); e != nil {
		return true
	} else {
		return exist
	}
}

func Lock() *data.CodeError {
	locker.Lock()
	defer locker.Unlock()

	if len(lockerPath) == 0 {
		return nil
	}

	process := fmt.Sprintf("%d", os.Getpid())
	err := os.WriteFile(lockerPath, []byte(process), os.ModePerm)
	log.DebugF("job lock, process:%s path:%s err:%v", process, lockerPath, err)
	return data.ConvertError(err)
}

func UnLock() *data.CodeError {
	locker.Lock()
	defer locker.Unlock()

	if len(lockerPath) == 0 {
		return nil
	}

	err := os.Remove(lockerPath)
	log.DebugF("job unlock, path:%s err:%v", lockerPath, err)
	return data.ConvertError(err)
}

func LockProcess() string {
	locker.Lock()
	defer locker.Unlock()

	if info, e := os.ReadFile(lockerPath); e != nil {
		log.ErrorF("get lock process error:%v", e)
		return ""
	} else {
		return string(info)
	}
}