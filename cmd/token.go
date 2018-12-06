package cmd

import (
	"fmt"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/conf"
	"github.com/qiniu/qshell/iqshell"
	"github.com/spf13/cobra"
	"io"
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
	cmdToken.AddCommand(cmdTokenQbox, cmdTokenQiniu)
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
