package locker

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
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

func TryLock() *data.CodeError {
	SetLockerPath(workspace.GetJobDir())
	if IsLock() {
		return data.NewEmptyError().AppendDescF("job doing by process:%s, If it is confirmed that this process doesn't exist or is not a qshell process, you can delete the .lock file in this folder:%s and try again",
			LockProcess(), workspace.GetJobDir())
	}

	if e := Lock(); e != nil {
		return data.NewEmptyError().AppendDescF("job lock error:%v", e)
	}
	return nil
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

func TryUnlock() *data.CodeError {
	if e := Unlock(); e != nil {
		return data.NewEmptyError().AppendDescF("job unlock error:%v", e)
	}
	return nil
}

func Unlock() *data.CodeError {
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