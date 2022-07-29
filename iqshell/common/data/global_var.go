package data

import "sync"

// 此处需要明确出值
const (
	StatusOK         = 0 // process success
	StatusError      = 1 // process error
	StatusUserCancel = 2 // 用户取消
)

var (
	mu        sync.Mutex
	cmdStatus = StatusOK
)

func SetCmdStatus(status int) {
	mu.Lock()
	defer mu.Unlock()
	if cmdStatus == StatusUserCancel {
		return
	}

	cmdStatus = status
}

func SetCmdStatusError() {
	SetCmdStatus(StatusError)
}

func SetCmdStatusUserCancel() {
	SetCmdStatus(StatusUserCancel)
}

func GetCmdStatus() int {
	mu.Lock()
	defer mu.Unlock()

	return cmdStatus
}
