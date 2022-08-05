package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	KB = 1024
	MB = 1024 * KB
	GB = 1024 * MB
	TB = 1024 * GB
)

// FormatFileSize 转化文件大小到人工可读的字符串，以相应的单位显示
func FormatFileSize(size int64) (result string) {
	if size >= TB {
		result = fmt.Sprintf("%.2fTB", float64(size)/float64(TB))
	} else if size >= GB {
		result = fmt.Sprintf("%.2fGB", float64(size)/float64(GB))
	} else if size >= MB {
		result = fmt.Sprintf("%.2fMB", float64(size)/float64(MB))
	} else if size >= KB {
		result = fmt.Sprintf("%.2fKB", float64(size)/float64(KB))
	} else {
		result = fmt.Sprintf("%dB", size)
	}
	return
}

func MarshalToFile(filePath string, v interface{}) *data.CodeError {
	if v == nil {
		return nil
	}

	if err := os.Remove(filePath); err != nil && os.IsExist(err) {
		return data.NewEmptyError().AppendDesc("marshal: delete origin file").AppendError(err)
	}

	if d, mErr := json.Marshal(v); mErr != nil {
		return data.NewEmptyError().AppendDesc("marshal: marshal").AppendError(mErr)
	} else if wErr := os.WriteFile(filePath, d, os.ModePerm); wErr != nil {
		return data.NewEmptyError().AppendDesc("marshal: write file").AppendError(mErr)
	} else {
		return nil
	}
}

func UnMarshalFromFile(filePath string, v interface{}) *data.CodeError {
	if v == nil {
		return nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		return data.NewEmptyError().AppendDesc("unmarshal: open file").AppendError(err)
	}

	d, err := ioutil.ReadAll(file)
	if err != nil {
		return data.NewEmptyError().AppendDesc("unmarshal: read file").AppendError(err)
	}

	err = json.Unmarshal(d, v)
	if err != nil {
		return data.NewEmptyError().AppendDesc("unmarshal: unmarshal").AppendError(err)
	}

	return nil
}

func IsNetworkSource(filePath string) bool {
	return strings.HasPrefix(filePath, "http://") || strings.HasPrefix(filePath, "https://")
}

func IsFileMatchFileSize(filePath string, fileSize int64) (match bool, err *data.CodeError) {
	if size, e := FileSize(filePath); e != nil {
		return false, e
	} else if size != fileSize {
		return false, data.NewEmptyError().AppendDescF("size don't match, except:%d but:%d", fileSize, size)
	} else {
		return true, nil
	}
}

func FileSize(filePath string) (fileSize int64, err *data.CodeError) {
	if IsNetworkSource(filePath) {
		return NetworkFileLength(filePath)
	} else {
		return LocalFileSize(filePath)
	}
}

func LocalFileSize(filePath string) (int64, *data.CodeError) {
	fileStatus, err := os.Stat(filePath)
	if err != nil {
		return 0, data.NewEmptyError().AppendDescF("get file size: get status error:%v", err)
	}
	return fileStatus.Size(), nil
}

func NetworkFileLength(srcResUrl string) (fileSize int64, err *data.CodeError) {
	resp, respErr := http.Head(srcResUrl)
	if respErr != nil {
		err = data.NewEmptyError().AppendDescF("New head request failed, %s", respErr.Error())
		return
	}
	defer resp.Body.Close()

	contentLength := resp.Header.Get("Content-Length")
	if contentLength == "" {
		err = data.NewEmptyError().AppendDesc("head request with no Content-Length found error")
		return
	}

	fileSize, _ = strconv.ParseInt(contentLength, 10, 64)

	return
}

func IsFileMatchFileModifyTime(filePath string, modifyTime int64) (match bool, err *data.CodeError) {
	if time, e := FileModify(filePath); e != nil {
		return false, e
	} else if time != modifyTime {
		return false, data.NewEmptyError().AppendDescF("modifyTime don't match, except:%d but:%d", modifyTime, time)
	} else {
		return true, nil
	}
}

func FileModify(filePath string) (int64, *data.CodeError) {
	if IsNetworkSource(filePath) {
		return NetworkFileModify(filePath)
	} else {
		return LocalFileModify(filePath)
	}
}

func NetworkFileModify(filePath string) (int64, *data.CodeError) {
	fileStatus, err := os.Stat(filePath)
	if err != nil {
		return 0, data.NewEmptyError().AppendDescF("get file : get status error:%v", err)
	}
	return fileStatus.ModTime().Unix(), nil
}

func LocalFileModify(filePath string) (int64, *data.CodeError) {
	fileStatus, err := os.Stat(filePath)
	if err != nil {
		return 0, data.NewEmptyError().AppendDescF("get file : get status error:%v", err)
	}
	return fileStatus.ModTime().Unix(), nil
}

func FileLineCounts(filePath string) (count int64, err *data.CodeError) {
	fp, openErr := os.Open(filePath)
	if openErr != nil {
		return 0, data.NewEmptyError().AppendError(openErr)
	}
	defer fp.Close()

	bScanner := bufio.NewScanner(fp)
	for bScanner.Scan() {
		count += 1
	}
	return
}

func CreateFileIfNotExist(path string) *data.CodeError {
	if exist, err := ExistFile(path); err == nil && exist {
		return nil
	}
	return CreateFileDirIfNotExist(path)
}

func CreateFileDirIfNotExist(path string) *data.CodeError {
	dir := filepath.Dir(path)
	if err := CreateDirIfNotExist(dir); err != nil {
		return err
	}
	return nil
}

func ExistFile(path string) (bool, *data.CodeError) {
	if s, err := os.Stat(path); err == nil {
		return !s.IsDir(), nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, data.NewEmptyError().AppendError(err)
	}
}

func CreateDirIfNotExist(path string) *data.CodeError {
	if exist, err := ExistDir(path); err == nil && exist {
		return nil
	}
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return data.NewEmptyError().AppendError(err)
	} else {
		return nil
	}
}

func ExistDir(path string) (bool, *data.CodeError) {
	if s, err := os.Stat(path); err == nil {
		return s.IsDir(), nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, data.NewEmptyError().AppendError(err)
	}
}
