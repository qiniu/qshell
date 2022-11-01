package utils

import "bytes"

var nl = []byte{'\n'}

func JsonDataTrimComments(data []byte) (data1 []byte) {
	lines := bytes.Split(data, nl)
	for k, line := range lines {
		lines[k] = trimCommentsLine(line)
	}
	return bytes.Join(lines, nl)
}

func trimCommentsLine(line []byte) []byte {

	var newLine []byte
	var i, quoteCount int
	lastIdx := len(line) - 1
	for i = 0; i <= lastIdx; i++ {
		if line[i] == '\\' {
			if i != lastIdx && (line[i+1] == '\\' || line[i+1] == '"') {
				newLine = append(newLine, line[i], line[i+1])
				i++
				continue
			}
		}
		if line[i] == '"' {
			quoteCount++
		}
		if line[i] == '#' {
			if quoteCount%2 == 0 {
				break
			}
		}
		if line[i] == '/' && i < lastIdx && line[i+1] == '/' {
			if quoteCount%2 == 0 {
				break
			}
		}
		newLine = append(newLine, line[i])
	}
	return newLine
}
