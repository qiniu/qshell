package qshell

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func Unzip(zipFilePath string, unzipPath string) (err error) {
	zipReader, zipErr := zip.OpenReader(zipFilePath)
	if zipErr != nil {
		err = zipErr
		return
	}
	defer zipReader.Close()

	zipFiles := zipReader.File
	for _, zipFile := range zipFiles {
		fileInfo := zipFile.FileHeader.FileInfo()
		fileName := zipFile.FileHeader.Name
		fullPath := filepath.Join(unzipPath, fileName)
		if fileInfo.IsDir() {
			mErr := os.MkdirAll(fullPath, 0775)
			if mErr != nil {
				err = mErr
				continue
			}
		} else {
			localFp, openErr := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY, 0666)
			if openErr != nil {
				err = openErr
				continue
			}
			defer localFp.Close()

			zipFp, openErr := zipFile.Open()
			if openErr != nil {
				err = openErr
				continue
			}
			defer zipFp.Close()

			_, wErr := io.Copy(localFp, zipFp)
			if wErr != nil {
				err = wErr
				continue
			}
		}
	}
	return
}
