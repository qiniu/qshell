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
	info := operations.BatchUploadInfo{}
	cmd := &cobra.Command{
		Use:   "qupload <LocalDownloadConfig>",
		Short: "Batch upload files to the qiniu bucket",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.QUploadType
			cfg.CmdCfg.Log.LogLevel = data.NewString(config.DebugKey)
			cfg.CmdCfg.Log.LogStdout = data.NewBool(true)
			cfg.CmdCfg.Log.LogRotate = data.NewInt(7)
			if len(args) > 0 {
				info.UploadConfigFile = args[0]
			}
			info.BatchInfo.Force = true
			operations.BatchUpload(cfg, info)
		},
	}
	cmd.Flags().StringVarP(&info.BatchInfo.SuccessExportFilePath, "success-list", "s", "", "upload success (all) file list")
	cmd.Flags().StringVarP(&info.BatchInfo.FailExportFilePath, "failure-list", "f", "", "upload failure file list")
	cmd.Flags().StringVarP(&info.BatchInfo.OverwriteExportFilePath, "overwrite-list", "w", "", "upload success (overwrite) file list")
	cmd.Flags().IntVarP(&info.BatchInfo.WorkerCount, "worker", "c", 1, "worker count")
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
		UploadConfig: operations.UploadConfig{
			Policy: &storage.PutPolicy{},
		},
	}
	cmd := &cobra.Command{
		Use:   "qupload2",
		Short: "Batch upload files to the qiniu bucket",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.QUpload2Type
			info.BatchInfo.Force = true
			cfg.CmdCfg.Log = &config.LogSetting{
				LogLevel:  data.NewString(LogLevel),
				LogFile:   data.NewString(LogFile),
				LogRotate: data.NewInt(LogRotate),
				LogStdout: data.NewBool(true),
			}
			operations.BatchUpload2(cfg, info)
		},
	}
	cmd.Flags().StringVar(&info.BatchInfo.SuccessExportFilePath, "success-list", "", "upload success file list")
	cmd.Flags().StringVar(&info.BatchInfo.FailExportFilePath, "failure-list", "", "upload failure file list")
	cmd.Flags().StringVar(&info.BatchInfo.OverwriteExportFilePath, "overwrite-list", "", "upload success (overwrite) file list")
	cmd.Flags().IntVar(&info.BatchInfo.WorkerCount, "thread-count", 1, "multiple thread count")

	cmd.Flags().BoolVarP(&info.ResumableAPIV2, "resumable-api-v2", "", false, "use resumable upload v2 APIs to upload")
	cmd.Flags().BoolVar(&info.IgnoreDir, "ignore-dir", false, "ignore the dir in the dest file key")
	cmd.Flags().BoolVar(&info.Overwrite, "overwrite", false, "overwrite the file of same key in bucket")
	cmd.Flags().BoolVar(&info.CheckExists, "check-exists", false, "check file key whether in bucket before upload")
	cmd.Flags().BoolVar(&info.CheckHash, "check-hash", false, "check hash")
	cmd.Flags().BoolVar(&info.CheckSize, "check-size", false, "check file size")
	cmd.Flags().BoolVar(&info.RescanLocal, "rescan-local", false, "rescan local dir to upload newly add files")

	cmd.Flags().Int64Var(&info.ResumableAPIV2PartSize, "resumable-api-v2-part-size", data.BLOCK_SIZE, "the part size when use resumable upload v2 APIs to upload")
	cmd.Flags().StringVar(&info.SrcDir, "src-dir", "", "src dir to upload")
	cmd.Flags().StringVar(&info.FileList, "file-list", "", "file list to upload")
	cmd.Flags().StringVar(&info.Bucket, "bucket", "", "bucket")
	cmd.Flags().Int64Var(&info.PutThreshold, "put-threshold", 0, "chunk upload threshold")
	cmd.Flags().StringVar(&info.KeyPrefix, "key-prefix", "", "key prefix prepended to dest file key")
	cmd.Flags().StringVar(&info.SkipFilePrefixes, "skip-file-prefixes", "", "skip files with these file prefixes")
	cmd.Flags().StringVar(&info.SkipPathPrefixes, "skip-path-prefixes", "", "skip files with these relative path prefixes")
	cmd.Flags().StringVar(&info.SkipFixedStrings, "skip-fixed-strings", "", "skip files with the fixed string in the name")
	cmd.Flags().StringVar(&info.SkipSuffixes, "skip-suffixes", "", "skip files with these suffixes")
	cmd.Flags().StringVar(&info.UpHost, "up-host", "", "upload host")
	cmd.Flags().StringVar(&info.RecordRoot, "record-root", "", "record root dir, and will save record info to the dir(db and log), default <UserRoot>/.qshell")
	cmd.Flags().StringVar(&LogFile, "log-file", "", "log file")
	cmd.Flags().StringVar(&LogLevel, "log-level", "debug", "log level")
	cmd.Flags().IntVar(&LogRotate, "log-rotate", 7, "log rotate days")

	cmd.Flags().IntVar(&info.FileType, "file-type", 0, "set storage file type")
	cmd.Flags().StringVarP(&info.Policy.CallbackURL, "callback-urls", "l", "", "upload callback urls, separated by comma")
	cmd.Flags().StringVarP(&info.Policy.CallbackHost, "callback-host", "T", "", "upload callback host")
	//cmd.Flags().StringVar(&cfg.CmdCfg.Up.BindUpIp, "bind-up-ip", "", "upload host ip to bind")
	//cmd.Flags().StringVar(&cfg.CmdCfg.Up.BindRsIp, "bind-rs-ip", "", "rs host ip to bind")
	//cmd.Flags().StringVar(&cfg.CmdCfg.Up.BindNicIp, "bind-nic-ip", "", "local network interface card to bind")
	return cmd
}

var syncCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.SyncInfo{
		Policy: &storage.PutPolicy{},
	}
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
	cmd.Flags().BoolVarP(&info.UseResumeV2, "resumable-api-v2", "", false, "use resumable upload v2 APIs to upload")
	cmd.Flags().StringVarP(&info.UpHost, "up-host", "u", "", "upload host")
	cmd.Flags().BoolVar(&info.Overwrite, "overwrite", false, "overwrite the file of same key in bucket")

	return cmd
}

var formUploadCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.UploadInfo{
		Policy: &storage.PutPolicy{},
	}
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
	cmd.Flags().BoolVar(&info.Overwrite, "overwrite", false, "overwrite the file of same key in bucket")
	cmd.Flags().StringVarP(&info.MimeType, "mimetype", "t", "", "file mime type")
	cmd.Flags().IntVarP(&info.FileType, "storage", "s", 0, "storage type")
	cmd.Flags().StringVarP(&info.UpHost, "up-host", "u", "", "uphost")
	cmd.Flags().StringVarP(&info.Policy.CallbackURL, "callback-urls", "l", "", "upload callback urls, separated by comma")
	cmd.Flags().StringVarP(&info.Policy.CallbackHost, "callback-host", "T", "", "upload callback host")

	return cmd
}

var resumeUploadCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.UploadInfo{
		Policy: &storage.PutPolicy{},
	}
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

	cmd.Flags().BoolVarP(&info.UseResumeV2, "resumable-api-v2", "", false, "use resumable upload v2 APIs to upload")
	cmd.Flags().BoolVar(&info.Overwrite, "overwrite", false, "overwrite the file of same key in bucket")
	cmd.Flags().Int64VarP(&info.ChunkSize, "v2-part-size", "", data.BLOCK_SIZE, "the part size when use resumable upload v2 APIs to upload, default 4M")
	cmd.Flags().IntVarP(&info.FileType, "storage", "s", 0, "storage type")
	cmd.Flags().IntVarP(&info.ResumeWorkerCount, "worker", "c", 16, "worker count")
	cmd.Flags().StringVarP(&info.UpHost, "up-host", "u", "", "uphost")
	cmd.Flags().StringVarP(&info.Policy.CallbackURL, "callback-urls", "l", "", "upload callback urls, separated by comma")
	cmd.Flags().StringVarP(&info.Policy.CallbackHost, "callback-host", "T", "", "upload callback host")
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
