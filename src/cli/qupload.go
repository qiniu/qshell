package cli

import (
	"encoding/json"
	"flag"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"os"
	"qshell"
	"strconv"
)

func QiniuUpload(cmd string, params ...string) {
	var successFname string
	var failureFname string
	var overwriteFname string
	flagSet := flag.NewFlagSet("qupload", flag.ExitOnError)
	flagSet.StringVar(&successFname, "success-list", "", "upload success (all) file list")
	flagSet.StringVar(&failureFname, "failure-list", "", "upload failure file list")
	flagSet.StringVar(&overwriteFname, "overwrite-list", "", "upload success (overwrite) file list")
	flagSet.Parse(params)
	cmdParams := flagSet.Args()
	if len(cmdParams) == 1 || len(cmdParams) == 2 {
		var uploadConfigFile string
		var threadCount int64
		var err error
		if len(cmdParams) == 2 {
			threadCount, err = strconv.ParseInt(cmdParams[0], 10, 64)
			if err != nil {
				logs.Error("Invalid <ThreadCount> value,", cmdParams[0])
				os.Exit(2)
			}
			uploadConfigFile = cmdParams[1]
		} else {
			uploadConfigFile = cmdParams[0]
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
		if threadCount < qshell.MIN_UPLOAD_THREAD_COUNT || threadCount > qshell.MAX_UPLOAD_THREAD_COUNT {
			logs.Info("Tip: you can set <ThreadCount> value between %d and %d to improve speed\n",
				qshell.MIN_UPLOAD_THREAD_COUNT, qshell.MAX_UPLOAD_THREAD_COUNT)

			if threadCount < qshell.MIN_UPLOAD_THREAD_COUNT {
				threadCount = qshell.MIN_UPLOAD_THREAD_COUNT
			} else if threadCount > qshell.MAX_UPLOAD_THREAD_COUNT {
				threadCount = qshell.MAX_UPLOAD_THREAD_COUNT
			}
		}

		fileExporter := qshell.FileExporter{
			SuccessFname:   successFname,
			FailureFname:   failureFname,
			OverwriteFname: overwriteFname,
		}
		qshell.QiniuUpload(int(threadCount), &uploadConfig, &fileExporter)
	} else {
		CmdHelp(cmd)
	}
}
