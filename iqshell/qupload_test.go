package iqshell

import (
	"os"
	"testing"
)

func TestGetFileTotalCount(t *testing.T) {
	fpath := "/Users/jemy/Worklab/qshelltest/qiniu.txt"
	totalCount := GetFileLineCount(fpath)
	if totalCount != 5 {
		t.Fail()
	}
}

func TestSimpleUpload(t *testing.T) {
	uploadConfig := UploadConfig{
		SrcDir: os.Getenv("SrcDir"),
		Bucket: os.Getenv("Bucket"),
	}

	QiniuUpload(1, &uploadConfig)
}

func TestSimpleUploadWithKeyPrefix(t *testing.T) {
	uploadConfig := UploadConfig{
		SrcDir:    os.Getenv("SrcDir"),
		Bucket:    os.Getenv("Bucket"),
		KeyPrefix: os.Getenv("Prefix"),
	}

	QiniuUpload(1, &uploadConfig)
}

func TestSimpleUploadIgnoreDir(t *testing.T) {
	uploadConfig := UploadConfig{
		SrcDir:    os.Getenv("SrcDir"),
		Bucket:    os.Getenv("Bucket"),
		KeyPrefix: os.Getenv("Prefix"),
		IgnoreDir: true,
	}

	QiniuUpload(1, &uploadConfig)
}

func TestOverwriteUpload(t *testing.T) {
	uploadConfig := UploadConfig{
		SrcDir:      os.Getenv("SrcDir"),
		Bucket:      os.Getenv("Bucket"),
		Overwrite:   true,
		RescanLocal: true,
	}

	QiniuUpload(1, &uploadConfig)
}

//use when files are delete from the buckets
func TestCheckExistsUpload(t *testing.T) {
	uploadConfig := UploadConfig{
		SrcDir:      os.Getenv("SrcDir"),
		Bucket:      os.Getenv("Bucket"),
		Overwrite:   true,
		RescanLocal: true,
		CheckExists: true,
	}

	QiniuUpload(1, &uploadConfig)
}

func TestUploadWithFileList(t *testing.T) {
	flist := os.Getenv("SrcFileList")
	uploadConfig := UploadConfig{
		SrcDir:   os.Getenv("SrcDir"),
		Bucket:   os.Getenv("Bucket"),
		FileList: flist,
	}

	QiniuUpload(1, &uploadConfig)
}
