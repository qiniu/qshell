package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/conf"
	"github.com/qiniu/api.v7/storage"
	"github.com/qiniu/qshell/iqshell"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var (
	contentType string
	httpBody    string
	ak          string
	sk          string
	method      string
)

var (
	cmdToken = &cobra.Command{
		Use:   "token",
		Short: "Token related command",
		Long:  "Create QBox token, Qiniu token, Upload token and Download token",
	}
	cmdTokenQbox = &cobra.Command{
		Use:   "qbox <Url>",
		Short: "Create QBox token",
		Args:  cobra.ExactArgs(1),
		Run:   QBoxToken,
	}
	cmdTokenQiniu = &cobra.Command{
		Use:   "qiniu <Url>",
		Short: "Create Qiniu Token",
		Args:  cobra.ExactArgs(1),
		Run:   QiniuToken,
	}
	cmdTokenUpload = &cobra.Command{
		Use:   "upload <PutPolicyJsonFile>",
		Short: "Create upload token using put policy",
		Args:  cobra.ExactArgs(1),
		Run:   UploadToken,
	}
	cmdTokenDownload = &cobra.Command{}
)

func init() {
	cmdTokenQbox.Flags().StringVarP(&contentType, "content-type", "t", conf.CONTENT_TYPE_JSON, "http request content type, application/json by default")
	cmdTokenQbox.Flags().StringVarP(&httpBody, "http-body", "b", "", "http request body, when content-type is application/x-www-form-urlencoded, http body must be specified")
	cmdTokenQbox.Flags().StringVarP(&ak, "access-key", "a", "", "access key")
	cmdTokenQbox.Flags().StringVarP(&sk, "secret-key", "s", "", "secret key")

	cmdTokenQiniu.Flags().StringVarP(&contentType, "content-type", "t", conf.CONTENT_TYPE_JSON, "http request content type, application/json by default")
	cmdTokenQiniu.Flags().StringVarP(&httpBody, "http-body", "b", "", "http request body, when content-type is application/x-www-form-urlencoded, http body must be specified")
	cmdTokenQiniu.Flags().StringVarP(&ak, "access-key", "a", "", "access key")
	cmdTokenQiniu.Flags().StringVarP(&sk, "secret-key", "s", "", "secret key")
	cmdTokenQiniu.Flags().StringVarP(&method, "method", "m", "GET", "http request method")

	RootCmd.AddCommand(cmdToken)
	cmdToken.AddCommand(cmdTokenQbox, cmdTokenQiniu, cmdTokenUpload, cmdTokenDownload)
}

func macRequest(ak, sk, url, body, contentType, method string) (mac *qbox.Mac, req *http.Request, err error) {

	var mErr error
	var rErr error

	if ak != "" && sk != "" {
		mac = qbox.NewMac(ak, sk)
	} else {
		mac, mErr = iqshell.GetMac()
		if mErr != nil {
			err = fmt.Errorf("get mac: %v\n", mErr)
			return
		}
	}
	var httpBody io.Reader
	if body == "" {
		httpBody = nil
	} else {
		httpBody = strings.NewReader(body)
	}
	req, rErr = http.NewRequest(method, url, httpBody)
	if rErr != nil {
		err = fmt.Errorf("create request: %v\n", rErr)
		return
	}
	headers := http.Header{}
	headers.Add("Content-Type", contentType)
	req.Header = headers

	return
}

func QBoxToken(cmd *cobra.Command, args []string) {
	mac, req, mErr := macRequest(ak, sk, args[0], httpBody, contentType, "")
	if mErr != nil {
		fmt.Fprintf(os.Stderr, "create mac and reqeust: %v\n", mErr)
		os.Exit(1)
	}
	token, signErr := mac.SignRequest(req)
	if signErr != nil {
		fmt.Fprintf(os.Stderr, "create qbox token for url: %s: %v\n", args[0], signErr)
		os.Exit(1)
	}
	fmt.Println("QBox " + token)
}

func QiniuToken(cmd *cobra.Command, args []string) {

	mac, req, mErr := macRequest(ak, sk, args[0], httpBody, contentType, method)
	if mErr != nil {
		fmt.Fprintf(os.Stderr, "create mac and reqeust: %v\n", mErr)
		os.Exit(1)
	}
	token, signErr := mac.SignRequestV2(req)
	if signErr != nil {
		fmt.Fprintf(os.Stderr, "create qiniu token for url: %s: %v\n", args[0], signErr)
		os.Exit(1)
	}
	fmt.Println("Qiniu " + token)
}

func UploadToken(cmd *cobra.Command, args []string) {
	fileName := args[0]

	fileInfo, oErr := os.Open(fileName)
	if oErr != nil {
		fmt.Fprintf(os.Stderr, "open file %s: %v\n", fileName, oErr)
		os.Exit(1)
	}
	defer fileInfo.Close()

	configData, rErr := ioutil.ReadAll(fileInfo)
	if rErr != nil {
		fmt.Fprintf(os.Stderr, "read putPolicy config file `%s`: %v\n", fileName, rErr)
		os.Exit(1)
	}
	//remove UTF-8 BOM
	configData = bytes.TrimPrefix(configData, []byte("\xef\xbb\xbf"))

	putPolicy := new(storage.PutPolicy)
	uErr := json.Unmarshal(configData, putPolicy)
	if uErr != nil {
		fmt.Fprintf(os.Stderr, "parse upload config file `%s`: %v\n", fileName, uErr)
		os.Exit(1)
	}
	mac, mErr := iqshell.GetMac()
	if mErr != nil {
		fmt.Fprintf(os.Stderr, "get mac: %v\n", mErr)
		os.Exit(1)
	}
	uploadToken := putPolicy.UploadToken(mac)
	fmt.Println("UpToken " + uploadToken)
}

func DownloadToken(cmd *cobra.Command, args []string) {

}
