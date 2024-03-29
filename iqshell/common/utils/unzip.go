package utils

import (
	"archive/zip"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"io"
	"os"
	"path/filepath"
	"unicode/utf8"

	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"golang.org/x/text/encoding/simplifiedchinese"
)

func Gbk2Utf8(text string) (string, *data.CodeError) {
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

func Unzip(zipFilePath string, unzipPath string) (err *data.CodeError) {
	zipReader, zipErr := zip.OpenReader(zipFilePath)
	if zipErr != nil {
		err = data.NewEmptyError().AppendDescF("Open zip file error, %s", zipErr)
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
			fileName, err = Gbk2Utf8(fileName)
			if err != nil {
				err = data.NewEmptyError().AppendDesc("Unsupported filename encoding")
				continue
			}
		}

		fullPath := filepath.Join(unzipPath, fileName)
		if fileInfo.IsDir() {
			log.Debug("Mkdir", fullPath)
			mErr := os.MkdirAll(fullPath, 0775)
			if mErr != nil {
				err = data.NewEmptyError().AppendDescF("Mkdir error, %s", mErr)
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
			fileName, err = Gbk2Utf8(fileName)
			if err != nil {
				err = data.NewEmptyError().AppendDesc("Unsupported filename encoding")
				continue
			}
		}

		fullPath := filepath.Join(unzipPath, fileName)
		if !fileInfo.IsDir() {
			//to be compatible with pkzip(4.5)
			fullPathDir := filepath.Dir(fullPath)
			mErr := os.MkdirAll(fullPathDir, 0755)
			if mErr != nil {
				err = data.NewEmptyError().AppendDescF("Mkdir error, %v", mErr)
				continue
			}

			log.Debug("Creating file", fullPath)
			localFp, openErr := os.OpenFile(fullPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, zipFile.Mode())
			if openErr != nil {
				err = data.NewEmptyError().AppendDescF("Open local file error, %v", openErr)
				continue
			}
			defer localFp.Close()

			zipFp, openErr := zipFile.Open()
			if openErr != nil {
				err = data.NewEmptyError().AppendDescF("Read zip content error, %v", openErr)
				continue
			}
			defer zipFp.Close()

			_, wErr := io.Copy(localFp, zipFp)
			if wErr != nil {
				err = data.NewEmptyError().AppendDescF("Save zip content error, %v", wErr)
				continue
			}
		}
	}
	return
}
