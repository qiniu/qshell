package workspace

import (
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

// config 配置信息
var cfg = &config.Config{}

// 获取之前需要先 Load
func GetConfig() config.Config {
	return *cfg
}
