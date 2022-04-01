package flow

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
)

type lineSeparateWorkCreator struct {
	separate    string
	creatorFunc func(items []string) (work Work, err *data.CodeError)
}

func (l *lineSeparateWorkCreator) Create(info string) (work Work, err *data.CodeError) {
	items := utils.SplitString(info, l.separate)
	return l.creatorFunc(items)
}

func NewLineSeparateWorkCreator(separate string, creatorFunc func(items []string) (work Work, err *data.CodeError)) WorkCreator {
	return &lineSeparateWorkCreator{
		separate:    separate,
		creatorFunc: creatorFunc,
	}
}
