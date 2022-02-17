package cmd

import (
	"github.com/qiniu/go-sdk/v7/conf"
	"github.com/qiniu/qshell/v2/docs"
	"github.com/qiniu/qshell/v2/iqshell/common/utils/operations"
	"github.com/spf13/cobra"
)

var rpcEncodeCmdBuilder = func() *cobra.Command {
	var info = operations.RpcInfo{}
	var cmd = &cobra.Command{
		Use:        "rpcencode <DataToEncode1> [<DataToEncode2> [...]]",
		Short:      "rpcencode of qiniu",
		Args:       cobra.MinimumNArgs(1),
		SuggestFor: []string{"rpc"},
		Run: func(cmd *cobra.Command, args []string) {
			info.Params = args
			prepare(cmd, nil)
			operations.RpcEncode(info)
		},
	}
	return cmd
}

var rpcDecodeCmdBuilder = func() *cobra.Command {
	var info = operations.RpcInfo{}
	var cmd = &cobra.Command{
		Use:        "rpcdecode [DataToEncode...]",
		Short:      "rpcdecode of qiniu",
		SuggestFor: []string{"rpc"},
		Run: func(cmd *cobra.Command, args []string) {
			info.Params = args
			prepare(cmd, nil)
			operations.RpcDecode(info)
		},
	}
	return cmd
}

var base64EncodeCmdBuilder = func() *cobra.Command {
	var info = operations.Base64Info{}
	var cmd = &cobra.Command{
		Use:        "b64encode <DataToEncode>",
		Short:      "Base64 Encode, default not url safe",
		Long:       "Base64 encode of data, url safe base64 is not turn on by default. Use -safe flag to turn it on",
		SuggestFor: []string{"b64"},
		Run: func(cmd *cobra.Command, args []string) {
			cmdId = docs.B64Encode
			if len(args) > 0 {
				info.Data = args[0]
			}
			prepare(cmd, &info)
			operations.Base64Encode(info)
		},
	}
	cmd.Flags().BoolVarP(&info.UrlSafe, "safe", "s", false, "use urlsafe base64 encode")
	return cmd
}

var base64DecodeCmdBuilder = func() *cobra.Command {
	var info = operations.Base64Info{}
	var cmd = &cobra.Command{
		Use:        "b64decode <DataToDecode>",
		Short:      "Base64 Decode, default nor url safe",
		Long:       "Base64 Decode of data, urlsafe base64 is not turn on by default. Use -safe flag to turn it on",
		SuggestFor: []string{"b64"},
		Run: func(cmd *cobra.Command, args []string) {
			cmdId = docs.B64Decode
			if len(args) > 0 {
				info.Data = args[0]
			}
			prepare(cmd, &info)
			operations.Base64Decode(info)
		},
	}
	cmd.Flags().BoolVarP(&info.UrlSafe, "safe", "s", false, "use urlsafe base64 decode")
	return cmd
}

var ts2dCmdCmdBuilder = func() *cobra.Command {
	var info = operations.TimestampInfo{}
	var cmd = &cobra.Command{
		Use:   "ts2d <TimestampInSeconds>",
		Short: "Convert timestamp in seconds to a date (TZ: Local)",
		Run: func(cmd *cobra.Command, args []string) {
			cmdId = docs.TS2dType
			if len(args) > 0 {
				info.Value = args[0]
			}
			prepare(cmd, &info)
			operations.Timestamp2Date(info)
		},
	}
	return cmd
}

var tms2dCmdBuilder = func() *cobra.Command {
	var info = operations.TimestampInfo{}
	var cmd = &cobra.Command{
		Use:   "tms2d <TimestampInMilliSeconds>",
		Short: "Convert timestamp in milliseconds to a date (TZ: Local)",
		Run: func(cmd *cobra.Command, args []string) {
			cmdId = docs.TMs2dType
			if len(args) > 0 {
				info.Value = args[0]
			}
			prepare(cmd, &info)
			operations.TimestampMilli2Date(info)
		},
	}
	return cmd
}

