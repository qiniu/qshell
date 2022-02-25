package cmd

import (
	"github.com/qiniu/qshell/v2/docs"
	"github.com/qiniu/qshell/v2/iqshell"
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
			if len(args) > 0 {
				cfg.UploadConfigFile = args[0]
			}
			info.GroupInfo.Force = true
			operations.BatchUpload(cfg, info)
		},
	}
	cmd.Flags().StringVarP(&info.GroupInfo.SuccessExportFilePath, "success-list", "s", "", "upload success (all) file list")
	cmd.Flags().StringVarP(&info.GroupInfo.FailExportFilePath, "failure-list", "f", "", "upload failure file list")
	cmd.Flags().StringVarP(&info.GroupInfo.OverrideExportFilePath, "overwrite-list", "w", "", "upload success (overwrite) file list")
	cmd.Flags().IntVarP(&info.GroupInfo.WorkCount, "worker", "c", 1, "worker count")
	cmd.Flags().StringVarP(&cfg.CmdCfg.Up.Policy.CallbackURL, "callback-urls", "l", "", "upload callback urls, separated by comma")
	cmd.Flags().StringVarP(&cfg.CmdCfg.Up.Policy.CallbackHost, "callback-host", "T", "", "upload callback host")
	return cmd
}

var uploadConfigMouldCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.BatchUploadConfigMouldInfo{}
	cmd := &cobra.Command{
		Use:   "qupload-config-mould",
		Short: "Get config mould of batch upload ",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			operations.BatchUploadConfigMould(cfg, info)
		},
	}
	return cmd
}

