package cli

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"qiniu/log"
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
				log.Error("Invalid value for <ThreadCount>", params[0])
				return
			}
			downloadConfigFile = params[1]
		}

		//read download config
		fp, err := os.Open(downloadConfigFile)
		if err != nil {
			log.Errorf("Open download config file `%s` error, %s", downloadConfigFile, err)
			return
		}
		defer fp.Close()
		configData, err := ioutil.ReadAll(fp)
		if err != nil {
			log.Errorf("Read download config file `%s` error, %s", downloadConfigFile, err)
			return
		}

		var downloadConfig qshell.DownloadConfig
		err = json.Unmarshal(configData, &downloadConfig)
		if err != nil {
			log.Errorf("Parse download config file `%s` error, %s", downloadConfigFile, err)
			return
		}

		destFileInfo, err := os.Stat(downloadConfig.DestDir)

		if err != nil {
			log.Error("Download config error for parameter `DestDir`,", err)
			return
		}

		if !destFileInfo.IsDir() {
			log.Error("Download dest dir should be a directory")
			return
		}

		if threadCount < qshell.MIN_DOWNLOAD_THREAD_COUNT || threadCount > qshell.MAX_DOWNLOAD_THREAD_COUNT {
			log.Infof("Tip: you can set <ThreadCount> value between %d and %d to improve speed\n",
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
