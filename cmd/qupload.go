package cmd

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/spf13/cobra"
	"github.com/tonycai653/iqshell/qshell"
	"io/ioutil"
	"os"
	"strconv"
)

var qUploadCmd = &cobra.Command{
	Use:   "qupload [<ThreadCount>] <LocalUploadConfig>",
	Short: "Batch upload files to the qiniu bucket",
	Run:   QiniuUpload,
}

var (
	successFname   string
	failureFname   string
	overwriteFname string
	upthreadCount  int64
)

func init() {
	qUploadCmd.Flags().StringVar(&successFname, "success-list", "", "upload success (all) file list")
	qUploadCmd.Flags().StringVar(&failureFname, "failure-list", "", "upload failure file list")
	qUploadCmd.Flags().StringVar(&overwriteFname, "overwrite-list", "", "upload success (overwrite) file list")
	qUploadCmd.Flags().Int64Var(&upthreadCount, "c", 1, "upload success (overwrite) file list")
	RootCmd.AddCommand(qUploadCmd)
}

func QiniuUpload(cmd *cobra.Command, params []string) {
	var uploadConfigFile string
	var err error
	if len(params) == 2 {
		upthreadCount, err = strconv.ParseInt(params[0], 10, 64)
		if err != nil {
			logs.Error("Invalid <ThreadCount> value,", params[0])
			os.Exit(2)
		}
		uploadConfigFile = params[1]
	} else {
		uploadConfigFile = params[0]
	}

	//read upload config
	fp, err := os.Open(uploadConfigFile)
	if err != nil {
		logs.Error("Open upload config file `%s` error due to `%s`", uploadConfigFile, err)
		os.Exit(qshell.STATUS_HALT)
	}
	defer fp.Close()
	configData, err := ioutil.ReadAll(fp)
	if err != nil {
		logs.Error("Read upload config file `%s` error due to `%s`", uploadConfigFile, err)
		os.Exit(qshell.STATUS_HALT)
	}
	var uploadConfig qshell.UploadConfig
	err = json.Unmarshal(configData, &uploadConfig)
	if err != nil {
		logs.Error("Parse upload config file `%s` errror due to `%s`", uploadConfigFile, err)
		os.Exit(qshell.STATUS_HALT)
	}

	if uploadConfig.FileType != 1 && uploadConfig.FileType != 0 {
		logs.Error("Wrong Filetype, It should be 0 or 1 ")
		os.Exit(qshell.STATUS_HALT)
	}

	srcFileInfo, err := os.Stat(uploadConfig.SrcDir)
	if err != nil {
		logs.Error("Upload config error for parameter `SrcDir`,", err)
		os.Exit(qshell.STATUS_HALT)
	}

	if !srcFileInfo.IsDir() {
		logs.Error("Upload src dir should be a directory")
		os.Exit(qshell.STATUS_HALT)
	}

	//upload
	if upthreadCount < qshell.MIN_UPLOAD_THREAD_COUNT || upthreadCount > qshell.MAX_UPLOAD_THREAD_COUNT {
		logs.Info("Tip: you can set <ThreadCount> value between %d and %d to improve speed\n",
			qshell.MIN_UPLOAD_THREAD_COUNT, qshell.MAX_UPLOAD_THREAD_COUNT)

		if upthreadCount < qshell.MIN_UPLOAD_THREAD_COUNT {
			upthreadCount = qshell.MIN_UPLOAD_THREAD_COUNT
		} else if upthreadCount > qshell.MAX_UPLOAD_THREAD_COUNT {
			upthreadCount = qshell.MAX_UPLOAD_THREAD_COUNT
		}
	}

	fileExporter := qshell.FileExporter{
		SuccessFname:   successFname,
		FailureFname:   failureFname,
		OverwriteFname: overwriteFname,
	}

	qshell.QiniuUpload(int(upthreadCount), &uploadConfig, &fileExporter)
}
