package cmd

import (
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/docs"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/config"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/upload/operations"
	"github.com/spf13/cobra"
)

var uploadCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var (
		callbackUrl  = ""
		callbackHost = ""
	)
	info := operations.BatchUploadInfo{}
	cmd := &cobra.Command{
		Use:   "qupload <LocalDownloadConfig>",
		Short: "Batch upload files to the qiniu bucket",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.QUploadType
			if len(args) > 0 {
				cfg.UploadConfigFile = args[0]
			}
			cfg.CmdCfg.Up = &config.Up{
				Policy:                 &storage.PutPolicy{
					Scope:               "",
					Expires:             0,
					IsPrefixalScope:     0,
					InsertOnly:          0,
					DetectMime:          0,
					FsizeMin:            0,
					FsizeLimit:          0,
					MimeLimit:           "",
					ForceSaveKey:        false,
					SaveKey:             "",
					CallbackFetchKey:    0,
					CallbackURL:         callbackUrl,
					CallbackHost:        callbackHost,
					CallbackBody:        "",
					CallbackBodyType:    "",
					ReturnURL:           "",
					ReturnBody:          "",
					PersistentOps:       "",
					PersistentNotifyURL: "",
					PersistentPipeline:  "",
					EndUser:             "",
					DeleteAfterDays:     0,
					FileType:            0,
				},
			}
			info.GroupInfo.Force = true
			operations.BatchUpload(cfg, info)
		},
	}
	cmd.Flags().StringVarP(&info.GroupInfo.SuccessExportFilePath, "success-list", "s", "", "upload success (all) file list")
	cmd.Flags().StringVarP(&info.GroupInfo.FailExportFilePath, "failure-list", "f", "", "upload failure file list")
	cmd.Flags().StringVarP(&info.GroupInfo.OverrideExportFilePath, "overwrite-list", "w", "", "upload success (overwrite) file list")
	cmd.Flags().IntVarP(&info.GroupInfo.WorkCount, "worker", "c", 1, "worker count")
	cmd.Flags().StringVarP(&callbackUrl, "callback-urls", "l", "", "upload callback urls, separated by comma")
	cmd.Flags().StringVarP(&callbackHost, "callback-host", "T", "", "upload callback host")
	return cmd
}

