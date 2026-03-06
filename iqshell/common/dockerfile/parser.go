// Package dockerfile 提供轻量级 Dockerfile 解析器。
//
// 将 Dockerfile 内容解析为结构化的指令列表，
// 支持 escape 指令、行续接、heredoc 语法和引号处理。
// 解析器遵循 buildkit 的词法处理规则，但不依赖 buildkit。
//
// 本包不包含任何 SDK 特定类型，使用方应将通用的 [Instruction] 类型
// 转换为各自的领域模型。
package dockerfile

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

// Instruction 表示一条解析后的 Dockerfile 指令。
type Instruction struct {
	// Name 是大写的指令名称（如 "RUN"、"COPY"、"ENV"）。
	Name string

	// Args 是指令名称后的原始参数字符串。
	Args string

	// Flags 是从指令中提取的 --flag=value 标记（如 --chown=user:group）。
	Flags map[string]string

	// Heredoc 在指令使用 heredoc 语法时包含 heredoc 主体内容。
	Heredoc string

	// Line 是该指令在原始 Dockerfile 中的行号（从 1 开始）。
	Line int
}

// ParseResult 保存 Dockerfile 解析的完整输出。
type ParseResult struct {
	// Instructions 是解析后的指令有序列表。
	Instructions []Instruction

	// Warnings 包含非致命的解析问题（如未知指令、多阶段构建）。
	Warnings []string

	// EscapeToken 是使用的转义字符（默认 '\\'，可覆盖为 '`'）。
	EscapeToken rune
}

// defaultEscapeToken 是默认的行续接字符。
const defaultEscapeToken = '\\'

// reHeredoc 匹配 heredoc 标记，如 <<EOF、0<<-EOF。
var reHeredoc = regexp.MustCompile(`^(\d*)<<(-?)\s*['"]?([a-zA-Z_]\w*)['"]?$`)

// Parse 将 Dockerfile 内容解析为 ParseResult。
func Parse(content string) (*ParseResult, error) {
	content = stripBOM(content)
	rawLines := strings.Split(content, "\n")

	escapeToken, rawLines := detectEscapeDirective(rawLines)
	lines := joinContinuationLines(rawLines, escapeToken)

	result := &ParseResult{EscapeToken: escapeToken}

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		name, rest := splitInstruction(line)
		name = strings.ToUpper(name)

		inst := Instruction{
			Name: name,
			Args: rest,
			Line: i + 1,
		}

		// 提取支持标志的指令中的 flags
		switch name {
		case "COPY", "ADD", "FROM":
			inst.Flags, inst.Args = extractFlags(rest)
		}

		// 处理支持 heredoc 的指令
		switch name {
		case "RUN", "COPY", "ADD":
			body, advance, err := maybeParseHeredoc(rest, lines[i+1:])
			if err != nil {
				return nil, fmt.Errorf("line %d: %w", i+1, err)
			}
			if advance > 0 {
				inst.Heredoc = body
				if name == "RUN" {
					inst.Args = body
				}
				i += advance
			}
		}

		// 对特定指令生成警告
		switch name {
		case "ONBUILD":
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("line %d: ONBUILD instruction is not supported and will be ignored", i+1))
		case "EXPOSE", "VOLUME", "LABEL", "STOPSIGNAL", "HEALTHCHECK", "SHELL",
			"FROM", "RUN", "COPY", "ADD", "WORKDIR", "USER", "ENV", "ARG",
			"CMD", "ENTRYPOINT", "MAINTAINER":
			// 已知指令，无需警告
		default:
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("line %d: unknown instruction %q will be ignored", i+1, name))
		}

		result.Instructions = append(result.Instructions, inst)
	}

	return result, nil
}

// ---------------------------------------------------------------------------
// 转义指令检测
// ---------------------------------------------------------------------------

// detectEscapeDirective 扫描 Dockerfile 顶部的 `# escape=X` 指令。
// 仅反斜杠和反引号是有效的转义字符。
// 指令必须出现在任何指令或非指令注释之前。
func detectEscapeDirective(lines []string) (rune, []string) {
	escapeToken := defaultEscapeToken

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || !strings.HasPrefix(trimmed, "#") {
			return escapeToken, lines
		}

		comment := strings.TrimSpace(strings.TrimPrefix(trimmed, "#"))
		if key, value, ok := strings.Cut(comment, "="); ok {
			key = strings.TrimSpace(strings.ToLower(key))
			value = strings.TrimSpace(value)
			if key == "escape" {
				if value == "\\" {
					escapeToken = '\\'
				} else if value == "`" {
					escapeToken = '`'
				}
				return escapeToken, append(lines[:i:i], lines[i+1:]...)
			}
			if key == "syntax" || key == "check" {
				continue
			}
		}

		// 非指令注释，停止扫描
		return escapeToken, lines
	}

	return escapeToken, lines
}

