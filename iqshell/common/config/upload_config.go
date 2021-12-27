package config

import (
	"encoding/json"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"os"
	"path/filepath"
	"strings"
)

type UploadConfig struct {
	FileEncoding string `json:"file_encoding"`
	//basic config
	SrcDir string `json:"src_dir"`
	Bucket string `json:"bucket"`

	//optional config
	ResumableAPIV2         bool   `json:"resumable_api_v2,omitempty"`
	ResumableAPIV2PartSize int64  `json:"resumable_api_v2_part_size,omitempty"`
	FileList               string `json:"file_list,omitempty"`
	PutThreshold           int64  `json:"put_threshold,omitempty"`
	KeyPrefix              string `json:"key_prefix,omitempty"`
	IgnoreDir              bool   `json:"ignore_dir,omitempty"`
	Overwrite              bool   `json:"overwrite,omitempty"`
	CheckExists            bool   `json:"check_exists,omitempty"`
	CheckHash              bool   `json:"check_hash,omitempty"`
	CheckSize              bool   `json:"check_size,omitempty"`
	SkipFilePrefixes       string `json:"skip_file_prefixes,omitempty"`
	SkipPathPrefixes       string `json:"skip_path_prefixes,omitempty"`
	SkipFixedStrings       string `json:"skip_fixed_strings,omitempty"`
	SkipSuffixes           string `json:"skip_suffixes,omitempty"`
	RescanLocal            bool   `json:"rescan_local,omitempty"`
	FileType               int    `json:"file_type,omitempty"`

	//advanced config
	UpHost string `json:"up_host,omitempty"`

	BindUpIp string `json:"bind_up_ip,omitempty"`
	BindRsIp string `json:"bind_rs_ip,omitempty"`
	//local network interface card config
	BindNicIp string `json:"bind_nic_ip,omitempty"`

	//log settings
	LogLevel  string `json:"log_level,omitempty"`
	LogFile   string `json:"log_file,omitempty"`
	LogRotate int    `json:"log_rotate,omitempty"`
	LogStdout bool   `json:"log_stdout,omitempty"`

	//more settings
	DeleteOnSuccess bool   `json:"delete_on_success,omitempty"`
	DisableResume   bool   `json:"disable_resume,omitempty"`
	CallbackUrls    string `json:"callback_urls,omitempty"`
	CallbackHost    string `json:"callback_host,omitempty"`
	PutPolicy       storage.PutPolicy
}

func (cfg *UploadConfig) Check() {
	// 验证大小
	if cfg.ResumableAPIV2PartSize <= 0 {
		cfg.ResumableAPIV2PartSize = data.BLOCK_SIZE
	} else if cfg.ResumableAPIV2PartSize < int64(utils.MB) {
		cfg.ResumableAPIV2PartSize = int64(utils.MB)
	} else if cfg.ResumableAPIV2PartSize > int64(utils.GB) {
		cfg.ResumableAPIV2PartSize = int64(utils.GB)
	}
}

func (cfg *UploadConfig) JobId() string {

	return utils.Md5Hex(fmt.Sprintf("%s:%s", cfg.SrcDir, cfg.Bucket))
}

func (cfg *UploadConfig) GetLogLevel() int {

	//init log level
	logLevel := log.LevelInfo
	switch cfg.LogLevel {
	case "debug":
		logLevel = log.LevelDebug
	case "info":
		logLevel = log.LevelInfo
	case "warn":
		logLevel = log.LevelWarning
	case "error":
		logLevel = log.LevelError
	default:
		logLevel = log.LevelInfo
	}
	return int(logLevel)
}

func (cfg *UploadConfig) GetLogRotate() int {
	logRotate := 1
	if cfg.LogRotate > 0 {
		logRotate = cfg.LogRotate
	}
	return logRotate
}

func (cfg *UploadConfig) HitByPathPrefixes(localFileRelativePath string) (hit bool, pathPrefix string) {

	if cfg.SkipPathPrefixes != "" {
		//unpack skip prefix
		pathPrefixes := strings.Split(cfg.SkipPathPrefixes, ",")
		for _, prefix := range pathPrefixes {
			if strings.TrimSpace(prefix) == "" {
				continue
			}

			if strings.HasPrefix(localFileRelativePath, prefix) {
				pathPrefix = prefix
				hit = true
				break
			}
		}
	}
	return
}

func (cfg *UploadConfig) HitByFilePrefixes(localFileRelativePath string) (hit bool, filePrefix string) {
	if cfg.SkipFilePrefixes != "" {
		//unpack skip prefix
		filePrefixes := strings.Split(cfg.SkipFilePrefixes, ",")
		for _, prefix := range filePrefixes {
			if strings.TrimSpace(prefix) == "" {
				continue
			}

			localFileName := filepath.Base(localFileRelativePath)
			if strings.HasPrefix(localFileName, prefix) {
				filePrefix = prefix
				hit = true
				break
			}
		}
	}
	return
}

func (cfg *UploadConfig) HitByFixesString(localFileRelativePath string) (hit bool, hitFixedStr string) {
	if cfg.SkipFixedStrings != "" {
		//unpack fixed strings
		fixedStrings := strings.Split(cfg.SkipFixedStrings, ",")
		for _, fixedStr := range fixedStrings {
			if strings.TrimSpace(fixedStr) == "" {
				continue
			}

			if strings.Contains(localFileRelativePath, fixedStr) {
				hitFixedStr = fixedStr
				hit = true
				break
			}
		}
	}
	return

}

