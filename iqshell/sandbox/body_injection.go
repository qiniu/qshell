package sandbox

import (
	"encoding/json"
	"fmt"
)

// ApplyHTTPBodyFields 将 JSON 对象解析为 HTTP 注入字段，拒绝非字符串值和其他注入类型。
// 错误消息不包含原始 JSON，避免意外输出密钥。
func ApplyHTTPBodyFields(parts *InjectionParts, body, ifBody string) error {
	if body == "" && ifBody == "" {
		return nil
	}
	if parts.HTTP == nil {
		return fmt.Errorf("body and if-body are only supported for http injections")
	}
	parse := func(raw, name string) (*map[string]string, error) {
		if raw == "" {
			return nil, nil
		}
		var decoded map[string]any
		if err := json.Unmarshal([]byte(raw), &decoded); err != nil || decoded == nil {
			return nil, fmt.Errorf("%s must be a JSON object containing only string values", name)
		}
		if len(decoded) > 20 {
			return nil, fmt.Errorf("%s supports at most 20 fields", name)
		}
		fields := make(map[string]string, len(decoded))
		for key, value := range decoded {
			text, ok := value.(string)
			if !ok {
				return nil, fmt.Errorf("%s must be a JSON object containing only string values", name)
			}
			fields[key] = text
		}
		return &fields, nil
	}
	bodyFields, err := parse(body, "body")
	if err != nil {
		return err
	}
	matchFields, err := parse(ifBody, "if-body")
	if err != nil {
		return err
	}
	parts.HTTP.Body = bodyFields
	parts.HTTP.IfBody = matchFields
	return nil
}
