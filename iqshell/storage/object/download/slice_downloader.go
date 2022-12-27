package download

import (
	"errors"
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type slice struct {
	index     int64 // 切片下表
	FromBytes int64 // 切片开始位置
	ToBytes   int64 // 切片终止位置
}

// 切片下载
type sliceDownloader struct {
	SliceSize              int64  `json:"slice_size"`
	FileHash               string `json:"file_hash"`
	UseGetFileApi          bool   `json:"use_get_file_api"`
	slicesDir              string
	concurrentCount        int
	totalSliceCount        int64
	slices                 chan slice
	downloadError          *data.CodeError
	currentReadSliceIndex  int64
	currentReadSliceOffset int64
	locker                 sync.Mutex
}

func (s *sliceDownloader) getDownloadError() *data.CodeError {
	s.locker.Lock()
	defer s.locker.Unlock()
	return s.downloadError
}

func (s *sliceDownloader) setDownloadError(err *data.CodeError) {
	s.locker.Lock()
	defer s.locker.Unlock()
	s.downloadError = err
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
	s.slices = make(chan slice, s.concurrentCount)
	s.FileHash = info.ServerFileHash
	toFile, err := filepath.Abs(info.ToFile)
	if err != nil {
		return data.NewEmptyError().AppendDescF("slice download, get abs file path:%s error:%v", info.ToFile, err)
	}

	// 临时文件夹
	s.slicesDir = filepath.Join(toFile + ".tmp.slices")

	if s.SliceSize <= 0 {
		s.SliceSize = 4 * utils.MB
	} else if s.SliceSize < 512*utils.KB {
		// 切片大小最小 512KB
		s.SliceSize = 512 * utils.KB
	}

	if s.concurrentCount <= 0 {
		s.concurrentCount = 10
	}

	// 配置文件
	configPath := filepath.Join(s.slicesDir, "config.json")
	oldConfig := &sliceDownloader{}
	// 读配置文件，不管存不存在都要读取，读取失败按照不存在处理，避免存在但因读取失败导致的后续问题
	if e := utils.UnMarshalFromFile(configPath, oldConfig); e != nil {
		log.WarningF("slice download UnMarshal config file error:%v", e)
	}
	// 分片大小不同会导致下载逻辑出错
	if oldConfig.SliceSize != s.SliceSize || oldConfig.FileHash != s.FileHash || oldConfig.UseGetFileApi != s.UseGetFileApi {
		// 不同则删除原来已下载但为合并的文件
		if e := os.RemoveAll(s.slicesDir); e != nil {
			log.WarningF("slice download remove all in dir:%s error:%v", s.slicesDir, e)
		} else {
			log.DebugF("slice download remove all in dir:%s", s.slicesDir)
		}
	}

	// 配置文件保存
	if e := utils.MarshalToFile(configPath, s); e != nil {
		log.WarningF("slice download marshal config file error:%v", e)
	}

	s.totalSliceCount = 0
	s.downloadError = nil
	s.currentReadSliceIndex = 0
	s.currentReadSliceOffset = 0
	if info.FromBytes > 0 {
		s.currentReadSliceIndex = info.FromBytes / s.SliceSize
		s.currentReadSliceOffset = info.FromBytes - s.currentReadSliceIndex*s.SliceSize
	}
	return nil
}

// 并发下载
func (s *sliceDownloader) download(info *ApiInfo) (response *http.Response, err *data.CodeError) {

	from := s.currentReadSliceIndex * s.SliceSize
	index := s.currentReadSliceIndex
	go func() {
		var to int64 = 0
		for ; ; index++ {
			from = index * s.SliceSize
			to = from + s.SliceSize - 1
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
				if workspace.IsCmdInterrupt() {
					s.setDownloadError(data.CancelError)
					break
				}

				if s.getDownloadError() != nil {
					break
				}

				if e := s.downloadSliceWithRetry(info, sl); e != nil {
					s.setDownloadError(e)
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

func (s *sliceDownloader) downloadSliceWithRetry(info *ApiInfo, sl slice) *data.CodeError {
	var downloadErr *data.CodeError = nil
	for i := 0; i < 3; i++ {
		downloadErr = s.downloadSlice(info, sl)
		if downloadErr == nil {
			break
		}
	}
	return downloadErr
}

func (s *sliceDownloader) downloadSlice(info *ApiInfo, sl slice) *data.CodeError {
	toFile := filepath.Join(s.slicesDir, fmt.Sprintf("%d", sl.index))
	f, err := createDownloadFiles(toFile, info.FileEncoding)
	if err != nil {
		return err
	}

	file, _ := os.Stat(toFile)
	if file != nil {
		if file.Size() == s.SliceSize {
			// 已下载
			return nil
		} else {
			if e := os.RemoveAll(toFile); e != nil {
				log.WarningF("delete slice:%s error:%v", toFile, e)
			}
		}
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
		ServerFileSize:       s.SliceSize,
		ServerFileHash:       s.FileHash,
		CheckHash:            false,
		FromBytes:            f.fromBytes,
		ToBytes:              sl.ToBytes,
		RemoveTempWhileError: false,
		UseGetFileApi:        info.UseGetFileApi,
		Progress:             nil,
	})
}

func (s *sliceDownloader) Read(p []byte) (int, error) {

	if s.getDownloadError() != nil {
		return 0, s.downloadError
	}

	if s.totalSliceCount > 0 && s.currentReadSliceIndex >= s.totalSliceCount {
		return 0, io.EOF
	}

	currentReadSlicePath := filepath.Join(s.slicesDir, fmt.Sprintf("%d", s.currentReadSliceIndex))
	for {
		if s.getDownloadError() != nil {
			return 0, s.downloadError
		}

		exist, _ := utils.ExistFile(currentReadSlicePath)
		if exist {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}

	file, err := os.Open(currentReadSlicePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	num, err := file.ReadAt(p, s.currentReadSliceOffset)
	if err != nil && !errors.Is(err, io.EOF) {
		return num, err
	}

	if err != nil && errors.Is(err, io.EOF) {
		s.currentReadSliceOffset = 0
		s.currentReadSliceIndex += 1

		if e := os.Remove(currentReadSlicePath); e != nil {
			log.ErrorF("slice download delete slice error:%v", e)
		}
	} else {
		s.currentReadSliceOffset += int64(num)
	}

	return num, nil
}

func (s *sliceDownloader) Close() error {
	if s.downloadError == nil {
		if err := os.RemoveAll(s.slicesDir); err != nil {
			log.ErrorF("slice download delete slice dir:%s error:%v", s.slicesDir, err)
		}
	}
	return nil
}
