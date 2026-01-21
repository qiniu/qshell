package cmd

import (
	"github.com/spf13/cobra"

	"github.com/qiniu/qshell/v2/docs"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/config"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/upload/operations"
)

var uploadCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.BatchUploadInfo{}
	cmd := &cobra.Command{
		Use:   "qupload <LocalUploadConfig>",
		Short: "Batch upload files to the qiniu bucket",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.QUploadType
			cfg.CmdCfg.Log.LogLevel = data.NewString(config.DebugKey)
			cfg.CmdCfg.Log.LogStdout = data.NewBool(true)
			cfg.CmdCfg.Log.LogRotate = data.NewInt(7)
			info.Force = true
			if len(args) > 0 {
				info.UploadConfigFile = args[0]
			}
			operations.BatchUpload(cfg, info)
		},
	}
	cmd.Flags().StringVarP(&info.SuccessExportFilePath, "success-list", "s", "", "specifies the file path where the successful file list is saved")

	cmd.Flags().StringVarP(&info.FailExportFilePath, "failure-list", "e", "", "specifies the file path where the failure file list is saved")
	cmd.Flags().StringVarP(&info.FailExportFilePath, "failure-list-old", "f", "", "specifies the file path where the failure file list is saved, deprecated")
	_ = cmd.Flags().MarkDeprecated("failure-list-old", "use --failure-list instead")

	cmd.Flags().StringVarP(&info.OverwriteExportFilePath, "overwrite-list", "w", "", "specifies the file path where the overwrite file list is saved")
	cmd.Flags().IntVarP(&info.Info.WorkerCount, "worker", "c", 1, "worker count")
	cmd.Flags().StringVarP(&info.CallbackUrl, "callback-urls", "l", "", "upload callback urls, separated by comma")
	cmd.Flags().StringVarP(&info.CallbackHost, "callback-host", "T", "", "upload callback host")
	return cmd
}

