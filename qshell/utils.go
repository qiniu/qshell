package qshell

import (
	"bufio"
	"os"
)

func GetFileLineCount(filePath string) (totalCount int64) {
	fp, openErr := os.Open(filePath)
	if openErr != nil {
		return
	}
	defer fp.Close()

	bScanner := bufio.NewScanner(fp)
	for bScanner.Scan() {
		totalCount += 1
	}
	return
}
