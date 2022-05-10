package utils

import (
	"github.com/mitchellh/go-homedir"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
)

func GetHomePath() (string, *data.CodeError) {
	if path, e := homedir.Dir(); e != nil {
		return "", data.NewEmptyError().AppendError(e)
	} else {
		return path, nil
	}
}
