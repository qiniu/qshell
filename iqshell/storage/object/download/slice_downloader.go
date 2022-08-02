package download

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type slice struct {
	index     int64 // 切片下表
	FromBytes int64 // 切片开始位置
	ToBytes   int64 // 切片终止位置
}

// 切片下载
type sliceDownloader struct {
	concurrentCount        int
	totalSliceCount        int64
	sliceSize              int64
	slicesDir              string
	slices                 chan slice
	downloadError          *data.CodeError
	currentReadSliceIndex  int64
	currentReadSliceOffset int64
}

func (s *sliceDownloader) Download(info *ApiInfo) (response *http.Response, err *data.CodeError) {
	err = s.initDownloadStatus(info)
	if err != nil {
		return
	}

	return s.download(info)
}

// 初始化状态
// 加载本地下载配置文件，没有则创建
func (s *sliceDownloader) initDownloadStatus(info *ApiInfo) *data.CodeError {
	if s.sliceSize <= 0 {
		s.sliceSize = 4 * utils.MB
	}
	if s.concurrentCount <= 0 {
		s.concurrentCount = 10
	}

	s.totalSliceCount = 0
	s.downloadError = nil
	s.currentReadSliceIndex = 0
	s.currentReadSliceOffset = 0
	if info.FromBytes > 0 {
		s.currentReadSliceIndex = info.FromBytes / s.sliceSize
		s.currentReadSliceOffset = info.FromBytes - s.currentReadSliceIndex*s.sliceSize
	}

	s.slices = make(chan slice, s.concurrentCount)
	// 临时文件夹
	s.slicesDir = filepath.Join(info.ToFile + "_download.slice")
	return utils.CreateDirIfNotExist(s.slicesDir)
}

// 并发下载
func (s *sliceDownloader) download(info *ApiInfo) (response *http.Response, err *data.CodeError) {

	go func() {
		var index int64 = 0
		var from int64 = 0
		var to int64 = 0
		for ; ; index++ {
			from = index * s.sliceSize
			to = from + s.sliceSize
			if from >= info.ServerFileSize {
				break
			}
			if to > info.ServerFileSize {
				to = info.ServerFileSize
			}
			s.slices <- slice{
				index:     index,
				FromBytes: from,
				ToBytes:   to,
			}
		}
		s.totalSliceCount = index
		close(s.slices)
	}()

	// 先尝试下载一个分片
	err = s.downloadSlice(info, <-s.slices)
	if err != nil {
		s.downloadError = err
		return nil, err
	}

	for i := 0; i < s.concurrentCount; i++ {
		go func() {
			for sl := range s.slices {
				if s.downloadError != nil {
					break
				}

				if e := s.downloadSlice(info, sl); e != nil {
					s.downloadError = e
					break
				}
			}
		}()
	}

	responseBodyContentLength := info.ServerFileSize - info.FromBytes
	responseHeader := http.Header{}
	responseHeader.Add("Content-Length", fmt.Sprintf("%d", responseBodyContentLength))
	return &http.Response{
		Status:        "slice download: 200",
		StatusCode:    200,
		Header:        responseHeader,
		Body:          s,
		ContentLength: responseBodyContentLength,
	}, nil
}

func (s *sliceDownloader) downloadSlice(info *ApiInfo, sl slice) *data.CodeError {
	toFile := filepath.Join(s.slicesDir, fmt.Sprintf("%d", sl.index))
	f, err := createDownloadFiles(toFile, info.FileEncoding)
	if err != nil {
		return err
	}

	file, _ := os.Stat(toFile)
	if file != nil && file.Size() == s.sliceSize {
		// 已下载
		return nil
	}

	f.fromBytes = sl.FromBytes + f.fromBytes

	log.DebugF("download slice, index:%d fromBytes:%d toBytes:%d", sl.index, sl.FromBytes, sl.ToBytes)
	return download(f, &ApiInfo{
		Bucket:               info.Bucket,
		Key:                  info.Key,
		IsPublic:             info.IsPublic,
		HostProvider:         info.HostProvider,
		DestDir:              info.DestDir,
		ToFile:               toFile,
		Referer:              info.Referer,
		FileEncoding:         info.FileEncoding,
		ServerFilePutTime:    0,
		ServerFileSize:       s.sliceSize,
		ServerFileHash:       "",
		FromBytes:            f.fromBytes,
		ToBytes:              sl.ToBytes,
		RemoveTempWhileError: false,
		UseGetFileApi:        false,
		Progress:             nil,
	})
}

func (s *sliceDownloader) Read(p []byte) (n int, err error) {
	if s.downloadError != nil {
		return 0, s.downloadError
	}

	if s.totalSliceCount > 0 && s.currentReadSliceIndex > s.totalSliceCount {
		return 0, io.EOF
	}

	currentReadSlicePath := filepath.Join(s.slicesDir, fmt.Sprintf("%d", s.currentReadSliceIndex))
	for {
		exist, _ := utils.ExistFile(currentReadSlicePath)
		if exist {
			break
		}
		time.Sleep(time.Millisecond * 500)
	}

	file, err := os.Open(currentReadSlicePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	n, err = file.ReadAt(p, s.currentReadSliceOffset)
	if err != nil {
		return n, err
	}

	s.currentReadSliceOffset += int64(n)

	if s.currentReadSliceOffset >= s.sliceSize {
		s.currentReadSliceOffset = 0
		s.currentReadSliceIndex += 1

		if e := os.Remove(currentReadSlicePath); e != nil {
			log.ErrorF("slice download delete slice error:%v", e)
		}
	}
	return
}

func (s *sliceDownloader) Close() error {
	if s.downloadError == nil {
		if err := os.RemoveAll(s.slicesDir); err != nil {
			log.ErrorF("slice download delete slice dir:%s error:%v", s.slicesDir, err)
		}
	}
	return nil
}
