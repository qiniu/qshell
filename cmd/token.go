package cmd

import (
	"fmt"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/conf"
	"github.com/qiniu/qshell/iqshell"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"strings"
)

var (
	contentType string
	httpBody    string
	ak          string
	sk          string
)

var (
	cmdToken = &cobra.Command{
		Use:   "token [-a <AccessKey> -s <SecretKey>]",
		Short: "Token related command",
	}
	cmdTokenQbox = &cobra.Command{
		Use:   "qbox <Url>",
		Short: "Sign QBox administration token using url, http body, contentType",
		Args:  cobra.ExactArgs(1),
		Run:   QBoxToken,
	}
)

func init() {
	cmdTokenQbox.Flags().StringVarP(&contentType, "content-type", "t", conf.CONTENT_TYPE_JSON, "http request content type, application/json by default")
	cmdTokenQbox.Flags().StringVarP(&httpBody, "http-body", "b", "", "http request body, when content-type is application/x-www-form-urlencoded, http body must be specified")

	cmdToken.Flags().StringVarP(&ak, "access-key", "a", "", "access key")
	cmdToken.Flags().StringVarP(&sk, "secret-key", "s", "", "secret key")
	RootCmd.AddCommand(cmdToken)
	cmdToken.AddCommand(cmdTokenQbox)
}

func QBoxToken(cmd *cobra.Command, args []string) {
	var mac *qbox.Mac
	var mErr error
	if ak != "" && sk != "" {
		mac = qbox.NewMac(ak, sk)
	} else {
		mac, mErr = iqshell.GetMac()
		if mErr != nil {
			fmt.Fprintf(os.Stderr, "get mac: %v\n", mErr)
			os.Exit(1)
		}
	}
	req, rErr := http.NewRequest("", args[0], strings.NewReader(httpBody))
	if rErr != nil {
		fmt.Fprintf(os.Stderr, "create request: %v\n", rErr)
		os.Exit(1)
	}

	headers := http.Header{}
	headers.Add("Content-Type", contentType)
	req.Header = headers

	token, signErr := mac.SignRequest(req)
	if signErr != nil {
		fmt.Fprintf(os.Stderr, "sign url: %v\n", signErr)
		os.Exit(1)
	}
	fmt.Println("QBox " + token)
}