var upload2CmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var (
		LogFile   = ""
		LogLevel  = ""
		LogRotate = 7
	)
	info := operations.BatchUpload2Info{
		UploadConfig: operations.UploadConfig{},
	}
	cmd := &cobra.Command{
		Use:   "qupload2",
		Short: "Batch upload files to the qiniu bucket",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.QUpload2Type
			info.Force = true
			cfg.CmdCfg.Log = &config.LogSetting{
				LogLevel:  data.NewString(LogLevel),
				LogFile:   data.NewString(LogFile),
				LogRotate: data.NewInt(LogRotate),
				LogStdout: data.NewBool(true),
			}
			operations.BatchUpload2(cfg, info)
		},
	}
	cmd.Flags().StringVarP(&info.SuccessExportFilePath, "success-list", "s", "", "upload success file list")
	cmd.Flags().StringVarP(&info.FailExportFilePath, "failure-list", "e", "", "upload failure file list")
	cmd.Flags().StringVarP(&info.OverwriteExportFilePath, "overwrite-list", "w", "", "upload success (overwrite) file list")
	cmd.Flags().IntVar(&info.Info.WorkerCount, "thread-count", 1, "multiple thread count")
	cmd.Flags().IntVar(&info.UploadConfig.WorkerCount, "worker-count", 3, "the number of concurrently uploaded parts of a single file in resumable upload")
	cmd.Flags().BoolVar(&info.UploadConfig.SequentialReadFile, "sequential-read-file", false, "File reading is sequential and does not involve skipping; when enabled, the uploading fragment data will be loaded into the memory. This option may increase file upload speed for mounted network filesystems.")

	cmd.Flags().BoolVarP(&info.ResumableAPIV2, "resumable-api-v2", "", true, "use resumable upload v2 APIs to upload, default is true")
	cmd.Flags().Int64Var(&info.ResumableAPIV2PartSize, "resumable-api-v2-part-size", data.BLOCK_SIZE, "the part size when use resumable upload v2 APIs to upload")
	cmd.Flags().BoolVar(&info.IgnoreDir, "ignore-dir", false, "ignore the dir in the dest file key")
	cmd.Flags().BoolVarP(&info.Overwrite, "overwrite", "", false, "overwrite the file of same key in bucket")
	cmd.Flags().BoolVar(&info.CheckExists, "check-exists", false, "check file key whether in bucket before upload")
	cmd.Flags().BoolVar(&info.CheckHash, "check-hash", false, "check hash")
	cmd.Flags().BoolVar(&info.CheckSize, "check-size", false, "check file size")
	cmd.Flags().BoolVar(&info.RescanLocal, "rescan-local", false, "rescan local dir to upload newly add files")

	cmd.Flags().StringVar(&info.SrcDir, "src-dir", "", "src dir to upload")
	cmd.Flags().StringVar(&info.FileList, "file-list", "", "file list to upload")
	cmd.Flags().StringVar(&info.Bucket, "bucket", "", "bucket")
	cmd.Flags().Int64Var(&info.PutThreshold, "put-threshold", 8*1024*1024, "chunk upload threshold, unit: B")
	cmd.Flags().StringVar(&info.KeyPrefix, "key-prefix", "", "key prefix prepended to dest file key")
	cmd.Flags().StringVar(&info.SkipFilePrefixes, "skip-file-prefixes", "", "skip files with these file prefixes")
	cmd.Flags().StringVar(&info.SkipPathPrefixes, "skip-path-prefixes", "", "skip files with these relative path prefixes")
	cmd.Flags().StringVar(&info.SkipFixedStrings, "skip-fixed-strings", "", "skip files with the fixed string in the name")
	cmd.Flags().StringVar(&info.SkipSuffixes, "skip-suffixes", "", "skip files with these suffixes")
	cmd.Flags().StringVar(&info.UpHost, "up-host", "", "upload host")
	cmd.Flags().BoolVarP(&info.Accelerate, "accelerate", "", false, "enable uploading acceleration")
	cmd.Flags().StringVar(&info.RecordRoot, "record-root", "", "record root dir, and will save record info to the dir(db and log), default <UserRoot>/.qshell")
	cmd.Flags().StringVar(&LogFile, "log-file", "", "log file")
	cmd.Flags().StringVar(&LogLevel, "log-level", "debug", "log level")
	cmd.Flags().IntVar(&LogRotate, "log-rotate", 7, "log rotate days")

	cmd.Flags().IntVarP(&info.FileType, "file-type", "", 0, "set storage type of file, 0:STANDARD storage, 1:IA storage, 2:ARCHIVE storage, 3:DEEP_ARCHIVE storage, 4:ARCHIVE_IR storage, 5:INTELLIGENT_TIERING")
	cmd.Flags().IntVarP(&info.FileType, "storage", "", 0, "set storage type of file, same to --file-type")
	_ = cmd.Flags().MarkDeprecated("storage", "use --file-type instead") // 废弃 storage

	// cmd.Flags().StringVar(&cfg.CmdCfg.Up.BindUpIp, "bind-up-ip", "", "upload host ip to bind")
	// cmd.Flags().StringVar(&cfg.CmdCfg.Up.BindRsIp, "bind-rs-ip", "", "rs host ip to bind")
	// cmd.Flags().StringVar(&cfg.CmdCfg.Up.BindNicIp, "bind-nic-ip", "", "local network interface card to bind")

	cmd.Flags().StringVarP(&info.EndUser, "end-user", "", "", "Owner identification")
	cmd.Flags().StringVarP(&info.CallbackURL, "callback-urls", "l", "", "upload callback urls, separated by comma")
	cmd.Flags().StringVarP(&info.CallbackHost, "callback-host", "T", "", "upload callback host")
	cmd.Flags().StringVarP(&info.CallbackBody, "callback-body", "", "", "upload callback body")
	cmd.Flags().StringVarP(&info.CallbackBodyType, "callback-body-type", "", "", "upload callback body type")
	cmd.Flags().StringVarP(&info.PersistentOps, "persistent-ops", "", "", "List of pre-transfer persistence processing instructions that are triggered after successful resource upload. This parameter is not supported when fileType=2 or 3 (upload archive storage or deep archive storage files). Supports magic variables and custom variables. Each directive is an API specification string, and multiple directives are separated by ;.")
	cmd.Flags().StringVarP(&info.PersistentNotifyURL, "persistent-notify-url", "", "", "URL to receive notification of persistence processing results. It must be a valid URL that can make POST requests normally on the public Internet and respond successfully. The content obtained by this URL is consistent with the processing result of the persistence processing status query. To send a POST request whose body format is application/json, you need to read the body of the request in the form of a read stream to obtain it.")
	cmd.Flags().StringVarP(&info.PersistentPipeline, "persistent-pipeline", "", "", "Transcoding queue name. After the resource is successfully uploaded, an independent queue is designated for transcoding when transcoding is triggered. If it is empty, it means that the public queue is used, and the processing speed is slower. It is recommended to use a dedicated queue.")
	cmd.Flags().IntVarP(&info.DetectMime, "detect-mime", "", 0, `Turn on the MimeType detection function and perform detection according to the following rules; if the correct value cannot be detected, application/octet-stream will be used by default.
If set to a value of 1, the file MimeType information passed by the uploader will be ignored, and the MimeType value will be detected in the following order:
	1. Detection content;
	2. Check the file extension;
	3. Check the Key extension.
The default value is set to 0.If the uploader specifies MimeType (except application/octet-stream), this value will be used directly. Otherwise, the MimeType value will be detected in the following order:
	1. Check the file extension;
	2. Check the Key extension;
	3. Detect content.
Set to a value of -1 and use this value regardless of what value is specified on the uploader.`)
	cmd.Flags().Uint64VarP(&info.TrafficLimit, "traffic-limit", "", 0, "Upload request single link speed limit to control client bandwidth usage. The speed limit value range is 819200 ~ 838860800, and the unit is bit/s.")
	return cmd
}

var syncCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.SyncInfo{}
	cmd := &cobra.Command{
		Use:   "sync <SrcResUrl> <Buckets> [-k <Key>]",
		Short: "Sync big file to qiniu bucket",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.SyncType
			info.DisableResume = true
			if len(args) > 0 {
				info.FilePath = args[0]
			}
			if len(args) > 1 {
				info.ToBucket = args[1]
			}
			operations.SyncFile(cfg, info)
		},
	}
	cmd.Flags().StringVarP(&info.SaveKey, "key", "k", "", "save as <key> in bucket")
	cmd.Flags().BoolVarP(&info.UseResumeV2, "resumable-api-v2", "", true, "use resumable upload v2 APIs to upload, default is true")
	cmd.Flags().Int64VarP(&info.ChunkSize, "resumable-api-v2-part-size", "", data.BLOCK_SIZE, "the part size when use resumable upload v2 APIs to upload, default 4M")
	cmd.Flags().StringVarP(&info.UpHost, "up-host", "u", "", "upload host")
	cmd.Flags().BoolVarP(&info.Accelerate, "accelerate", "", false, "enable uploading acceleration")

	cmd.Flags().IntVarP(&info.FileType, "file-type", "", 0, "set storage type of file, 0:STANDARD storage, 1:IA storage, 2:ARCHIVE storage, 3:DEEP_ARCHIVE storage, 4:ARCHIVE_IR storage, 5:INTELLIGENT_TIERING")
	cmd.Flags().IntVarP(&info.FileType, "storage", "s", 0, "set storage type of file, same to --file-type")
	_ = cmd.Flags().MarkDeprecated("storage", "use --file-type instead") // 废弃 storage

	cmd.Flags().BoolVarP(&info.Overwrite, "overwrite", "", false, "overwrite the file of same key in bucket")

	cmd.Flags().StringVarP(&info.Policy.EndUser, "end-user", "", "", "Owner identification")
	cmd.Flags().StringVarP(&info.Policy.CallbackURL, "callback-urls", "l", "", "upload callback urls, separated by comma")
	cmd.Flags().StringVarP(&info.Policy.CallbackHost, "callback-host", "T", "", "upload callback host")
	cmd.Flags().StringVarP(&info.Policy.CallbackBody, "callback-body", "", "", "upload callback body")
	cmd.Flags().StringVarP(&info.Policy.CallbackBodyType, "callback-body-type", "", "", "upload callback body type")
	cmd.Flags().StringVarP(&info.Policy.PersistentOps, "persistent-ops", "", "", "List of pre-transfer persistence processing instructions that are triggered after successful resource upload. This parameter is not supported when fileType=2 or 3 (upload archive storage or deep archive storage files). Supports magic variables and custom variables. Each directive is an API specification string, and multiple directives are separated by ;.")
	cmd.Flags().StringVarP(&info.Policy.PersistentNotifyURL, "persistent-notify-url", "", "", "URL to receive notification of persistence processing results. It must be a valid URL that can make POST requests normally on the public Internet and respond successfully. The content obtained by this URL is consistent with the processing result of the persistence processing status query. To send a POST request whose body format is application/json, you need to read the body of the request in the form of a read stream to obtain it.")
	cmd.Flags().StringVarP(&info.Policy.PersistentPipeline, "persistent-pipeline", "", "", "Transcoding queue name. After the resource is successfully uploaded, an independent queue is designated for transcoding when transcoding is triggered. If it is empty, it means that the public queue is used, and the processing speed is slower. It is recommended to use a dedicated queue.")
	cmd.Flags().IntVarP(&info.Policy.DetectMime, "detect-mime", "", 0, `Turn on the MimeType detection function and perform detection according to the following rules; if the correct value cannot be detected, application/octet-stream will be used by default.
If set to a value of 1, the file MimeType information passed by the uploader will be ignored, and the MimeType value will be detected in the following order:
	1. Detection content;
	2. Check the file extension;
	3. Check the Key extension.
The default value is set to 0.If the uploader specifies MimeType (except application/octet-stream), this value will be used directly. Otherwise, the MimeType value will be detected in the following order:
	1. Check the file extension;
	2. Check the Key extension;
	3. Detect content.
Set to a value of -1 and use this value regardless of what value is specified on the uploader.`)
	cmd.Flags().Uint64VarP(&info.Policy.TrafficLimit, "traffic-limit", "", 0, "Upload request single link speed limit to control client bandwidth usage. The speed limit value range is 819200 ~ 838860800, and the unit is bit/s.")
	return cmd
}

var formUploadCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.UploadInfo{}
	cmd := &cobra.Command{
		Use:   "fput <Bucket> <Key> <LocalFile>",
		Short: "Form upload a local file",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.FormPutType
			info.DisableResume = true
			if len(args) > 0 {
				info.ToBucket = args[0]
			}
			if len(args) > 1 {
				info.SaveKey = args[1]
			}
			if len(args) > 2 {
				info.FilePath = args[2]
			}
			operations.UploadFile(cfg, info)
		},
	}
	cmd.Flags().BoolVarP(&info.Overwrite, "overwrite", "", false, "overwrite the file of same key in bucket")
	cmd.Flags().StringVarP(&info.MimeType, "mimetype", "t", "", "file mime type")

	cmd.Flags().IntVarP(&info.FileType, "file-type", "", 0, "set storage type of file, 0:STANDARD storage, 1:IA storage, 2:ARCHIVE storage, 3:DEEP_ARCHIVE storage, 4:ARCHIVE_IR storage, 5:INTELLIGENT_TIERING")
	cmd.Flags().IntVarP(&info.FileType, "storage", "s", 0, "set storage type of file, same to --file-type")
	_ = cmd.Flags().MarkDeprecated("storage", "use --file-type instead") // 废弃 storage

	cmd.Flags().StringVarP(&info.UpHost, "up-host", "u", "", "uphost")
	cmd.Flags().BoolVarP(&info.Accelerate, "accelerate", "", false, "enable uploading acceleration")

	cmd.Flags().StringVarP(&info.Policy.EndUser, "end-user", "", "", "Owner identification")
	cmd.Flags().StringVarP(&info.Policy.CallbackURL, "callback-urls", "l", "", "upload callback urls, separated by comma")
	cmd.Flags().StringVarP(&info.Policy.CallbackHost, "callback-host", "T", "", "upload callback host")
	cmd.Flags().StringVarP(&info.Policy.CallbackBody, "callback-body", "", "", "upload callback body")
	cmd.Flags().StringVarP(&info.Policy.CallbackBodyType, "callback-body-type", "", "", "upload callback body type")
	cmd.Flags().StringVarP(&info.Policy.PersistentOps, "persistent-ops", "", "", "List of pre-transfer persistence processing instructions that are triggered after successful resource upload. This parameter is not supported when fileType=2 or 3 (upload archive storage or deep archive storage files). Supports magic variables and custom variables. Each directive is an API specification string, and multiple directives are separated by ;.")
	cmd.Flags().StringVarP(&info.Policy.PersistentNotifyURL, "persistent-notify-url", "", "", "URL to receive notification of persistence processing results. It must be a valid URL that can make POST requests normally on the public Internet and respond successfully. The content obtained by this URL is consistent with the processing result of the persistence processing status query. To send a POST request whose body format is application/json, you need to read the body of the request in the form of a read stream to obtain it.")
	cmd.Flags().StringVarP(&info.Policy.PersistentPipeline, "persistent-pipeline", "", "", "Transcoding queue name. After the resource is successfully uploaded, an independent queue is designated for transcoding when transcoding is triggered. If it is empty, it means that the public queue is used, and the processing speed is slower. It is recommended to use a dedicated queue.")
	cmd.Flags().IntVarP(&info.Policy.DetectMime, "detect-mime", "", 0, `Turn on the MimeType detection function and perform detection according to the following rules; if the correct value cannot be detected, application/octet-stream will be used by default.
If set to a value of 1, the file MimeType information passed by the uploader will be ignored, and the MimeType value will be detected in the following order:
	1. Detection content;
	2. Check the file extension;
	3. Check the Key extension.
The default value is set to 0.If the uploader specifies MimeType (except application/octet-stream), this value will be used directly. Otherwise, the MimeType value will be detected in the following order:
	1. Check the file extension;
	2. Check the Key extension;
	3. Detect content.
Set to a value of -1 and use this value regardless of what value is specified on the uploader.`)
	cmd.Flags().Uint64VarP(&info.Policy.TrafficLimit, "traffic-limit", "", 0, "Upload request single link speed limit to control client bandwidth usage. The speed limit value range is 819200 ~ 838860800, and the unit is bit/s.")
	return cmd
}

var resumeUploadCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.UploadInfo{}
	cmd := &cobra.Command{
		Use:   "rput <Bucket> <Key> <LocalFile>",
		Short: "Resumable upload a local file",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.RPutType
			info.DisableForm = true
			if len(args) > 0 {
				info.ToBucket = args[0]
			}
			if len(args) > 1 {
				info.SaveKey = args[1]
			}
			if len(args) > 2 {
				info.FilePath = args[2]
			}
			operations.UploadFile(cfg, info)
		},
	}

	cmd.Flags().StringVarP(&info.MimeType, "mimetype", "t", "", "file mime type")
	cmd.Flags().BoolVarP(&info.Overwrite, "overwrite", "", false, "overwrite the file of same key in bucket")
	cmd.Flags().BoolVarP(&info.UseResumeV2, "resumable-api-v2", "", true, "use resumable upload v2 APIs to upload, default is true")
	cmd.Flags().BoolVar(&info.SequentialReadFile, "sequential-read-file", false, "File reading is sequential and does not involve skipping; when enabled, the uploading fragment data will be loaded into the memory. This option may increase file upload speed for mounted network filesystems.")

	cmd.Flags().Int64VarP(&info.ChunkSize, "resumable-api-v2-part-size", "", data.BLOCK_SIZE, "the part size when use resumable upload v2 APIs to upload, default 4M")
	cmd.Flags().Int64VarP(&info.ChunkSize, "v2-part-size", "", data.BLOCK_SIZE, "the part size when use resumable upload v2 APIs to upload, same to --resumable-api-v2-part-size")
	_ = cmd.Flags().MarkDeprecated("v2-part-size", "use --resumable-api-v2-part-size instead")

	cmd.Flags().IntVarP(&info.FileType, "file-type", "", 0, "set storage type of file, 0:STANDARD storage, 1:IA storage, 2:ARCHIVE storage, 3:DEEP_ARCHIVE storage, 4:ARCHIVE_IR storage, 5:INTELLIGENT_TIERING")
	cmd.Flags().IntVarP(&info.FileType, "storage", "s", 0, "set storage type of file, same to --file-type")
	_ = cmd.Flags().MarkDeprecated("storage", "use --file-type instead") // 废弃 storage

	cmd.Flags().IntVarP(&info.ResumeWorkerCount, "worker", "c", 3, "worker count")
	cmd.Flags().StringVarP(&info.UpHost, "up-host", "u", "", "uphost")
	cmd.Flags().BoolVarP(&info.Accelerate, "accelerate", "", false, "enable uploading acceleration")

	cmd.Flags().StringVarP(&info.Policy.EndUser, "end-user", "", "", "Owner identification")
	cmd.Flags().StringVarP(&info.Policy.CallbackURL, "callback-urls", "l", "", "upload callback urls, separated by comma")
	cmd.Flags().StringVarP(&info.Policy.CallbackHost, "callback-host", "T", "", "upload callback host")
	cmd.Flags().StringVarP(&info.Policy.CallbackBody, "callback-body", "", "", "upload callback body")
	cmd.Flags().StringVarP(&info.Policy.CallbackBodyType, "callback-body-type", "", "", "upload callback body type")
	cmd.Flags().StringVarP(&info.Policy.PersistentOps, "persistent-ops", "", "", "List of pre-transfer persistence processing instructions that are triggered after successful resource upload. This parameter is not supported when fileType=2 or 3 (upload archive storage or deep archive storage files). Supports magic variables and custom variables. Each directive is an API specification string, and multiple directives are separated by ;.")
	cmd.Flags().StringVarP(&info.Policy.PersistentNotifyURL, "persistent-notify-url", "", "", "URL to receive notification of persistence processing results. It must be a valid URL that can make POST requests normally on the public Internet and respond successfully. The content obtained by this URL is consistent with the processing result of the persistence processing status query. To send a POST request whose body format is application/json, you need to read the body of the request in the form of a read stream to obtain it.")
	cmd.Flags().StringVarP(&info.Policy.PersistentPipeline, "persistent-pipeline", "", "", "Transcoding queue name. After the resource is successfully uploaded, an independent queue is designated for transcoding when transcoding is triggered. If it is empty, it means that the public queue is used, and the processing speed is slower. It is recommended to use a dedicated queue.")
	cmd.Flags().IntVarP(&info.Policy.DetectMime, "detect-mime", "", 0, `Turn on the MimeType detection function and perform detection according to the following rules; if the correct value cannot be detected, application/octet-stream will be used by default.
If set to a value of 1, the file MimeType information passed by the uploader will be ignored, and the MimeType value will be detected in the following order:
	1. Detection content;
	2. Check the file extension;
	3. Check the Key extension.
The default value is set to 0. If the uploader specifies MimeType (except application/octet-stream), this value will be used directly. Otherwise, the MimeType value will be detected in the following order:
	1. Check the file extension;
	2. Check the Key extension;
	3. Detect content.
Set to a value of -1 and use this value regardless of what value is specified on the uploader.`)
	cmd.Flags().Uint64VarP(&info.Policy.TrafficLimit, "traffic-limit", "", 0, "Upload request single link speed limit to control client bandwidth usage. The speed limit value range is 819200 ~ 838860800, and the unit is bit/s.")
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
