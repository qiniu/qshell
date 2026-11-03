package operations

// splitInlineInjectionFields 按顶层逗号拆分规约，保留 JSON 对象中的逗号和转义引号。
func splitInlineInjectionFields(spec string) []string {
	var fields []string
	start, depth := 0, 0
	quoted, escaped := false, false
	for i, ch := range spec {
		if quoted {
			if escaped {
				escaped = false
			} else if ch == '\\' {
				escaped = true
			} else if ch == '"' {
				quoted = false
			}
			continue
		}
		switch ch {
		case '"':
			if depth > 0 {
				quoted = true
			}
		case '{', '[':
			depth++
		case '}', ']':
			if depth > 0 {
				depth--
			}
		case ',':
			if depth == 0 {
				fields = append(fields, spec[start:i])
				start = i + 1
			}
		}
	}
	return append(fields, spec[start:])
}
