package qshell

import (
	"os"
	"testing"
)

func TestGetFileTotalCount(t *testing.T) {
	fpath := "/Users/jemy/Temp/t.go"
	totalCount := getFileLineCount(fpath)
	t.Log(totalCount)
}

func TestSimpleUpload(t *testing.T) {
	uploadConfig := UploadConfig{
		SrcDir:    "/Users/jemy/Temp/test",
		AccessKey: os.Getenv("AccessKey"),
		SecretKey: os.Getenv("SecretKey"),
		Bucket:    "if-pbl",
	}

	QiniuUpload(1, &uploadConfig)
}

func TestOverwriteUpload(t *testing.T) {
	uploadConfig := UploadConfig{
		SrcDir:      "/Users/jemy/Temp/test",
		AccessKey:   os.Getenv("AccessKey"),
		SecretKey:   os.Getenv("SecretKey"),
		Bucket:      "if-pbl",
		Overwrite:   true,
		RescanLocal: true,
	}

	QiniuUpload(1, &uploadConfig)
}

//use when files are delete from the buckets
func TestCheckExistsUpload(t *testing.T) {
	uploadConfig := UploadConfig{
		SrcDir:      "/Users/jemy/Temp/test",
		AccessKey:   os.Getenv("AccessKey"),
		SecretKey:   os.Getenv("SecretKey"),
		Bucket:      "if-pbl",
		Overwrite:   true,
		RescanLocal: true,
		CheckExists: true,
	}

	QiniuUpload(1, &uploadConfig)
}

func TestUploadWithFileList(t *testing.T) {
	flist := "/Users/jemy/Temp/test/flist.txt"
	uploadConfig := UploadConfig{
		SrcDir:    "/Users/jemy/Temp/test",
		AccessKey: os.Getenv("AccessKey"),
		SecretKey: os.Getenv("SecretKey"),
		Bucket:    "if-pbl",
		FileList:  flist,
	}

	QiniuUpload(1, &uploadConfig)
}