func (cfg *UploadConfig) HitBySuffixes(localFileRelativePath string) (hit bool, hitSuffix string) {
	if cfg.SkipSuffixes != "" {
		suffixes := strings.Split(cfg.SkipSuffixes, ",")
		for _, suffix := range suffixes {
			if strings.TrimSpace(suffix) == "" {
				continue
			}

			if strings.HasSuffix(localFileRelativePath, suffix) {
				hitSuffix = suffix
				hit = true
				break
			}
		}
	}
	return
}

func (cfg *UploadConfig) CacheFileNameAndCount(storePath, jobId string) (cacheResultName string, totalFileCount int64, cacheErr error) {
	//find the local file list, by specified or by config

	_, localFileStatErr := os.Stat(cfg.FileList)
	if cfg.FileList != "" && localFileStatErr == nil {
		//use specified file list
		cacheResultName = cfg.FileList
		totalFileCount = utils.GetFileLineCount(cacheResultName)
	} else {
		cacheResultName = filepath.Join(storePath, fmt.Sprintf("%s.cache", jobId))
		cacheCountName := filepath.Join(storePath, fmt.Sprintf("%s.count", jobId))
		totalFileCount, cacheErr = prepareCacheFileList(cacheResultName, cacheCountName, cfg.SrcDir, cfg.RescanLocal)
		if cacheErr != nil {
			return
		}
	}
	return

}

type UploadInfo struct {
	TotalFileCount int64 `json:"total_file_count"`
}
func prepareCacheFileList(cacheResultName, cacheCountName, srcDir string, rescanLocal bool) (totalFileCount int64, cacheErr error) {
	//cache file
	cacheTempName := fmt.Sprintf("%s.temp", cacheResultName)

	rescanLocalDir := false
	if _, statErr := os.Stat(cacheResultName); statErr == nil {
		//file exists
		rescanLocalDir = rescanLocal
	} else {
		rescanLocalDir = true
	}

	if rescanLocalDir {
		log.Info("Listing local sync dir, this can take a long time for big directory, please wait patiently")
		totalFileCount, cacheErr = utils.DirCache(srcDir, cacheTempName)
		if cacheErr != nil {
			return
		}

		if rErr := os.Rename(cacheTempName, cacheResultName); rErr != nil {
			log.Error("Rename the temp cached file error", rErr)
			cacheErr = rErr
			return
		}
		//write the total count to local file
		if cFp, cErr := os.Create(cacheCountName); cErr == nil {
			func() {
				defer cFp.Close()
				uploadInfo := UploadInfo{
					TotalFileCount: totalFileCount,
				}
				uploadInfoBytes, mErr := json.Marshal(&uploadInfo)
				if mErr == nil {
					if _, wErr := cFp.Write(uploadInfoBytes); wErr != nil {
						log.WarningF("Write local cached count file error %s", cErr)
					} else {
						cFp.Close()
					}
				}
			}()
		} else {
			log.ErrorF("Open local cached count file error %s,", cErr)
		}
	} else {
		log.Info("Use the last cached local sync dir file list")
		//read from local cache
		if rFp, rErr := os.Open(cacheCountName); rErr == nil {
			func() {
				defer rFp.Close()
				uploadInfo := UploadInfo{}
				decoder := json.NewDecoder(rFp)
				if dErr := decoder.Decode(&uploadInfo); dErr == nil {
					totalFileCount = uploadInfo.TotalFileCount
				}
			}()
		} else {
			log.WarningF("Open local cached count file error %s,", rErr)
			totalFileCount = utils.GetFileLineCount(cacheResultName)
		}
	}

	return
}

func (cfg *UploadConfig) DefaultLogFile(storePath, jobId string) (defaultLogFile string, err error) {
	//local storage path
	if mkdirErr := os.MkdirAll(storePath, 0775); mkdirErr != nil {
		log.ErrorF("Failed to mkdir `%s` due to `%s`", storePath, mkdirErr)
		err = mkdirErr
		return
	}
	defaultLogFile = filepath.Join(storePath, fmt.Sprintf("%s.log", jobId))
	return
}

func (cfg *UploadConfig) PrepareLogger(storePath, jobId string) {

	defaultLogFile, err := cfg.DefaultLogFile(storePath, jobId)
	if err != nil {
		os.Exit(data.STATUS_HALT)
	}
	logLevel := cfg.GetLogLevel()
	logRotate := cfg.GetLogRotate()

	//init log writer
	if cfg.LogFile == "" {
		//set default log file
		cfg.LogFile = defaultLogFile
	}

	//Todo: 处理 stdout 不输出
	if !cfg.LogStdout {
		// log.GetBeeLogger().DelLogger(logs.AdapterConsole)
	}
	//open log file
	fmt.Println("Writing upload log to file", cfg.LogFile)

	//daily rotate
	logCfg := log.Config{
		Filename: cfg.LogFile,
		Level:    logLevel,
		Daily:    true,
		MaxDays:  logRotate,
	}
	log.LoadFileLogger(logCfg)
	fmt.Println()
}

func (cfg *UploadConfig) UploadToken(mac *qbox.Mac, uploadFileKey string) string {

	policy := cfg.PutPolicy
	policy.Scope = cfg.Bucket
	if cfg.Overwrite {
		policy.Scope = fmt.Sprintf("%s:%s", cfg.Bucket, uploadFileKey)
		policy.InsertOnly = 0
	}

	policy.FileType = cfg.FileType
	policy.Expires = 7 * 24 * 3600
	return policy.UploadToken(mac)
}
