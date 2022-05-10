package account

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"

	"github.com/qiniu/qshell/v2/iqshell/common/log"
)

type LoadInfo struct {
	AccountPath    string
	OldAccountPath string
	AccountDBPath  string
}

var info LoadInfo

// Load 保证 AccountPath、OldAccountPath、AccountDBPath 均不为空
func Load(i LoadInfo) *data.CodeError {
	if i.AccountDBPath == "" {
		return data.NewEmptyError().AppendDescF("empty account db path\n")
	}

	if i.AccountPath == "" {
		return data.NewEmptyError().AppendDescF("empty account path\n")
	}

	if i.OldAccountPath == "" {
		return data.NewEmptyError().AppendDescF("empty old account db path\n")
	}

	info = i

	log.Debug("account db path:" + info.AccountDBPath)
	log.Debug("account path:" + info.AccountPath)
	log.Debug("account old path:" + info.OldAccountPath)
	return nil
}