var tns2dCmdBuilder = func() *cobra.Command {
	var info = operations.TimestampInfo{}
	var cmd = &cobra.Command{
		Use:   "tns2d <TimestampInNanoSeconds>",
		Short: "Convert timestamp in Nanoseconds to a date (TZ: Local)",
		Run: func(cmd *cobra.Command, args []string) {
			cmdId = docs.TNs2dType
			if len(args) > 0 {
				info.Value = args[0]
			}
			prepare(cmd, &info)
			operations.TimestampNano2Date(info)
		},
	}
	return cmd
}

var d2tsCmdBuilder = func() *cobra.Command {
	var info = operations.TimestampInfo{}
	var cmd = &cobra.Command{
		Use:   "d2ts <SecondsToNow>",
		Short: "Create a timestamp in seconds using seconds to now",
		Run: func(cmd *cobra.Command, args []string) {
			cmdId = docs.DateToTimestampType
			if len(args) > 0 {
				info.Value = args[0]
			}
			prepare(cmd, &info)
			operations.Date2Timestamp(info)
		},
	}
	return cmd
}

var urlEncodeCmdBuilder = func() *cobra.Command {
	var info = operations.UrlInfo{}
	var cmd = &cobra.Command{
		Use:   "urlencode <DataToEncode>",
		Short: "Url Encode",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				info.Url = args[0]
			}
			prepare(cmd, nil)
			operations.UrlEncode(info)
		},
	}
	return cmd
}

var urlDecodeCmdBuilder = func() *cobra.Command {
	var info = operations.UrlInfo{}
	var cmd = &cobra.Command{
		Use:   "urldecode <DataToDecode>",
		Short: "Url Decode",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				info.Url = args[0]
			}
			prepare(cmd, nil)
			operations.UrlDecode(info)
		},
	}
	return cmd
}

var etagCmdBuilder = func() *cobra.Command {
	var info = operations.EtagInfo{}
	var cmd = &cobra.Command{
		Use:   "qetag <LocalFilePath>",
		Short: "Calculate the hash of local file using the algorithm of qiniu qetag",
		Run: func(cmd *cobra.Command, args []string) {
			cmdId = docs.QTagType
			if len(args) > 0 {
				info.FilePath = args[0]
			}
			prepare(cmd, &info)
			operations.CreateEtag(info)
		},
	}
	return cmd
}

var unzipCmdBuilder = func() *cobra.Command {
	var info = operations.ZipInfo{}
	var cmd = &cobra.Command{
		Use:   "unzip <QiniuZipFilePath> [<UnzipToDir>]",
		Short: "Unzip the archive file created by the qiniu mkzip API",
		Long:  "Unzip to current workding directory by default, use -dir to specify a directory",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			cmdId = docs.UnzipType
			if len(args) > 0 {
				info.ZipFilePath = args[0]
			}
			if len(args) > 1 {
				info.UnzipPath = args[1]
			}
			prepare(cmd, &info)
			operations.Unzip(info)
		},
	}
	cmd.Flags().StringVar(&info.UnzipPath, "dir", "", "unzip directory")
	return cmd
}

var reqIdCmdBuilder = func() *cobra.Command {
	var info = operations.ReqIdInfo{}
	var cmd = &cobra.Command{
		Use:   "reqid <ReqIdToDecode>",
		Short: "Decode qiniu reqid",
		Run: func(cmd *cobra.Command, args []string) {
			cmdId = docs.ReqIdType
			if len(args) > 0 {
				info.ReqId = args[0]
			}
			prepare(cmd, &info)
			operations.DecodeReqId(info)
		},
	}
	return cmd
}

var IpCmdBuilder = func() *cobra.Command {
	var info = operations.IpQueryInfo{}
	var cmd = &cobra.Command{
		Use:   "ip <Ip1> [<Ip2> [<Ip3> ...]]]",
		Short: "Query the ip information",
		Run: func(cmd *cobra.Command, args []string) {
			cmdId = docs.IPType
			info.Ips = args
			prepare(cmd, &info)
			operations.IpQuery(info)
		},
	}
	return cmd
}

var TokenCmdBuilder = func() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "token",
		Short: "Token related command",
		Long:  "Create QBox token, Qiniu token, Upload token and Download token",
		Run: func(cmd *cobra.Command, args []string) {
			cmdId = docs.TokenType
			prepare(cmd, nil)
		},
	}

	cmd.AddCommand(
		QboxTokenCmdBuilder(),
		QiniuTokenCmdBuilder(),
		UploadTokenCmdBuilder(),
	)
	return cmd
}

