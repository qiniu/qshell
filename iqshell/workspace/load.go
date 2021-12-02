package workspace

import (
	"errors"
	"github.com/qiniu/qshell/v2/iqshell/config"
	"github.com/qiniu/qshell/v2/iqshell/utils"
	"path/filepath"
)

type Option func(w *workspace)

func Workspace(path string) Option {
	return func(w *workspace) {
		w.workspace = path
	}
}

func UserConfigPath(path string) Option {
	return func(w *workspace) {
		w.userConfigPath = path
	}
}

// 加载工作环境
func Load(options ...Option) (err error) {
	ws := &workspace{}

	// 设置配置
	for _, option := range options {
		option(ws)
	}

	// 检查工作目录
	if len(ws.workspace) == 0 {
		err = errors.New("can't get home dir")
		return
	}
	err = utils.CreateDirIfNotExist(ws.workspace)
	if err != nil {
		return
	}

	// 设置配置文件路径
	config.Load(config.UserConfigPath(ws.userConfigPath), config.GlobalConfigPath(ws.globalConfigPath))

	// 加载配置
	cfg.Merge(config.GetUser())
	cfg.Merge(config.GetGlobal())
	cfg.Merge(DefaultConfig())

	return
}

type workspace struct {
	workspace        string
	userConfigPath   string
	globalConfigPath string
}

func (w *workspace) init() {
	home, err := utils.GetHomePath()
	if len(home) == 0 || err != nil {
		return
	}

	w.workspace = filepath.Join(home, workspaceName)
	w.userConfigPath = filepath.Join(w.workspace, usersDirName, configFileName)
	w.globalConfigPath = filepath.Join(home, configFileName)
}
