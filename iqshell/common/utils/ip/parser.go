package ip

import "github.com/qiniu/qshell/v2/iqshell/common/data"

type ParserResult interface {
}

type Parser interface {
	Parse(ip string) (ParserResult, *data.CodeError)
}
