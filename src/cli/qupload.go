package cli

import (
	"encoding/json"

	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"os"
	"qshell"
	"strconv"
)

func QiniuUpload(cmd string, params ...string) {
	if len(params) == 1 || len(params) == 2 {
		var uploadConfigFile string
		var threadCount int64
		var err error
		if len(params) == 2 {
			threadCount, err = strconv.ParseInt(params[0], 10, 64)
			if err != nil {
				logs.Error("Invalid <ThreadCount> value,", params[0])
				return
			}
			uploadConfigFile = params[1]
		} else {
			uploadConfigFile = params[0]
		}

		//read upload config
		fp, err := os.Open(uploadConfigFile)
		if err != nil {
			logs.Error("Open upload config file `%s` error due to `%s`", uploadConfigFile, err)
			return
		}
		defer fp.Close()
		configData, err := ioutil.ReadAll(fp)
		if err != nil {
			logs.Error("Read upload config file `%s` error due to `%s`", uploadConfigFile, err)
			return
		}
		var uploadConfig qshell.UploadConfig
		err = json.Unmarshal(configData, &uploadConfig)
		if err != nil {
			logs.Error("Parse upload config file `%s` errror due to `%s`", uploadConfigFile, err)
			return
		}
		srcFileInfo, err := os.Stat(uploadConfig.SrcDir)

		if err != nil {
			logs.Error("Upload config error for parameter `SrcDir`,", err)
			return
		}

		if !srcFileInfo.IsDir() {
			logs.Error("Upload src dir should be a directory")
			return
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

		qshell.QiniuUpload(int(threadCount), &uploadConfig)
	} else {
		CmdHelp(cmd)
	}
}
