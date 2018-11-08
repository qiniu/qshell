package qshell

import (
	"archive/zip"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"os"
	"path/filepath"
	"unicode/utf8"
)

func gbk2Utf8(text string) (string, error) {
	var gDecoder = simplifiedchinese.GBK.NewDecoder()
	utf8Dst := make([]byte, len(text)*3)
	_, _, err := gDecoder.Transform(utf8Dst, []byte(text), true)
	if err != nil {
		return "", nil
	}
	gDecoder.Reset()
	utf8Bytes := make([]byte, 0)
	for _, b := range utf8Dst {
		if b != 0 {
			utf8Bytes = append(utf8Bytes, b)
		}
	}
	return string(utf8Bytes), nil
}

func Unzip(zipFilePath string, unzipPath string) (err error) {
	zipReader, zipErr := zip.OpenReader(zipFilePath)
	if zipErr != nil {
		err = fmt.Errorf("Open zip file error, %s", zipErr)
		return
	}
	defer zipReader.Close()

	zipFiles := zipReader.File

	//list dir
	for _, zipFile := range zipFiles {
		fileInfo := zipFile.FileHeader.FileInfo()
		fileName := zipFile.FileHeader.Name

		//check charset utf8 or gbk
		if !utf8.Valid([]byte(fileName)) {
			fileName, err = gbk2Utf8(fileName)
			if err != nil {
				err = errors.New("Unsupported filename encoding")
				continue
			}
		}

		fullPath := filepath.Join(unzipPath, fileName)
		if fileInfo.IsDir() {
			logs.Debug("Mkdir", fullPath)
			mErr := os.MkdirAll(fullPath, 0775)
			if mErr != nil {
				err = fmt.Errorf("Mkdir error, %s", mErr)
				continue
			}
		}
	}

	//list file
	for _, zipFile := range zipFiles {
		fileInfo := zipFile.FileHeader.FileInfo()
		fileName := zipFile.FileHeader.Name

		//check charset utf8 or gbk
		if !utf8.Valid([]byte(fileName)) {
			fileName, err = gbk2Utf8(fileName)
			if err != nil {
				err = errors.New("Unsupported filename encoding")
				continue
			}
		}

		fullPath := filepath.Join(unzipPath, fileName)
		if !fileInfo.IsDir() {
			//to be compatible with pkzip(4.5)
			fullPathDir := filepath.Dir(fullPath)
			mErr := os.MkdirAll(fullPathDir, 0755)
			if mErr != nil {
				err = fmt.Errorf("Mkdir error, %s", mErr)
				continue
			}

			logs.Debug("Creating file", fullPath)
			localFp, openErr := os.OpenFile(fullPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, zipFile.Mode())
			if openErr != nil {
				err = fmt.Errorf("Open local file error, %s", openErr)
				continue
			}
			defer localFp.Close()

			zipFp, openErr := zipFile.Open()
			if openErr != nil {
				err = fmt.Errorf("Read zip content error, %s", openErr)
				continue
			}
			defer zipFp.Close()

			_, wErr := io.Copy(localFp, zipFp)
			if wErr != nil {
				err = fmt.Errorf("Save zip content error, %s", wErr)
				continue
			}
		}
	}
	return
}
