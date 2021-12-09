package workspace

import (
	"context"
	"github.com/qiniu/qshell/v2/iqshell/common/config"
)

const (
	workspaceName         = ".qshell"
	usersDirName          = "users"
	usersDBName           = "account.db"
	currentUserFileName   = "account.json"
	oldUserFileName       = "old_account.json"
	usersWorkspaceDirName = "workspace"
	taskDirName           = "task"
	taskDBName            = "task.db"
	configFileName        = ".qshell.json"
)

var (
	// config 配置信息
	cfg = &config.Config{}

	// 工作路径
	workspacePath = ""

	cancelFunc func()
)

// 获取之前需要先 Load
func GetConfig() config.Config {
	return *cfg
}

func GetWorkspace() string {
	return workspacePath
}

func GetContext() context.Context {
	ctx := context.Background()
	ctx, cancelFunc = context.WithCancel(ctx)
	return ctx
}

func Cancel() {
	if cancelFunc != nil {
		cancelFunc()
	}
}