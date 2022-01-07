package test

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"os"
	"path/filepath"
)

func RootPath() (string, error) {
	if r, err := homedir.Dir(); err != nil {
		return "", err
	} else {
		return filepath.Join(r, "qshell_test"), nil
	}
}

func CreateFileWithContent(fileName, content string) (string, error) {
	rootPath, err := RootPath()
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(rootPath, "file")
	err = os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		return "", err
	}

	filePath = filepath.Join(filePath, fileName)
	_ = RemoveFile(filePath)
	if f, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0600); err != nil {
		return "", err
	} else {
		_, err = f.Write([]byte(content))
		return filePath, err
	}
}

func CreateTempFile(size int) (string, error) {
	rootPath, err := RootPath()
	if err != nil {
		return "", err
	}

	fileName := fmt.Sprintf("%vK", size)
	fileName = filepath.Join(rootPath, "temp", fileName)
	fi, err := os.Stat(fileName)

	if err == nil {
		if !fi.IsDir() {
			if fi.Size() == int64(size) {
				return fileName, nil
			} else {
				if err = RemoveFile(fileName); err != nil {
					return "", err
				}
			}
		}
	} else if !os.IsNotExist(err) {
		return "", err
	}

	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return "", err
	}
	defer f.Close()

	oneKB := [1024]byte{8, 8, 8, 8, 8, 8, 8}
	written := 0
	for written < size {
		if _, err = f.Write(oneKB[:]); err != nil {
			break
		}
		written += 1
	}

	if err != nil {
		_ = RemoveFile(fileName)
		fileName = ""
	}

	return fileName, err
}

func RemoveFile(filePath string) error {
	return os.Remove(filePath)
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

func ExistDir(path string) (bool, error) {
	if s, err := os.Stat(path); err == nil {
		return s.IsDir(), nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}
