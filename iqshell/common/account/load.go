package account

import (
	"fmt"

	"github.com/qiniu/qshell/v2/iqshell/common/log"
)

type LoadInfo struct {
	AccountPath    string
	OldAccountPath string
	AccountDBPath  string
}

var info LoadInfo
// Load 保证 AccountPath、OldAccountPath、AccountDBPath 均不为空
func Load(i LoadInfo) error {
	if info.AccountDBPath == "" {
		return fmt.Errorf("empty account db path\n")
	}

	if info.AccountPath == "" {
		return fmt.Errorf("empty account path\n")
	}

	if info.OldAccountPath == "" {
		return fmt.Errorf("empty old account db path\n")
	}

	info = i

	log.Debug("account db path:" + info.AccountDBPath)
	log.Debug("account path:" + info.AccountPath)
	log.Debug("account old path:" + info.OldAccountPath)
	return nil
}