var upload2CmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var (
		ResumableAPIV2 = false
		IgnoreDir      = false
		Overwrite      = false
		CheckExists    = false
		CheckHash      = false
		CheckSize      = false
		RescanLocal    = false

		ResumableAPIV2PartSize int64 = data.BLOCK_SIZE
		SrcDir                       = ""
		FileList                     = ""
		Bucket                       = ""
		PutThreshold           int64 = 0
		KeyPrefix                    = ""
		SkipFilePrefixes             = ""
		SkipPathPrefixes             = ""
		SkipFixedStrings             = ""
		SkipSuffixes                 = ""
		UpHost                       = ""
		RecordRoot                   = ""
		LogFile                      = ""
		LogLevel                     = ""
		LogRotate                    = 7
		FileType                     = 0
		callbackUrl                  = ""
		callbackHost                 = ""
	)
	info := operations.BatchUploadInfo{}
	cmd := &cobra.Command{
		Use:   "qupload2",
		Short: "Batch upload files to the qiniu bucket",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.QUpload2Type
			cfg.CmdCfg.Up = &config.Up{
				LogSetting:             &config.LogSetting{
					LogLevel:  data.NewString(LogLevel),
					LogFile:   data.NewString(LogFile),
					LogRotate: data.NewInt(LogRotate),
					LogStdout: nil,
				},
				UpHost:                 data.NewString(UpHost),
				BindUpIp:               nil,
				BindRsIp:               nil,
				BindNicIp:              nil,
				SrcDir:                 data.NewString(SrcDir),
				FileList:               data.NewString(FileList),
				IgnoreDir:              data.NewBool(IgnoreDir),
				SkipFilePrefixes:       data.NewString(SkipFilePrefixes),
				SkipPathPrefixes:       data.NewString(SkipPathPrefixes),
				SkipFixedStrings:       data.NewString(SkipFixedStrings),
				SkipSuffixes:           data.NewString(SkipSuffixes),
				FileEncoding:           nil,
				Bucket:                 data.NewString(Bucket),
				ResumableAPIV2:         data.NewBool(ResumableAPIV2),
				ResumableAPIV2PartSize: data.NewInt64(data.BLOCK_SIZE),
				PutThreshold:           data.NewInt64(PutThreshold),
				KeyPrefix:              data.NewString(KeyPrefix),
				Overwrite:              data.NewBool(Overwrite),
				CheckExists:            data.NewBool(CheckExists),
				CheckHash:              data.NewBool(CheckHash),
				CheckSize:              data.NewBool(CheckSize),
				RescanLocal:            data.NewBool(RescanLocal),
				FileType:               data.NewInt(FileType),
				DeleteOnSuccess:        nil,
				DisableResume:          nil,
				DisableForm:            nil,
				WorkerCount:            nil,
				RecordRoot:             data.NewString(RecordRoot),
				Tasks:                  nil,
				Retry:                  nil,
				Policy:                 &storage.PutPolicy{
					Scope:               "",
					Expires:             0,
					IsPrefixalScope:     0,
					InsertOnly:          0,
					DetectMime:          0,
					FsizeMin:            0,
					FsizeLimit:          0,
					MimeLimit:           "",
					ForceSaveKey:        false,
					SaveKey:             "",
					CallbackFetchKey:    0,
					CallbackURL:         callbackUrl,
					CallbackHost:        callbackHost,
					CallbackBody:        "",
					CallbackBodyType:    "",
					ReturnURL:           "",
					ReturnBody:          "",
					PersistentOps:       "",
					PersistentNotifyURL: "",
					PersistentPipeline:  "",
					EndUser:             "",
					DeleteAfterDays:     0,
					FileType:            0,
				},
			}
			info.GroupInfo.Force = true
			operations.BatchUpload2(cfg, info)
		},
	}
	cmd.Flags().StringVar(&info.GroupInfo.SuccessExportFilePath, "success-list", "", "upload success file list")
	cmd.Flags().StringVar(&info.GroupInfo.FailExportFilePath, "failure-list", "", "upload failure file list")
	cmd.Flags().StringVar(&info.GroupInfo.OverrideExportFilePath, "overwrite-list", "", "upload success (overwrite) file list")
	cmd.Flags().IntVar(&info.GroupInfo.WorkCount, "thread-count", 1, "multiple thread count")

	cmd.Flags().BoolVarP(&ResumableAPIV2, "resumable-api-v2", "", false, "use resumable upload v2 APIs to upload")
	cmd.Flags().BoolVar(&IgnoreDir, "ignore-dir", false, "ignore the dir in the dest file key")
	cmd.Flags().BoolVar(&Overwrite, "overwrite", false, "overwrite the file of same key in bucket")
	cmd.Flags().BoolVar(&CheckExists, "check-exists", false, "check file key whether in bucket before upload")
	cmd.Flags().BoolVar(&CheckHash, "check-hash", false, "check hash")
	cmd.Flags().BoolVar(&CheckSize, "check-size", false, "check file size")
	cmd.Flags().BoolVar(&RescanLocal, "rescan-local", false, "rescan local dir to upload newly add files")

	cmd.Flags().Int64Var(&ResumableAPIV2PartSize, "resumable-api-v2-part-size", data.BLOCK_SIZE, "the part size when use resumable upload v2 APIs to upload")
	cmd.Flags().StringVar(&SrcDir, "src-dir", "", "src dir to upload")
	cmd.Flags().StringVar(&FileList, "file-list", "", "file list to upload")
	cmd.Flags().StringVar(&Bucket, "bucket", "", "bucket")
	cmd.Flags().Int64Var(&PutThreshold, "put-threshold", 0, "chunk upload threshold")
	cmd.Flags().StringVar(&KeyPrefix, "key-prefix", "", "key prefix prepended to dest file key")
	cmd.Flags().StringVar(&SkipFilePrefixes, "skip-file-prefixes", "", "skip files with these file prefixes")
	cmd.Flags().StringVar(&SkipPathPrefixes, "skip-path-prefixes", "", "skip files with these relative path prefixes")
	cmd.Flags().StringVar(&SkipFixedStrings, "skip-fixed-strings", "", "skip files with the fixed string in the name")
	cmd.Flags().StringVar(&SkipSuffixes, "skip-suffixes", "", "skip files with these suffixes")
	cmd.Flags().StringVar(&UpHost, "up-host", "", "upload host")
	cmd.Flags().StringVar(&RecordRoot, "record-root", "", "record root dir, and will save record info to the dir(db and log), default <UserRoot>/.qshell")
	cmd.Flags().StringVar(&LogFile, "log-file", "", "log file")
	cmd.Flags().StringVar(&LogLevel, "log-level", "debug", "log level")
	cmd.Flags().IntVar(&LogRotate, "log-rotate", 7, "log rotate days")

	cmd.Flags().IntVar(&FileType, "file-type", 0, "set storage file type")
	cmd.Flags().StringVarP(&callbackUrl, "callback-urls", "l", "", "upload callback urls, separated by comma")
	cmd.Flags().StringVarP(&callbackHost, "callback-host", "T", "", "upload callback host")
	//cmd.Flags().StringVar(&cfg.CmdCfg.Up.BindUpIp, "bind-up-ip", "", "upload host ip to bind")
	//cmd.Flags().StringVar(&cfg.CmdCfg.Up.BindRsIp, "bind-rs-ip", "", "rs host ip to bind")
	//cmd.Flags().StringVar(&cfg.CmdCfg.Up.BindNicIp, "bind-nic-ip", "", "local network interface card to bind")
	return cmd
}

var syncCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var (
		upHost      = ""
		resumeAPIV2 = false
	)
	info := operations.SyncInfo{}
	cmd := &cobra.Command{
		Use:   "sync <SrcResUrl> <Buckets> [-k <Key>]",
		Short: "Sync big file to qiniu bucket",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.SyncType
			cfg.CmdCfg.Up = &config.Up{
				LogSetting:             nil,
				UpHost:                 data.NewString(upHost),
				FileEncoding:           nil,
				Bucket:                 nil,
				ResumableAPIV2:         data.NewBool(resumeAPIV2),
				ResumableAPIV2PartSize: nil,
				PutThreshold:           nil,
				Overwrite:              nil,
				CheckExists:            nil,
				CheckHash:              nil,
				CheckSize:              nil,
				FileType:               nil,
				DeleteOnSuccess:        nil,
				DisableResume:          data.NewBool(true),
				DisableForm:            nil,
				WorkerCount:            nil,
				Policy: &storage.PutPolicy{
					Scope:               "",
					Expires:             0,
					IsPrefixalScope:     0,
					InsertOnly:          0,
					DetectMime:          0,
					FsizeMin:            0,
					FsizeLimit:          0,
					MimeLimit:           "",
					ForceSaveKey:        false,
					SaveKey:             "",
					CallbackFetchKey:    0,
					CallbackURL:         "",
					CallbackHost:        "",
					CallbackBody:        "",
					CallbackBodyType:    "",
					ReturnURL:           "",
					ReturnBody:          "",
					PersistentOps:       "",
					PersistentNotifyURL: "",
					PersistentPipeline:  "",
					EndUser:             "",
					DeleteAfterDays:     0,
					FileType:            0,
				},
			}
			if len(args) > 0 {
				info.FilePath = args[0]
			}
			if len(args) > 1 {
				info.Bucket = args[1]
			}
			operations.SyncFile(cfg, info)
		},
	}
	cmd.Flags().StringVarP(&info.Key, "key", "k", "", "save as <key> in bucket")
	cmd.Flags().BoolVarP(&resumeAPIV2, "resumable-api-v2", "", false, "use resumable upload v2 APIs to upload")
	cmd.Flags().StringVarP(&upHost, "up-host", "u", "", "upload host")
	return cmd
}

var formUploadCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var (
		upHost       = ""
		fileType     = 0
		overwrite    = false
		callbackUrl  = ""
		callbackHost = ""
	)

	info := operations.UploadInfo{}
	cmd := &cobra.Command{
		Use:   "fput <Bucket> <Key> <LocalFile>",
		Short: "Form upload a local file",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.FormPutType
			cfg.CmdCfg.Up = &config.Up{
				LogSetting:             nil,
				UpHost:                 data.NewString(upHost),
				FileEncoding:           nil,
				Bucket:                 nil,
				ResumableAPIV2:         nil,
				ResumableAPIV2PartSize: nil,
				PutThreshold:           nil,
				Overwrite:              data.NewBool(overwrite),
				CheckExists:            nil,
				CheckHash:              nil,
				CheckSize:              nil,
				FileType:               data.NewInt(fileType),
				DeleteOnSuccess:        nil,
				DisableResume:          data.NewBool(true),
				DisableForm:            nil,
				WorkerCount:            nil,
				Policy: &storage.PutPolicy{
					Scope:               "",
					Expires:             0,
					IsPrefixalScope:     0,
					InsertOnly:          0,
					DetectMime:          0,
					FsizeMin:            0,
					FsizeLimit:          0,
					MimeLimit:           "",
					ForceSaveKey:        false,
					SaveKey:             "",
					CallbackFetchKey:    0,
					CallbackURL:         callbackUrl,
					CallbackHost:        callbackHost,
					CallbackBody:        "",
					CallbackBodyType:    "",
					ReturnURL:           "",
					ReturnBody:          "",
					PersistentOps:       "",
					PersistentNotifyURL: "",
					PersistentPipeline:  "",
					EndUser:             "",
					DeleteAfterDays:     0,
					FileType:            0,
				},
			}
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			if len(args) > 1 {
				info.Key = args[1]
			}
			if len(args) > 2 {
				info.FilePath = args[2]
			}
			operations.UploadFile(cfg, info)
		},
	}
	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite the file of same key in bucket")
	cmd.Flags().StringVarP(&info.MimeType, "mimetype", "t", "", "file mime type")
	cmd.Flags().IntVarP(&fileType, "storage", "s", 0, "storage type")
	cmd.Flags().StringVarP(&upHost, "up-host", "u", "", "uphost")
	cmd.Flags().StringVarP(&callbackUrl, "callback-urls", "l", "", "upload callback urls, separated by comma")
	cmd.Flags().StringVarP(&callbackHost, "callback-host", "T", "", "upload callback host")

	return cmd
}

var resumeUploadCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var (
		resumeAPIV2                  = false
		overwrite                    = false
		FileType                     = 0
		WorkerCount                  = 16
		UpHost                       = ""
		ResumableAPIV2PartSize int64 = data.BLOCK_SIZE
		callbackUrl                  = ""
		callbackHost                 = ""
	)
	info := operations.UploadInfo{}
	cmd := &cobra.Command{
		Use:   "rput <Bucket> <Key> <LocalFile>",
		Short: "Resumable upload a local file",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.RPutType
			cfg.CmdCfg.Up = &config.Up{
				LogSetting:             nil,
				UpHost:                 data.NewString(UpHost),
				FileEncoding:           nil,
				Bucket:                 nil,
				ResumableAPIV2:         data.NewBool(resumeAPIV2),
				ResumableAPIV2PartSize: data.NewInt64(ResumableAPIV2PartSize),
				PutThreshold:           nil,
				Overwrite:              data.NewBool(overwrite),
				CheckExists:            nil,
				CheckHash:              nil,
				CheckSize:              nil,
				FileType:               data.NewInt(FileType),
				DeleteOnSuccess:        nil,
				DisableResume:          nil,
				DisableForm:            data.NewBool(true),
				WorkerCount:            data.NewInt(WorkerCount),
				Policy: &storage.PutPolicy{
					Scope:               "",
					Expires:             0,
					IsPrefixalScope:     0,
					InsertOnly:          0,
					DetectMime:          0,
					FsizeMin:            0,
					FsizeLimit:          0,
					MimeLimit:           "",
					ForceSaveKey:        false,
					SaveKey:             "",
					CallbackFetchKey:    0,
					CallbackURL:         callbackUrl,
					CallbackHost:        callbackHost,
					CallbackBody:        "",
					CallbackBodyType:    "",
					ReturnURL:           "",
					ReturnBody:          "",
					PersistentOps:       "",
					PersistentNotifyURL: "",
					PersistentPipeline:  "",
					EndUser:             "",
					DeleteAfterDays:     0,
					FileType:            0,
				},
			}
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			if len(args) > 1 {
				info.Key = args[1]
			}
			if len(args) > 2 {
				info.FilePath = args[2]
			}
			operations.UploadFile(cfg, info)
		},
	}

	cmd.Flags().StringVarP(&info.MimeType, "mimetype", "t", "", "file mime type")

	cmd.Flags().BoolVarP(&resumeAPIV2, "resumable-api-v2", "", false, "use resumable upload v2 APIs to upload")
	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite the file of same key in bucket")
	cmd.Flags().Int64VarP(&ResumableAPIV2PartSize, "v2-part-size", "", data.BLOCK_SIZE, "the part size when use resumable upload v2 APIs to upload, default 4M")
	cmd.Flags().IntVarP(&FileType, "storage", "s", 0, "storage type")
	cmd.Flags().IntVarP(&WorkerCount, "worker", "c", 16, "worker count")
	cmd.Flags().StringVarP(&UpHost, "up-host", "u", "", "uphost")
	cmd.Flags().StringVarP(&callbackUrl, "callback-urls", "l", "", "upload callback urls, separated by comma")
	cmd.Flags().StringVarP(&callbackHost, "callback-host", "T", "", "upload callback host")
	return cmd
}

func init() {
	registerLoader(uploadCmdLoader)
}

func uploadCmdLoader(superCmd *cobra.Command, cfg *iqshell.Config) {
	superCmd.AddCommand(
		upload2CmdBuilder(cfg),
		uploadCmdBuilder(cfg),
		syncCmdBuilder(cfg),
		formUploadCmdBuilder(cfg),
		resumeUploadCmdBuilder(cfg),
	)
}
