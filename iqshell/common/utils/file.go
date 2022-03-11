package utils

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
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

func UnMarshalFromFile(filePath string, v interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		return errors.New("unmarshal: open file error:" + err.Error())
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return errors.New("unmarshal: read file error:" + err.Error())
	}

	err = json.Unmarshal(data, v)
	if err != nil {
		return errors.New("unmarshal: unmarshal error:" + err.Error())
	}

	return nil
}

func IsNetworkSource(filePath string) bool {
	return strings.HasPrefix(filePath, "http://") || strings.HasPrefix(filePath, "https://")
}

func FileSize(filePath string) (fileSize int64, err error) {
	if IsNetworkSource(filePath) {
		return NetworkFileLength(filePath)
	} else {
		return LocalFileSize(filePath)
	}
}

func LocalFileSize(filePath string) (fileSize int64, err error) {
	fileStatus, err := os.Stat(filePath)
	if err != nil {
		err = errors.New("get file size: get status error:" + err.Error())
		return
	}

	fileSize = fileStatus.Size()
	return
}

func NetworkFileLength(srcResUrl string) (fileSize int64, err error) {
	resp, respErr := http.Head(srcResUrl)
	if respErr != nil {
		err = fmt.Errorf("New head request failed, %s", respErr.Error())
		return
	}
	defer resp.Body.Close()

	contentLength := resp.Header.Get("Content-Length")
	if contentLength == "" {
		err = errors.New("head request with no Content-Length found error")
		return
	}

	fileSize, _ = strconv.ParseInt(contentLength, 10, 64)

	return
}

func FileLineCounts(filePath string) (count int64, err error) {
	fp, openErr := os.Open(filePath)
	if openErr != nil {
		return 0, openErr
	}
	defer fp.Close()

	bScanner := bufio.NewScanner(fp)
	for bScanner.Scan() {
		count += 1
	}
	return
}

func CreateFileIfNotExist(path string) error {
	if exist, err := ExistFile(path); err == nil && exist {
		return nil
	}
	return CreateFileDirIfNotExist(path)
}

func CreateFileDirIfNotExist(path string) error {
	dir := filepath.Dir(path)
	if err := CreateDirIfNotExist(dir); err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	defer f.Close()
	return err
}

func ExistFile(path string) (bool, error) {
	if s, err := os.Stat(path); err == nil {
		return !s.IsDir(), nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

func CreateDirIfNotExist(path string) error {
	if exist, err := ExistDir(path); err == nil && exist {
		return nil
	}
	return os.MkdirAll(path, os.ModePerm)
}

func ExistDir(path string) (bool, error) {
	if s, err := os.Stat(path); err == nil {
		return s.IsDir(), nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}
