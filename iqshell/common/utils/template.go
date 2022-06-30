package utils

import (
	"bytes"
	"encoding/json"
	"github.com/Masterminds/sprig/v3"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"html/template"
)

type Template struct {
	t *template.Template
}

func NewTemplate(templateString string) (*Template, *data.CodeError) {
	if t, err := template.New("QshellTemplate").Funcs(sprig.FuncMap()).Parse(templateString); err != nil {
		return nil, data.NewEmptyError().AppendDescF("create template by %s fail, %v", templateString, err)
	} else {
		return &Template{
			t: t,
		}, nil
	}
}

func (t *Template) Run(paramJson string) (string, *data.CodeError) {
	if t == nil {
		return "", alert.CannotEmptyError("Template", "")
	}

	var param interface{}
	paramMap := make(map[string]interface{})
	paramArray := make([]interface{}, 0, 0)
	if umErr := json.Unmarshal([]byte(paramJson), &paramMap); umErr == nil {
		param = paramMap
		log.DebugF("param is map:%v", param)
	} else if uaErr := json.Unmarshal([]byte(paramJson), &paramArray); uaErr == nil {
		param = paramArray
		log.DebugF("not map, %v", umErr)
		log.DebugF("param is array:%v", param)
	} else {
		param = paramJson
		log.DebugF("not map, %v", umErr)
		log.DebugF("not array, %v", uaErr)
		log.DebugF("param is string:%v", param)
	}

	buff := &bytes.Buffer{}
	if err := t.t.Execute(buff, param); err != nil {
		return "", data.ConvertError(err)
	} else {
		return buff.String(), nil
	}
}