// ---------------------------------------------------------------------------
// 行续接
// ---------------------------------------------------------------------------

// joinContinuationLines 将以转义字符结尾的行合并。
// 续接块内的注释会被跳过（匹配 buildkit 行为）。
func joinContinuationLines(lines []string, escapeToken rune) []string {
	var result []string
	var current strings.Builder
	escStr := string(escapeToken)

	for _, line := range lines {
		trimmed := strings.TrimRight(line, " \t\r")

		if current.Len() > 0 {
			if t := strings.TrimSpace(trimmed); t == "" || strings.HasPrefix(t, "#") {
				continue
			}
		}

		if before, found := strings.CutSuffix(trimmed, escStr); found {
			if !isEscapedEscape(before, escapeToken) {
				current.WriteString(before)
				current.WriteString(" ")
				continue
			}
		}

		if current.Len() > 0 {
			current.WriteString(line)
			result = append(result, current.String())
			current.Reset()
		} else {
			result = append(result, line)
		}
	}
	if current.Len() > 0 {
		result = append(result, current.String())
	}
	return result
}

// isEscapedEscape 检查末尾的转义字符本身是否被转义。
// 使用 utf8.DecodeLastRuneInString 从末尾正确迭代 rune。
func isEscapedEscape(s string, escapeToken rune) bool {
	count := 0
	for len(s) > 0 {
		r, size := utf8.DecodeLastRuneInString(s)
		if r == escapeToken {
			count++
			s = s[:len(s)-size]
		} else {
			break
		}
	}
	return count%2 == 1
}

// ---------------------------------------------------------------------------
// Heredoc 支持
// ---------------------------------------------------------------------------

// maybeParseHeredoc 检查 rest 是否包含 heredoc 标记，
// 并消费后续行直到找到所有终止符。
func maybeParseHeredoc(rest string, followingLines []string) (string, int, error) {
	if !strings.Contains(rest, "<<") {
		return rest, 0, nil
	}

	words := strings.Fields(rest)
	var terminators []string
	for _, w := range words {
		if m := reHeredoc.FindStringSubmatch(w); m != nil {
			terminators = append(terminators, m[3])
		}
	}

	if len(terminators) == 0 {
		return rest, 0, nil
	}

	var bodies []string
	consumed := 0
	termIdx := 0

	for termIdx < len(terminators) && consumed < len(followingLines) {
		line := strings.TrimRight(followingLines[consumed], "\r\n")
		consumed++
		if strings.TrimSpace(line) == terminators[termIdx] {
			termIdx++
			continue
		}
		bodies = append(bodies, line)
	}

	if termIdx < len(terminators) {
		return "", 0, fmt.Errorf("unterminated heredoc: expected %q", terminators[termIdx])
	}

	return strings.Join(bodies, "\n"), consumed, nil
}

// ---------------------------------------------------------------------------
// 指令解析辅助函数
// ---------------------------------------------------------------------------

// splitInstruction 将一行拆分为指令名称和剩余部分。
func splitInstruction(line string) (string, string) {
	parts := strings.SplitN(line, " ", 2)
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], strings.TrimSpace(parts[1])
}

// extractFlags 从参数开头提取 --flag=value 标记。
// 仅提取前导的 flag 标记，遇到第一个非 flag 参数后停止提取。
// 返回提取的标志映射和剩余的参数字符串。
func extractFlags(args string) (map[string]string, string) {
	flags := make(map[string]string)
	fields := strings.Fields(args)

	i := 0
	for i < len(fields) {
		f := fields[i]
		if !strings.HasPrefix(f, "--") {
			break
		}
		if key, value, ok := strings.Cut(f, "="); ok {
			flags[strings.TrimPrefix(key, "--")] = value
		} else {
			flags[strings.TrimPrefix(f, "--")] = ""
		}
		i++
	}

	return flags, strings.Join(fields[i:], " ")
}

// ---------------------------------------------------------------------------
// 引号感知的值解析（导出供复用）
// ---------------------------------------------------------------------------

// ParseEnvValues 解析 ENV 风格的 KEY=VALUE 对，正确处理引号。
// 支持双引号（带转义处理）、单引号（字面量）和无引号值。
// 返回扁平的键值切片：["K1", "V1", "K2", "V2", ...]。
func ParseEnvValues(rest string, escapeToken rune) ([]string, error) {
	if rest == "" {
		return nil, fmt.Errorf("empty ENV instruction")
	}

	// 检查是否使用 KEY=VALUE 格式
	if strings.Contains(strings.Fields(rest)[0], "=") {
		return parseEnvKeyValue(rest, escapeToken), nil
	}

	// 旧格式：ENV KEY VALUE
	key, value, _ := strings.Cut(rest, " ")
	return []string{key, strings.TrimSpace(value)}, nil
}

