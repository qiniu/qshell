package alert

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"
)

func CannotEmptyError(words string, suggest string) *data.CodeError {
	return data.NewError(data.ErrorCodeParamNotExist, CannotEmpty(words, suggest))
}

func Error(desc string, suggest string) *data.CodeError {
	return data.NewError(data.ErrorCodeParamNotExist, Description(desc, suggest))
}

func CannotEmpty(words string, suggest string) string {
	desc := words
	if len(words) > 0 {
		desc += " can't be empty"
	}
	return Description(desc, suggest)
}

func Description(desc string, suggest string) string {
	ret := desc
	if len(suggest) > 0 {
		ret += ", you can do like this:\n" + suggest
	}
	return ret
}
