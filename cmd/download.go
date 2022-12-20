package cmd

import (
	"github.com/qiniu/qshell/v2/docs"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/config"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/download/operations"
	"github.com/spf13/cobra"
)

var downloadCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.BatchDownloadWithConfigInfo{}
	cmd := &cobra.Command{
		Use:   "qdownload [-c <ThreadCount>] <LocalDownloadConfig>",
		Short: "Batch download files from the qiniu bucket",
		Long: `By default qdownload use 5 goroutines to download, it can be customized use -c <count> flag.
And qdownload will use batch stat api or list api to get files info so that it have knowledge to tell whether files
have already in local disk and need to skip download or not.`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.QDownloadType
			cfg.CmdCfg.Log.LogLevel = data.NewString(config.DebugKey)
			info.Force = true
			if len(args) > 0 {
				info.LocalDownloadConfig = args[0]
			}
			operations.BatchDownloadWithConfig(cfg, info)
		},
	}
	cmd.Flags().IntVarP(&info.WorkerCount, "thread", "c", 5, "num of threads to download files")
	return cmd
}

var download2CmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	var (
		LogFile   = ""
		LogLevel  = ""
		LogRotate = 7
	)

	info := operations.BatchDownloadInfo{}
	cmd := &cobra.Command{
		Use:   "qdownload2 [-c <ThreadCount>] ",
		Short: "Batch download files from the qiniu bucket",
		Long: `By default qdownload2 use 5 goroutines to download, it can be customized use -c <count> flag.
And qdownload2 will use batch stat api or list api to get files info so that it have knowledge to tell whether files
have already in local disk and need to skip download or not.`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.QDownload2Type
			info.Force = true
			cfg.CmdCfg.Log = &config.LogSetting{
				LogLevel:  data.NewString(LogLevel),
				LogFile:   data.NewString(LogFile),
				LogRotate: data.NewInt(LogRotate),
				LogStdout: data.NewBool(true),
			}
			operations.BatchDownload(cfg, info)
		},
	}
	cmd.Flags().IntVarP(&info.WorkerCount, "thread", "c", 5, "num of threads to download files")
	cmd.Flags().StringVarP(&info.DownloadCfg.DestDir, "dest-dir", "", "", "local storage path, full path. default current dir")
	cmd.Flags().BoolVarP(&info.DownloadCfg.GetFileApi, "get-file-api", "", false, "public storage cloud not support, private storage cloud support when has getfile api.")
	cmd.Flags().StringVarP(&info.DownloadCfg.Bucket, "bucket", "", "", "storage bucket")
	cmd.Flags().StringVarP(&info.DownloadCfg.Prefix, "prefix", "", "", "only download files with the specified prefix")
	cmd.Flags().StringVarP(&info.DownloadCfg.Suffixes, "suffixes", "", "", "only download files with the specified suffixes")
	cmd.Flags().StringVarP(&info.DownloadCfg.KeyFile, "key-file", "", "", "configure a file and specify the keys to be downloaded; if not configured, download all the files in the bucket")
	cmd.Flags().StringVarP(&info.DownloadCfg.SavePathHandler, "save-path-handler", "", "", "specify a callback function; when constructing the save path of the file, this option is preferred for construction. If not configured, $dest_dir + $ file separator + $Key will be used for construction. This function is implemented through the template of the Go language. The func command is used for function verification. For the specific syntax, please refer to the description of the func command.")
	cmd.Flags().BoolVarP(&info.DownloadCfg.CheckHash, "check-hash", "", false, "whether to verify the hash, if it is enabled, it may take a long time")
	cmd.Flags().StringVarP(&info.IoHost, "io-host", "", "", "io host of request")
	cmd.Flags().StringVarP(&info.DownloadCfg.CdnDomain, "cdn-domain", "", "", "set the CDN domain name for downloading, the default is empty, which means downloading from the storage source site")
	cmd.Flags().StringVarP(&info.DownloadCfg.Referer, "referer", "", "", "if the CDN domain name is configured with domain name whitelist anti-leech, you need to specify a referer address that allows access")
	cmd.Flags().BoolVarP(&info.DownloadCfg.Public, "public", "", false, "whether the space is a public space")
	cmd.Flags().BoolVarP(&info.DownloadCfg.EnableSlice, "enable-slice", "", false, "whether to enable slice download, you need to pay attention to the configuration of `--slice-file-size-threshold` slice threshold option. Only when slice download is enabled and the size of the downloaded file is greater than the slice threshold will the slice download be started")
	cmd.Flags().Int64VarP(&info.DownloadCfg.SliceSize, "slice-size", "", 4*utils.MB, "slice size; when using slice download, the size of each slice; unit:B")
	cmd.Flags().IntVarP(&info.DownloadCfg.SliceConcurrentCount, "slice-concurrent-count", "", 10, "concurrency of slice downloads")
	cmd.Flags().Int64VarP(&info.DownloadCfg.SliceFileSizeThreshold, "slice-file-size-threshold", "", 40*utils.MB, "file threshold for downloading slices. When slice downloading is enabled and the file size is greater than this threshold, slice downloading will be enabled; unit:B")
	cmd.Flags().BoolVarP(&info.DownloadCfg.RemoveTempWhileError, "remove-temp-while-error", "", false, "when the download encounters an error, delete the previously downloaded part of the file cache")
	cmd.Flags().StringVarP(&info.DownloadCfg.RecordRoot, "record-root", "", "", "path to save download record information, including log files and download progress files; the default is `qshell` download directory")

	cmd.Flags().StringVarP(&LogLevel, "log-level", "", "debug", "download log output level, optional values are debug,info,warn and error")
	cmd.Flags().StringVarP(&LogFile, "log-file", "", "debug", "the output file of the download log is output to the file specified by record_root by default, and the specific file path can be seen in the terminal output")
	cmd.Flags().IntVarP(&LogRotate, "log-rotate", "", 7, "the switching period of the download log file, the unit is day,")

	return cmd
}

var getCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.DownloadInfo{}
	var cmd = &cobra.Command{
		Use:   "get <Bucket> <Key>",
		Short: "Download a single file from bucket",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.GetType
			if len(args) > 0 {
				info.Bucket = args[0]
			}
			if len(args) > 1 {
				info.Key = args[1]
			}
			operations.DownloadFile(cfg, info)
		},
	}

	cmd.Flags().StringVarP(&info.ToFile, "outfile", "o", "", "save file as specified by this option")
	cmd.Flags().StringVarP(&info.Domain, "domain", "", "", "domain of request")
	cmd.Flags().BoolVarP(&info.CheckHash, "check-hash", "", false, "check the consistency of the hash of the local file and the server file after downloading.")
	cmd.Flags().BoolVarP(&info.UseGetFileApi, "get-file-api", "", false, "public storage cloud not support, private storage cloud support when has getfile api.")
	cmd.Flags().BoolVarP(&info.EnableSlice, "enable-slice", "", false, "file download using slices, you need to pay attention to the setting of --slice-file-size-threshold. default is false")
	cmd.Flags().Int64VarP(&info.SliceFileSizeThreshold, "slice-file-size-threshold", "", 40*utils.MB, "the file size threshold that download using slices. when you use --enable-slice option, files larger than this size will be downloaded using slices. Unit: B")
	cmd.Flags().Int64VarP(&info.SliceSize, "slice-size", "", 4*utils.MB, "slice size that download using slices. when you use --enable-slice option, the file will be cut into data blocks according to the slice size, then the data blocks will be downloaded concurrently, and finally these data blocks will be spliced into a file. Unit: B")
	cmd.Flags().IntVarP(&info.SliceConcurrentCount, "slice-concurrent-count", "", 10, "the count of concurrently downloaded slices.")
	cmd.Flags().BoolVarP(&info.RemoveTempWhileError, "remove-temp-while-error", "", false, "remove download temp file while error happened, default is false")
	return cmd
}

func init() {
	registerLoader(downloadCmdLoader)
}

func downloadCmdLoader(superCmd *cobra.Command, cfg *iqshell.Config) {
	superCmd.AddCommand(
		getCmdBuilder(cfg),
		downloadCmdBuilder(cfg),
		download2CmdBuilder(cfg),
	)
}