var QboxTokenCmdBuilder = func() *cobra.Command {
	var info = operations.QBoxTokenInfo{}
	var cmd = &cobra.Command{
		Use:   "qbox <Url>",
		Short: "Create QBox token",
		Run: func(cmd *cobra.Command, args []string) {
			cmdId = docs.TokenType
			if len(args) > 0 {
				info.Url = args[0]
			}
			prepare(cmd, &info)
			operations.CreateQBoxToken(info)
		},
	}

	cmd.Flags().StringVarP(&info.ContentType, "content-type", "t", conf.CONTENT_TYPE_JSON, "http request content type, application/json by default")
	cmd.Flags().StringVarP(&info.Body, "http-body", "b", "", "http request body, when content-type is application/x-www-form-urlencoded, http body must be specified")
	cmd.Flags().StringVarP(&info.AccessKey, "access-key", "a", "", "access key")
	cmd.Flags().StringVarP(&info.SecretKey, "secret-key", "s", "", "secret key")

	return cmd
}

var QiniuTokenCmdBuilder = func() *cobra.Command {
	var info = operations.QiniuTokenInfo{}
	var cmd = &cobra.Command{
		Use:   "qiniu <Url>",
		Short: "Create Qiniu Token",
		Run: func(cmd *cobra.Command, args []string) {
			cmdId = docs.TokenType
			if len(args) > 0 {
				info.Url = args[0]
			}
			prepare(cmd, &info)
			operations.CreateQiniuToken(info)
		},
	}

	cmd.Flags().StringVarP(&info.ContentType, "content-type", "t", conf.CONTENT_TYPE_JSON, "http request content type, application/json by default")
	cmd.Flags().StringVarP(&info.Body, "http-body", "b", "", "http request body, when content-type is application/x-www-form-urlencoded, http body must be specified")
	cmd.Flags().StringVarP(&info.AccessKey, "access-key", "a", "", "access key")
	cmd.Flags().StringVarP(&info.SecretKey, "secret-key", "s", "", "secret key")
	cmd.Flags().StringVarP(&info.Method, "method", "m", "GET", "http request method")

	return cmd
}

var UploadTokenCmdBuilder = func() *cobra.Command {
	var info = operations.UploadTokenInfo{}
	var cmd = &cobra.Command{
		Use:   "upload <PutPolicyJsonFile>",
		Short: "Create upload token using put policy",
		Run: func(cmd *cobra.Command, args []string) {
			cmdId = docs.TokenType
			if len(args) > 0 {
				info.Url = args[0]
			}
			prepare(cmd, &info)
			operations.CreateUploadToken(info)
		},
	}

	cmd.Flags().StringVarP(&info.AccessKey, "access-key", "a", "", "access key")
	cmd.Flags().StringVarP(&info.SecretKey, "secret-key", "s", "", "secret key")

	return cmd
}

var dirCacheCmdBuilder = func() *cobra.Command {
	var info = operations.DirCacheInfo{}
	var cmd = &cobra.Command{
		Use:   "dircache <DirCacheRootPath>",
		Short: "Cache the directory structure of a file path",
		Long:  "Cache the directory structure of a file path to a file, \nif <DirCacheResultFile> not specified, cache to stdout",
		Run: func(cmd *cobra.Command, args []string) {
			cmdId = docs.DirCacheType
			if len(args) > 0 {
				info.Dir = args[0]
			}
			prepare(cmd, &info)
			operations.DirCache(info)
		},
	}
	cmd.Flags().StringVarP(&info.SaveToFile, "outfile", "o", "", "output filepath")
	return cmd
}

func init() {
	rootCmd.AddCommand(
		rpcEncodeCmdBuilder(),
		rpcDecodeCmdBuilder(),
		base64EncodeCmdBuilder(),
		base64DecodeCmdBuilder(),
		ts2dCmdCmdBuilder(),
		tms2dCmdBuilder(),
		tns2dCmdBuilder(),
		d2tsCmdBuilder(),
		urlEncodeCmdBuilder(),
		urlDecodeCmdBuilder(),
		etagCmdBuilder(),
		unzipCmdBuilder(),
		reqIdCmdBuilder(),
		IpCmdBuilder(),
		TokenCmdBuilder(),
		dirCacheCmdBuilder(),
	)
}
