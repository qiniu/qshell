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
	ToFile                 string `json:"-"`
	FileEncoding           string `json:"-"`
	ConcurrentCount        int    `json:"-"`
	slicesDir              string
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

func (s *sliceDownloader) Download(info *DownloadApiInfo) (response *http.Response, err *data.CodeError) {
	err = s.initDownloadStatus(info)
	if err != nil {
		return
	}

	return s.download(info)
}

// 初始化状态
// 加载本地下载配置文件，没有则创建
func (s *sliceDownloader) initDownloadStatus(info *DownloadApiInfo) *data.CodeError {
	s.slices = make(chan slice, s.ConcurrentCount)
	s.FileHash = info.FileHash

	// 临时文件夹
	s.slicesDir = filepath.Join(s.ToFile + ".tmp.slices")

	if s.SliceSize <= 0 {
		s.SliceSize = 4 * utils.MB
	} else if s.SliceSize < 512*utils.KB {
		// 切片大小最小 512KB
		s.SliceSize = 512 * utils.KB
	}

	if s.ConcurrentCount <= 0 {
		s.ConcurrentCount = 10
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
	if info.RangeFromBytes > 0 {
		s.currentReadSliceIndex = info.RangeFromBytes / s.SliceSize
		s.currentReadSliceOffset = info.RangeFromBytes - s.currentReadSliceIndex*s.SliceSize
	}
	return nil
}

// 并发下载
func (s *sliceDownloader) download(info *DownloadApiInfo) (response *http.Response, err *data.CodeError) {

	from := s.currentReadSliceIndex * s.SliceSize
	index := s.currentReadSliceIndex
	go func() {
		var to int64 = 0
		for ; ; index++ {
			from = index * s.SliceSize
			to = from + s.SliceSize - 1
			if from >= info.FileSize {
				break
			}
			if to >= info.FileSize && info.FileSize > 0 {
				to = info.FileSize - 1
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

	for i := 0; i < s.ConcurrentCount; i++ {
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

	responseBodyContentLength := info.FileSize - info.RangeFromBytes
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

func (s *sliceDownloader) downloadSliceWithRetry(info *DownloadApiInfo, sl slice) *data.CodeError {
	var downloadErr *data.CodeError = nil
	for i := 0; i < 3; i++ {
		downloadErr = s.downloadSlice(info, sl)
		if downloadErr == nil {
			break
		}
	}
	return downloadErr
}

func (s *sliceDownloader) downloadSlice(info *DownloadApiInfo, sl slice) *data.CodeError {
	toFile := filepath.Join(s.slicesDir, fmt.Sprintf("%d", sl.index))
	f, err := createDownloadFiles(toFile, s.FileEncoding)
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
	err = downloadTempFileWithDownloader(&downloaderFile{}, f, &DownloadApiInfo{
		Url:            info.Url,
		Host:           info.Host,
		Referer:        info.Referer,
		RangeFromBytes: f.fromBytes,
		RangeToBytes:   sl.ToBytes,
		CheckSize:      false,
		FileSize:       info.FileSize,
		CheckHash:      false,
		FileHash:       info.FileHash,
		Progress:       nil,
	})
	if err != nil {
		return err
	}

	return renameTempFile(f)
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
