package ip

import "github.com/qiniu/qshell/v2/iqshell/common/data"

type groupIPParser struct {
	parserList []Parser
}

func DefaultParser() Parser {
	return NewGroupParser(NewAliIPParser())
}

func NewGroupParser(parsers ...Parser) Parser {
	return &groupIPParser{
		parserList: parsers,
	}
}

var _ Parser = (*groupIPParser)(nil)

func (g *groupIPParser) Parse(ip string) (result ParserResult, err *data.CodeError) {
	if g == nil || len(g.parserList) == 0 {
		return nil, data.NewEmptyError().AppendDesc("no group parser")
	}
	for _, parser := range g.parserList {
		result, err = parser.Parse(ip)
		if err == nil && result != nil {
			break
		}
	}
	return
}
