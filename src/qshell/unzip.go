package qshell

import (
	"archive/zip"
	"errors"
	"fmt"
	"github.com/qiniu/iconv"
	"io"
	"os"
	"path/filepath"
	"unicode/utf8"
)

func Unzip(zipFilePath string, unzipPath string) (err error) {
	zipReader, zipErr := zip.OpenReader(zipFilePath)
	if zipErr != nil {
		err = errors.New(fmt.Sprintf("Open zip file error, %s", zipErr))
		return
	}
	defer zipReader.Close()

	zipFiles := zipReader.File
	cd, cErr := iconv.Open("utf-8", "gbk")
	if cErr != nil {
		err = errors.New(fmt.Sprintf("Create charset converter error, %s", cErr))
		return
	}
	defer cd.Close()

	//list dir
	for _, zipFile := range zipFiles {
		fileInfo := zipFile.FileHeader.FileInfo()
		fileName := zipFile.FileHeader.Name

		//check charset utf8 or gbk
		if !utf8.Valid([]byte(fileName)) {
			fileName = cd.ConvString(fileName)
		}

		fullPath := filepath.Join(unzipPath, fileName)

		if fileInfo.IsDir() {
			mErr := os.MkdirAll(fullPath, 0775)
			if mErr != nil {
				err = errors.New(fmt.Sprintf("Mkdir error, %s", mErr))
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
			fileName = cd.ConvString(fileName)
		}

		fullPath := filepath.Join(unzipPath, fileName)
		if !fileInfo.IsDir() {
			localFp, openErr := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY, 0666)
			if openErr != nil {
				err = errors.New(fmt.Sprintf("Open local file error, %s", openErr))
				continue
			}
			defer localFp.Close()

			zipFp, openErr := zipFile.Open()
			if openErr != nil {
				err = errors.New(fmt.Sprintf("Read zip content error, %s", openErr))
				continue
			}
			defer zipFp.Close()

			_, wErr := io.Copy(localFp, zipFp)
			if wErr != nil {
				err = errors.New(fmt.Sprintf("Save zip content error, %s", wErr))
				continue
			}
		}
	}
	return
}
