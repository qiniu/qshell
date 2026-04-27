package config

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// templateIDLineRegex 匹配未注释的 template_id 赋值行（允许前导空白）。
var templateIDLineRegex = regexp.MustCompile(`^\s*template_id\s*=`)

// WriteTemplateID 将 template_id 写入指定 TOML 文件。
// 若文件中已有未注释的 template_id 行，则替换其值；否则在文件头插入一行。
// 保留文件原有的权限、注释、字段顺序、缩进、换行风格和空白。
func WriteTemplateID(path, templateID string) error {
	stat, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat %s: %w", path, err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}

	// 识别原文件换行风格，回写时保留 CRLF / LF。
	lineEnd := "\n"
	if bytes.Contains(data, []byte("\r\n")) {
		lineEnd = "\r\n"
	}

	var out bytes.Buffer
	scanner := bufio.NewScanner(bytes.NewReader(data))
	replaced := false
	inRootTable := true
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if !isCommentLine(line) && strings.HasPrefix(trimmed, "[") {
			inRootTable = false
		}
		if inRootTable && !replaced && templateIDLineRegex.MatchString(line) && !isCommentLine(line) {
			// 保留原行的前导缩进
			idx := strings.Index(line, "template_id")
			indent := ""
			if idx > 0 {
				indent = line[:idx]
			}
			fmt.Fprintf(&out, "%stemplate_id = %q%s", indent, templateID, lineEnd)
			replaced = true
			continue
		}
		out.WriteString(line)
		out.WriteString(lineEnd)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan %s: %w", path, err)
	}

	if !replaced {
		// 未找到赋值行：在文件头插入
		var prefix bytes.Buffer
		fmt.Fprintf(&prefix, "template_id = %q%s", templateID, lineEnd)
		prefix.Write(out.Bytes())
		out = prefix
	}

	// 保持与原文件结尾换行一致
	final := out.Bytes()
	if len(data) > 0 && !bytes.HasSuffix(data, []byte(lineEnd)) && bytes.HasSuffix(final, []byte(lineEnd)) {
		final = final[:len(final)-len(lineEnd)]
	}

	if err := os.WriteFile(path, final, stat.Mode().Perm()); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

// isCommentLine 判断一行是否以 # 开头（允许前导空白）。
func isCommentLine(line string) bool {
	return strings.HasPrefix(strings.TrimLeft(line, " \t"), "#")
}
