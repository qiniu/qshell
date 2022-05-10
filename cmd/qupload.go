package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/astaxie/beego/logs"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/spf13/cobra"
)

var qUploadCmd = &cobra.Command{
	Use:   "qupload <quploadConfigFile>",
	Short: "Batch upload files to the qiniu bucket",
	Args:  cobra.ExactArgs(1),
	Run:   QiniuUpload,
}

var (
	successFname   string
	failureFname   string
	overwriteFname string
	upthreadCount  int64
	uploadConfig   iqshell.UploadConfig
)

func init() {
	qUploadCmd.Flags().StringVarP(&successFname, "success-list", "s", "", "upload success (all) file list")
	qUploadCmd.Flags().StringVarP(&failureFname, "failure-list", "f", "", "upload failure file list")
	qUploadCmd.Flags().StringVarP(&overwriteFname, "overwrite-list", "w", "", "upload success (overwrite) file list")
	qUploadCmd.Flags().Int64VarP(&upthreadCount, "worker", "c", 1, "worker count")
	qUploadCmd.Flags().StringVarP(&callbackUrls, "callback-urls", "l", "", "upload callback urls, separated by comma")
	qUploadCmd.Flags().StringVarP(&callbackHost, "callback-host", "T", "", "upload callback host")
	RootCmd.AddCommand(qUploadCmd)
}

func parseUploadConfigFile(uploadConfigFile string, uploadConfig *iqshell.UploadConfig) (err error) {
	//read upload config
	if uploadConfigFile == "" {
		err = fmt.Errorf("config filename is empty")
		return
	}
	fp, oErr := os.Open(uploadConfigFile)
	if oErr != nil {
		err = fmt.Errorf("Open upload config file ``%s`: %v\n", uploadConfigFile, oErr)
		return
	}
	defer fp.Close()

	configData, rErr := ioutil.ReadAll(fp)
	if rErr != nil {
		err = fmt.Errorf("Read upload config file `%s`: %v\n", uploadConfigFile, rErr)
		return
	}
	//remove UTF-8 BOM
	configData = bytes.TrimPrefix(configData, []byte("\xef\xbb\xbf"))
	uErr := json.Unmarshal(configData, uploadConfig)
	if uErr != nil {
		err = fmt.Errorf("Parse upload config file `%s`: %v\n", uploadConfigFile, uErr)
		return
	}
	return
}

// [qupload]命令， 上传本地文件到七牛存储中
// 该命令会读取配置文件， 上传本地文件系统的文件到七牛存储中; 可以设置多线程上传，默认的线程区间在[iqshell.MIN_UPLOAD_THREAD_COUNT, iqshell.MAX_UPLOAD_THREAD_COUNT]
func QiniuUpload(cmd *cobra.Command, params []string) {

	configFile := params[0]

	pErr := parseUploadConfigFile(configFile, &uploadConfig)
	if pErr != nil {
		logs.Error(fmt.Sprintf("parse config file: %s: %v\n", configFile, pErr))
		os.Exit(iqshell.STATUS_HALT)
	}

	if uploadConfig.FileType != 1 && uploadConfig.FileType != 0 {
		logs.Error("Wrong Filetype, It should be 0 or 1 ")
		os.Exit(iqshell.STATUS_HALT)
	}

	srcFileInfo, err := os.Stat(uploadConfig.SrcDir)
	if err != nil {
		logs.Error("Upload config error for parameter `SrcDir`,", err)
		os.Exit(iqshell.STATUS_HALT)
	}

	if !srcFileInfo.IsDir() {
		logs.Error("Upload src dir should be a directory")
		os.Exit(iqshell.STATUS_HALT)
	}
	policy := storage.PutPolicy{}

	if (callbackUrls == "" && callbackHost != "") || (callbackUrls != "" && callbackHost == "") {
		fmt.Fprintf(os.Stderr, "callbackUrls and callback must exist at the same time\n")
		os.Exit(1)
	}
	if (uploadConfig.CallbackUrls == "" && uploadConfig.CallbackHost != "") || (uploadConfig.CallbackUrls != "" && uploadConfig.CallbackHost == "") {
		fmt.Fprintf(os.Stderr, "callbackUrls and callback must exist at the same time\n")
		os.Exit(1)
	}
	if (callbackHost != "" && callbackUrls != "") || (uploadConfig.CallbackHost != "" && uploadConfig.CallbackUrls != "") {
		callbackUrls = strings.Replace(callbackUrls, ",", ";", -1)
		policy.CallbackHost = callbackHost
		policy.CallbackURL = callbackUrls
		policy.CallbackBody = "key=$(key)&hash=$(etag)"
		policy.CallbackBodyType = "application/x-www-form-urlencoded"
	}
	uploadConfig.PutPolicy = policy

	//upload
	if upthreadCount < iqshell.MIN_UPLOAD_THREAD_COUNT || upthreadCount > iqshell.MAX_UPLOAD_THREAD_COUNT {
		logs.Info("Tip: you can set <ThreadCount> value between %d and %d to improve speed\n",
			iqshell.MIN_UPLOAD_THREAD_COUNT, iqshell.MAX_UPLOAD_THREAD_COUNT)

		if upthreadCount < iqshell.MIN_UPLOAD_THREAD_COUNT {
			upthreadCount = iqshell.MIN_UPLOAD_THREAD_COUNT
		} else if upthreadCount > iqshell.MAX_UPLOAD_THREAD_COUNT {
			upthreadCount = iqshell.MAX_UPLOAD_THREAD_COUNT
		}
	}

	fileExporter, fErr := iqshell.NewFileExporter(successFname, failureFname, overwriteFname)
	if fErr != nil {
		logs.Error("initialize fileExporter: ", fErr)
		os.Exit(iqshell.STATUS_HALT)
	}
	iqshell.QiniuUpload(int(upthreadCount), &uploadConfig, fileExporter)
}
