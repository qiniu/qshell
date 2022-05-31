package bucket

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"strconv"
	"strings"
	"sync"
)

const (
	listObjectFieldsKey         = "Key"
	listObjectFieldsHash        = "Hash"
	listObjectFieldsFileSize    = "FileSize"
	listObjectFieldsPutTime     = "PutTime"
	listObjectFieldsMimeType    = "MimeType"
	listObjectFieldsStorageType = "StorageType"
	listObjectFieldsEndUser     = "EndUser"
)

var listObjectFields = []string{
	listObjectFieldsKey,
	listObjectFieldsFileSize,
	listObjectFieldsHash,
	listObjectFieldsPutTime,
	listObjectFieldsMimeType,
	listObjectFieldsStorageType,
	listObjectFieldsEndUser,
}

type ListLineParser struct {
	mu          sync.Mutex
	isFirstLine bool
	fields      []string // 输入的字段信息
}

func NewListLineParser() *ListLineParser {
	return &ListLineParser{
		mu:          sync.Mutex{},
		isFirstLine: true,
		fields:      nil,
	}
}

func (p *ListLineParser) Parse(items []string) (*ListObject, *data.CodeError) {
	if p.isFirstLine {
		p.isFirstLine = false
		isKeys, fields := getKeyItems(items)
		p.fields = fields
		if isKeys {
			return nil, data.NewError(data.ErrorCodeParamMissing, "This is key line")
		} else {
			return p.convert(items)
		}
	}
	return p.convert(items)
}

func (p *ListLineParser) convert(items []string) (*ListObject, *data.CodeError) {
	if len(items) == 0 {
		return nil, data.NewError(0, "empty line")
	}

	object := &ListObject{}
	for index, field := range p.fields {
		if index >= len(items) {
			return object, nil
		}

		if err := listObjectSetFieldWithStringValue(object, field, items[index]); err != nil {
			return nil, err
		}
	}
	return object, nil
}

type ListLineCreator struct {
	Fields   []string // 需要输出的字段
	Sep      string   // 分隔符
	Readable bool     // 是否可读
}

func (l *ListLineCreator)Create(object *ListObject) string {
	return listObjectDescWithFields(object, l.Fields, l.Sep, l.Readable)
}

func getKeyItems(items []string) (bool, []string) {
	if len(items) == 0 {
		return false, listObjectFields
	}

	var fields []string
	var isKeys = true
	for _, item := range items {
		field := listObjectField(item)
		if len(field) == 0 {
			isKeys = false
			break
		} else {
			fields = append(fields, field)
		}
	}

	if isKeys {
		return true, fields
	} else {
		return false, listObjectFields
	}
}

func listObjectField(field string) string {
	for _, f := range listObjectFields {
		if strings.EqualFold(field, f) {
			return f
		}
	}
	return ""
}

func listObjectDescWithFields(object *ListObject, fields []string, outputFieldsSep string, readable bool) string {
	var values []string
	for _, field := range fields {
		values = append(values, listObjectFieldStringValue(object, field, readable))
	}
	return strings.Join(values, outputFieldsSep)
}

func listObjectFieldStringValue(object *ListObject, field string, readable bool) string {
	if object == nil {
		return ""
	}

	var value interface{}
	switch field {
	case listObjectFieldsKey:
		value = object.Key
		break
	case listObjectFieldsFileSize:
		if readable {
			value = utils.BytesToReadable(object.Fsize)
		} else {
			value = object.Fsize
		}
		break
	case listObjectFieldsHash:
		value = object.Hash
		break
	case listObjectFieldsPutTime:
		value = object.PutTime
		break
	case listObjectFieldsMimeType:
		value = object.MimeType
		break
	case listObjectFieldsStorageType:
		value = object.Type
		break
	case listObjectFieldsEndUser:
		value = object.EndUser
		break
	default:
	}
	return fmt.Sprintf("%v", value)
}

func listObjectSetFieldWithStringValue(object *ListObject, field string, value string) *data.CodeError {
	if object == nil {
		return nil
	}

	var err *data.CodeError
	switch field {
	case listObjectFieldsKey:
		object.Key = value
		break
	case listObjectFieldsFileSize:
		if i, e := strconv.ParseInt(value, 10, 64); e != nil {
			err = data.NewEmptyError().AppendDescF("get file size error:%v", e)
		} else {
			object.Fsize = i
		}
		break
	case listObjectFieldsHash:
		object.Hash = value
		break
	case listObjectFieldsPutTime:
		if i, e := strconv.ParseInt(value, 10, 64); e != nil {
			err = data.NewEmptyError().AppendDescF("get put time error:%v", e)
		} else {
			object.PutTime = i
		}
		break
	case listObjectFieldsMimeType:
		object.MimeType = value
		break
	case listObjectFieldsStorageType:
		if i, e := strconv.Atoi(value); e != nil {
			err = data.NewEmptyError().AppendDescF("get storage type error:%v", e)
		} else {
			object.Type = i
		}
		break
	case listObjectFieldsEndUser:
		object.EndUser = value
		break
	default:
	}

	return err
}