// parseEnvKeyValue 解析 KEY=VALUE 对，正确处理引号和转义。
func parseEnvKeyValue(rest string, escapeToken rune) []string {
	var result []string
	pos := 0

	for pos < len(rest) {
		for pos < len(rest) && (rest[pos] == ' ' || rest[pos] == '\t') {
			pos++
		}
		if pos >= len(rest) {
			break
		}

		eqIdx := strings.IndexByte(rest[pos:], '=')
		if eqIdx < 0 {
			result = append(result, rest[pos:], "")
			break
		}
		key := rest[pos : pos+eqIdx]
		pos += eqIdx + 1

		value, newPos := parseQuotedValue(rest, pos, escapeToken)
		pos = newPos
		result = append(result, key, value)
	}

	return result
}

// parseQuotedValue 从 pos 位置开始提取值，处理双引号、单引号和无引号值。
func parseQuotedValue(s string, pos int, escapeToken rune) (string, int) {
	if pos >= len(s) {
		return "", pos
	}

	switch ch := s[pos]; ch {
	case '"':
		return parseDoubleQuoted(s, pos+1, escapeToken)
	case '\'':
		return parseSingleQuoted(s, pos+1)
	default:
		return parseUnquoted(s, pos)
	}
}

// parseDoubleQuoted 解析双引号字符串，处理转义。
func parseDoubleQuoted(s string, pos int, escapeToken rune) (string, int) {
	var value strings.Builder
	for pos < len(s) {
		r, size := utf8.DecodeRuneInString(s[pos:])
		if r == '"' {
			return value.String(), pos + size
		}
		if r == escapeToken && pos+size < len(s) {
			next, nextSize := utf8.DecodeRuneInString(s[pos+size:])
			value.WriteRune(next)
			pos += size + nextSize
			continue
		}
		value.WriteRune(r)
		pos += size
	}
	return value.String(), pos
}

// parseSingleQuoted 解析单引号字符串（所有字符均为字面量）。
func parseSingleQuoted(s string, pos int) (string, int) {
	end := strings.IndexByte(s[pos:], '\'')
	if end < 0 {
		return s[pos:], len(s)
	}
	return s[pos : pos+end], pos + end + 1
}

// parseUnquoted 解析无引号值，遇到空白字符结束。
func parseUnquoted(s string, pos int) (string, int) {
	start := pos
	for pos < len(s) && s[pos] != ' ' && s[pos] != '\t' {
		pos++
	}
	return s[start:pos], pos
}

// ParseCommand 解析 CMD 或 ENTRYPOINT 参数。
// 支持 exec 格式 ["cmd", "arg1", ...] 和 shell 格式。
// exec 格式中含 shell 特殊字符的参数会被单引号包裹以防止 bash 解析。
func ParseCommand(rest string) string {
	rest = strings.TrimSpace(rest)
	if rest == "" {
		return ""
	}

	if inner, ok := strings.CutPrefix(rest, "["); ok {
		inner = strings.TrimSuffix(strings.TrimSpace(inner), "]")
		var parts []string
		for item := range strings.SplitSeq(inner, ",") {
			item = strings.TrimSpace(item)
			item = strings.Trim(item, "\"'")
			if item != "" {
				parts = append(parts, shellQuote(item))
			}
		}
		return strings.Join(parts, " ")
	}

	return rest
}

// shellQuote 对含 shell 特殊字符的字符串用单引号包裹。
// 若字符串仅含安全字符则原样返回。
func shellQuote(s string) string {
	for _, c := range s {
		if !isShellSafe(c) {
			// 用单引号包裹，其中已有的单引号用 '\'' 转义
			return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
		}
	}
	return s
}

// isShellSafe 判断字符是否对 shell 安全（无需引号）。
func isShellSafe(c rune) bool {
	if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' {
		return true
	}
	switch c {
	case '-', '_', '.', '/', ':', ',', '+', '=', '@', '%':
		return true
	}
	return false
}

// StripHeredocMarkers 从字符串中移除 <<WORD 标记。
func StripHeredocMarkers(s string) string {
	words := strings.Fields(s)
	var filtered []string
	for _, w := range words {
		if reHeredoc.MatchString(w) {
			continue
		}
		filtered = append(filtered, w)
	}
	return strings.Join(filtered, " ")
}

// ---------------------------------------------------------------------------
// 辅助函数
// ---------------------------------------------------------------------------

// stripBOM 移除内容开头的 UTF-8 BOM。
func stripBOM(s string) string {
	return strings.TrimPrefix(s, "\xef\xbb\xbf")
}
