package cli

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"os"
	"qshell"
	"strconv"
)

func QiniuDownload(cmd string, params ...string) {
	if len(params) == 1 || len(params) == 2 {
		var threadCount int64 = 5
		var downloadConfigFile string
		var err error
		if len(params) == 1 {
			downloadConfigFile = params[0]
		} else {
			threadCount, err = strconv.ParseInt(params[0], 10, 64)
			if err != nil {
				logs.Error("Invalid value for <ThreadCount>", params[0])
				os.Exit(qshell.STATUS_HALT)
			}
			downloadConfigFile = params[1]
		}

		//read download config
		fp, err := os.Open(downloadConfigFile)
		if err != nil {
			logs.Error("Open download config file `%s` error, %s", downloadConfigFile, err)
			os.Exit(qshell.STATUS_HALT)
		}
		defer fp.Close()
		configData, err := ioutil.ReadAll(fp)
		if err != nil {
			logs.Error("Read download config file `%s` error, %s", downloadConfigFile, err)
			os.Exit(qshell.STATUS_HALT)
		}

		var downloadConfig qshell.DownloadConfig
		err = json.Unmarshal(configData, &downloadConfig)
		if err != nil {
			logs.Error("Parse download config file `%s` error, %s", downloadConfigFile, err)
			os.Exit(qshell.STATUS_HALT)
		}

		destFileInfo, err := os.Stat(downloadConfig.DestDir)

		if err != nil {
			logs.Error("Download config error for parameter `DestDir`,", err)
			os.Exit(qshell.STATUS_HALT)
		}

		if !destFileInfo.IsDir() {
			logs.Error("Download dest dir should be a directory")
			os.Exit(qshell.STATUS_HALT)
		}

		if threadCount < qshell.MIN_DOWNLOAD_THREAD_COUNT || threadCount > qshell.MAX_DOWNLOAD_THREAD_COUNT {
			logs.Info("Tip: you can set <ThreadCount> value between %d and %d to improve speed\n",
				qshell.MIN_DOWNLOAD_THREAD_COUNT, qshell.MAX_DOWNLOAD_THREAD_COUNT)

			if threadCount < qshell.MIN_DOWNLOAD_THREAD_COUNT {
				threadCount = qshell.MIN_DOWNLOAD_THREAD_COUNT
			} else if threadCount > qshell.MAX_DOWNLOAD_THREAD_COUNT {
				threadCount = qshell.MAX_DOWNLOAD_THREAD_COUNT
			}
		}

		qshell.QiniuDownload(int(threadCount), &downloadConfig)
	} else {
		CmdHelp(cmd)
	}
}
