package workspace

import (
	"errors"
	"github.com/qiniu/qshell/v2/iqshell/utils"
	"path/filepath"
)

const (
	workspaceName         = ".qshell"
	usersDirName          = "users"
	usersDBName           = "users.db"
	usersWorkspaceDirName = "workspace"
	taskDirName           = "task"
	taskDBName            = "task.db"
	configFileName        = "config.json"
)

var (
	workspace = func() string {
		home := utils.GetHomePath()
		if len(home) == 0 {
			return ""
		}
		return filepath.Join(home, workspaceName)
	}()
)

// 检查工作目录
func Load() (err error) {
	if len(workspace) == 0 {
		err = errors.New("can't get home dir")
		return
	}

	err = utils.CreateDirIfNotExist(workspace)

	return
}