var upload2CmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.BatchUploadInfo{}
	cmd := &cobra.Command{
		Use:   "qupload2",
		Short: "Batch upload files to the qiniu bucket",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.QUpload2Type
			info.GroupInfo.Force = true
			operations.BatchUpload2(cfg, info)
		},
	}
	cmd.Flags().StringVar(&info.GroupInfo.SuccessExportFilePath, "success-list", "", "upload success file list")
	cmd.Flags().StringVar(&info.GroupInfo.FailExportFilePath, "failure-list", "", "upload failure file list")
	cmd.Flags().StringVar(&info.GroupInfo.OverrideExportFilePath, "overwrite-list", "", "upload success (overwrite) file list")
	cmd.Flags().IntVar(&info.GroupInfo.WorkCount, "thread-count", 1, "multiple thread count")

	cfg.CmdCfg.Up.ResumableAPIV2 = data.NewBool(false)
	cmd.Flags().BoolVarP((*bool)(cfg.CmdCfg.Up.ResumableAPIV2), "resumable-api-v2", "", false, "use resumable upload v2 APIs to upload")
	cfg.CmdCfg.Up.IgnoreDir = data.NewBool(false)
	cmd.Flags().BoolVar((*bool)(cfg.CmdCfg.Up.IgnoreDir), "ignore-dir", false, "ignore the dir in the dest file key")
	cfg.CmdCfg.Up.Overwrite = data.NewBool(false)
	cmd.Flags().BoolVar((*bool)(cfg.CmdCfg.Up.Overwrite), "overwrite", false, "overwrite the file of same key in bucket")
	cfg.CmdCfg.Up.CheckExists = data.NewBool(false)
	cmd.Flags().BoolVar((*bool)(cfg.CmdCfg.Up.CheckExists), "check-exists", false, "check file key whether in bucket before upload")
	cfg.CmdCfg.Up.CheckHash = data.NewBool(false)
	cmd.Flags().BoolVar((*bool)(cfg.CmdCfg.Up.CheckHash), "check-hash", false, "check hash")
	cfg.CmdCfg.Up.CheckSize = data.NewBool(false)
	cmd.Flags().BoolVar((*bool)(cfg.CmdCfg.Up.CheckSize), "check-size", false, "check file size")
	cfg.CmdCfg.Up.RescanLocal = data.NewBool(false)
	cmd.Flags().BoolVar((*bool)(cfg.CmdCfg.Up.RescanLocal), "rescan-local", false, "rescan local dir to upload newly add files")

	cfg.CmdCfg.Up.ResumableAPIV2PartSize = data.NewInt64(data.BLOCK_SIZE)
	cmd.Flags().Int64Var((*int64)(cfg.CmdCfg.Up.ResumableAPIV2PartSize), "resumable-api-v2-part-size", data.BLOCK_SIZE, "the part size when use resumable upload v2 APIs to upload")
	cfg.CmdCfg.Up.SrcDir = data.NewString("")
	cmd.Flags().StringVar((*string)(cfg.CmdCfg.Up.SrcDir), "src-dir", "", "src dir to upload")
	cfg.CmdCfg.Up.FileList = data.NewString("")
	cmd.Flags().StringVar((*string)(cfg.CmdCfg.Up.FileList), "file-list", "", "file list to upload")
	cfg.CmdCfg.Up.Bucket = data.NewString("")
	cmd.Flags().StringVar((*string)(cfg.CmdCfg.Up.Bucket), "bucket", "", "bucket")
	cfg.CmdCfg.Up.PutThreshold = data.NewInt64(0)
	cmd.Flags().Int64Var((*int64)(cfg.CmdCfg.Up.PutThreshold), "put-threshold", 0, "chunk upload threshold")
	cfg.CmdCfg.Up.KeyPrefix = data.NewString("")
	cmd.Flags().StringVar((*string)(cfg.CmdCfg.Up.KeyPrefix), "key-prefix", "", "key prefix prepended to dest file key")
	cfg.CmdCfg.Up.SkipFilePrefixes = data.NewString("")
	cmd.Flags().StringVar((*string)(cfg.CmdCfg.Up.SkipFilePrefixes), "skip-file-prefixes", "", "skip files with these file prefixes")
	cfg.CmdCfg.Up.SkipPathPrefixes = data.NewString("")
	cmd.Flags().StringVar((*string)(cfg.CmdCfg.Up.SkipPathPrefixes), "skip-path-prefixes", "", "skip files with these relative path prefixes")
	cfg.CmdCfg.Up.SkipFixedStrings = data.NewString("")
	cmd.Flags().StringVar((*string)(cfg.CmdCfg.Up.SkipFixedStrings), "skip-fixed-strings", "", "skip files with the fixed string in the name")
	cfg.CmdCfg.Up.SkipSuffixes = data.NewString("")
	cmd.Flags().StringVar((*string)(cfg.CmdCfg.Up.SkipSuffixes), "skip-suffixes", "", "skip files with these suffixes")
	cfg.CmdCfg.Up.UpHost = data.NewString("")
	cmd.Flags().StringVar((*string)(cfg.CmdCfg.Up.UpHost), "up-host", "", "upload host")
	cfg.CmdCfg.Up.LogFile = data.NewString("")
	cmd.Flags().StringVar((*string)(cfg.CmdCfg.Up.LogFile), "log-file", "", "log file")
	cfg.CmdCfg.Up.LogLevel = data.NewString("")
	cmd.Flags().StringVar((*string)(cfg.CmdCfg.Up.LogLevel), "log-level", "debug", "log level")
	cfg.CmdCfg.Up.LogRotate = data.NewInt(7)
	cmd.Flags().IntVar((*int)(cfg.CmdCfg.Up.LogRotate), "log-rotate", 7, "log rotate days")
	cfg.CmdCfg.Up.FileType = data.NewInt(0)
	cmd.Flags().IntVar((*int)(cfg.CmdCfg.Up.FileType), "file-type", 0, "set storage file type")
	cmd.Flags().StringVarP(&cfg.CmdCfg.Up.Policy.CallbackURL, "callback-urls", "l", "", "upload callback urls, separated by comma")
	cmd.Flags().StringVarP(&cfg.CmdCfg.Up.Policy.CallbackHost, "callback-host", "T", "", "upload callback host")
	//cmd.Flags().StringVar(&cfg.CmdCfg.Up.BindUpIp, "bind-up-ip", "", "upload host ip to bind")
	//cmd.Flags().StringVar(&cfg.CmdCfg.Up.BindRsIp, "bind-rs-ip", "", "rs host ip to bind")
	//cmd.Flags().StringVar(&cfg.CmdCfg.Up.BindNicIp, "bind-nic-ip", "", "local network interface card to bind")
	return cmd
}

var syncCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	resumeAPIV2 := false
	info := operations.SyncInfo{}
	cmd := &cobra.Command{
		Use:   "sync <SrcResUrl> <Buckets> [-k <Key>]",
		Short: "Sync big file to qiniu bucket",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.SyncType
			cfg.CmdCfg.Up.ResumableAPIV2 = data.NewBool(resumeAPIV2)
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
	cfg.CmdCfg.Up.UpHost = data.NewString("")
	cmd.Flags().StringVarP((*string)(cfg.CmdCfg.Up.UpHost), "uphost", "u", "", "upload host")
	return cmd
}

var formUploadCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	overwrite := false
	info := operations.UploadInfo{}
	cmd := &cobra.Command{
		Use:   "fput <Bucket> <Key> <LocalFile>",
		Short: "Form upload a local file",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.FormPutType
			cfg.CmdCfg.Up.DisableResume = data.NewBool(true)
			cfg.CmdCfg.Up.Overwrite = data.NewBool(overwrite)
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
	cfg.CmdCfg.Up.FileType = data.NewInt(0)
	cmd.Flags().IntVarP((*int)(cfg.CmdCfg.Up.FileType), "storage", "s", 0, "storage type")
	cfg.CmdCfg.Up.UpHost = data.NewString("")
	cmd.Flags().StringVarP((*string)(cfg.CmdCfg.Up.UpHost), "up-host", "u", "", "uphost")
	cmd.Flags().StringVarP(&cfg.CmdCfg.Up.Policy.CallbackURL, "callback-urls", "l", "", "upload callback urls, separated by comma")
	cmd.Flags().StringVarP(&cfg.CmdCfg.Up.Policy.CallbackHost, "callback-host", "T", "", "upload callback host")

	return cmd
}

var resumeUploadCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var (
		resumeAPIV2 = false
		overwrite   = false
	)
	info := operations.UploadInfo{}
	cmd := &cobra.Command{
		Use:   "rput <Bucket> <Key> <LocalFile>",
		Short: "Resumable upload a local file",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.RPutType
			cfg.CmdCfg.Up.DisableForm = data.NewBool(true)
			cfg.CmdCfg.Up.ResumableAPIV2 = data.NewBool(resumeAPIV2)
			cfg.CmdCfg.Up.Overwrite = data.NewBool(overwrite)
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

	cfg.CmdCfg.Up.ResumableAPIV2PartSize = data.NewInt64(data.BLOCK_SIZE)
	cmd.Flags().Int64VarP((*int64)(cfg.CmdCfg.Up.ResumableAPIV2PartSize), "v2-part-size", "", data.BLOCK_SIZE, "the part size when use resumable upload v2 APIs to upload, default 4M")
	cfg.CmdCfg.Up.FileType = data.NewInt(0)
	cmd.Flags().IntVarP((*int)(cfg.CmdCfg.Up.FileType), "storage", "s", 0, "storage type")
	cfg.CmdCfg.Up.WorkerCount = data.NewInt(0)
	cmd.Flags().IntVarP((*int)(cfg.CmdCfg.Up.WorkerCount), "worker", "c", 16, "worker count")
	cfg.CmdCfg.Up.UpHost = data.NewString("")
	cmd.Flags().StringVarP((*string)(cfg.CmdCfg.Up.UpHost), "up-host", "u", "", "uphost")
	cmd.Flags().StringVarP(&cfg.CmdCfg.Up.Policy.CallbackURL, "callback-urls", "l", "", "upload callback urls, separated by comma")
	cmd.Flags().StringVarP(&cfg.CmdCfg.Up.Policy.CallbackHost, "callback-host", "T", "", "upload callback host")
	return cmd
}

func init() {
	registerLoader(uploadCmdLoader)
}

func uploadCmdLoader(superCmd *cobra.Command, cfg *iqshell.Config) {
	superCmd.AddCommand(
		upload2CmdBuilder(cfg),
		uploadCmdBuilder(cfg),
		uploadConfigMouldCmdBuilder(cfg),
		syncCmdBuilder(cfg),
		formUploadCmdBuilder(cfg),
		resumeUploadCmdBuilder(cfg),
	)
}
