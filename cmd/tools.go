package cmd

import (
	"github.com/qiniu/go-sdk/v7/conf"
	"github.com/qiniu/qshell/v2/docs"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/utils/operations"
	"github.com/spf13/cobra"
)

var rpcEncodeCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.RpcInfo{}
	cmd := &cobra.Command{
		Use:        "rpcencode <DataToEncode1> [<DataToEncode2> [...]]",
		Short:      "RPC Encode of qiniu",
		SuggestFor: []string{"rpc"},
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.RpcEncodeType
			info.Params = args
			operations.RpcEncode(cfg, info)
		},
	}
	return cmd
}

var rpcDecodeCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.RpcInfo{}
	cmd := &cobra.Command{
		Use:        "rpcdecode [DataToEncode...]",
		Short:      "RPC Decode of qiniu",
		SuggestFor: []string{"rpc"},
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.RpcDecodeType
			info.Params = args
			operations.RpcDecode(cfg, info)
		},
	}
	return cmd
}

var base64EncodeCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.Base64Info{}
	cmd := &cobra.Command{
		Use:        "b64encode <DataToEncode>",
		Short:      "Base64 Encode, default not url safe",
		Long:       "Base64 encode of data, url safe base64 is not turn on by default. Use --safe flag to turn it on",
		SuggestFor: []string{"b64"},
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.B64Encode
			if len(args) > 0 {
				info.Data = args[0]
			}
			operations.Base64Encode(cfg, info)
		},
	}
	cmd.Flags().BoolVarP(&info.UrlSafe, "safe", "s", false, "use urlsafe base64 encode")
	return cmd
}

var base64DecodeCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.Base64Info{}
	cmd := &cobra.Command{
		Use:        "b64decode <DataToDecode>",
		Short:      "Base64 Decode, default nor url safe",
		Long:       "Base64 Decode of data, urlsafe base64 is not turn on by default. Use --safe flag to turn it on",
		SuggestFor: []string{"b64"},
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.B64Decode
			if len(args) > 0 {
				info.Data = args[0]
			}
			operations.Base64Decode(cfg, info)
		},
	}
	cmd.Flags().BoolVarP(&info.UrlSafe, "safe", "s", false, "use urlsafe base64 decode")
	return cmd
}

var ts2dCmdCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.TimestampInfo{}
	cmd := &cobra.Command{
		Use:   "ts2d <TimestampInSeconds>",
		Short: "Convert timestamp in seconds to a date (TZ: Local)",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.TS2dType
			if len(args) > 0 {
				info.Value = args[0]
			}
			operations.Timestamp2Date(cfg, info)
		},
	}
	return cmd
}

var tms2dCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.TimestampInfo{}
	cmd := &cobra.Command{
		Use:   "tms2d <TimestampInMilliSeconds>",
		Short: "Convert timestamp in milliseconds to a date (TZ: Local)",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.TMs2dType
			if len(args) > 0 {
				info.Value = args[0]
			}
			operations.TimestampMilli2Date(cfg, info)
		},
	}
	return cmd
}

var tns2dCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.TimestampInfo{}
	cmd := &cobra.Command{
		Use:   "tns2d <TimestampInNanoSeconds>",
		Short: "Convert timestamp in Nanoseconds to a date (TZ: Local)",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.TNs2dType
			if len(args) > 0 {
				info.Value = args[0]
			}
			operations.TimestampNano2Date(cfg, info)
		},
	}
	return cmd
}

var d2tsCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.TimestampInfo{}
	cmd := &cobra.Command{
		Use:   "d2ts <SecondsToNow>",
		Short: "Create a timestamp of some seconds later, unit of timestamp is second",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.DateToTimestampType
			if len(args) > 0 {
				info.Value = args[0]
			}
			operations.Date2Timestamp(cfg, info)
		},
	}
	return cmd
}

var urlEncodeCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.UrlInfo{}
	cmd := &cobra.Command{
		Use:   "urlencode <DataToEncode>",
		Short: "Url Encode",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.UrlEncodeType
			if len(args) > 0 {
				info.Url = args[0]
			}
			operations.UrlEncode(cfg, info)
		},
	}
	return cmd
}

var urlDecodeCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.UrlInfo{}
	cmd := &cobra.Command{
		Use:   "urldecode <DataToDecode>",
		Short: "Url Decode",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.UrlDecodeType
			if len(args) > 0 {
				info.Url = args[0]
			}
			operations.UrlDecode(cfg, info)
		},
	}
	return cmd
}

var etagCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.EtagInfo{}
	cmd := &cobra.Command{
		Use:   "qetag <LocalFilePath>",
		Short: "Calculate the hash of local file using the algorithm of qiniu qetag",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.QTagType
			if len(args) > 0 {
				info.FilePath = args[0]
			}
			operations.CreateEtag(cfg, info)
		},
	}
	return cmd
}

var unzipCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.ZipInfo{}
	cmd := &cobra.Command{
		Use:   "unzip <QiniuZipFilePath> [<UnzipToDir>]",
		Short: "Unzip the archive file created by the qiniu mkzip API",
		Long:  "Unzip to current workding directory by default, use -dir to specify a directory",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.UnzipType
			if len(args) > 0 {
				info.ZipFilePath = args[0]
			}
			if len(args) > 1 {
				info.UnzipPath = args[1]
			}
			operations.Unzip(cfg, info)
		},
	}
	cmd.Flags().StringVar(&info.UnzipPath, "dir", "", "unzip directory")
	return cmd
}

var reqIdCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.ReqIdInfo{}
	cmd := &cobra.Command{
		Use:   "reqid <ReqIdToDecode>",
		Short: "Decode qiniu reqid",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.ReqIdType
			if len(args) > 0 {
				info.ReqId = args[0]
			}
			operations.DecodeReqId(cfg, info)
		},
	}
	return cmd
}

var IpCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.IpQueryInfo{}
	cmd := &cobra.Command{
		Use:   "ip <Ip1> [<Ip2> [<Ip3> ...]]]",
		Short: "Query the ip information",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.IPType
			info.Ips = args
			operations.IpQuery(cfg, info)
		},
	}
	return cmd
}

var TokenCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token",
		Short: "Token related command",
		Long:  "Create QBox token, Qiniu token and Upload token",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.TokenType
			operations.Token(cfg)
		},
	}

	cmd.AddCommand(
		QboxTokenCmdBuilder(cfg),
		QiniuTokenCmdBuilder(cfg),
		UploadTokenCmdBuilder(cfg),
	)
	return cmd
}

var QboxTokenCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.QBoxTokenInfo{}
	cmd := &cobra.Command{
		Use:   "qbox <Url>",
		Short: "Create QBox token",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.TokenType
			if len(args) > 0 {
				info.Url = args[0]
			}
			operations.CreateQBoxToken(cfg, info)
		},
	}

	cmd.Flags().StringVarP(&info.ContentType, "content-type", "t", conf.CONTENT_TYPE_JSON, "http request content type, application/json by default")
	cmd.Flags().StringVarP(&info.Body, "http-body", "b", "", "http request body, when content-type is application/x-www-form-urlencoded, http body must be specified")
	cmd.Flags().StringVarP(&info.AccessKey, "access-key", "a", "", "access key")
	cmd.Flags().StringVarP(&info.SecretKey, "secret-key", "s", "", "secret key")

	return cmd
}

var QiniuTokenCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.QiniuTokenInfo{}
	cmd := &cobra.Command{
		Use:   "qiniu <Url>",
		Short: "Create Qiniu Token",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.TokenType
			if len(args) > 0 {
				info.Url = args[0]
			}
			operations.CreateQiniuToken(cfg, info)
		},
	}

	cmd.Flags().StringVarP(&info.ContentType, "content-type", "t", conf.CONTENT_TYPE_JSON, "http request content type, application/json by default")
	cmd.Flags().StringVarP(&info.Body, "http-body", "b", "", "http request body, when content-type is application/x-www-form-urlencoded, http body must be specified")
	cmd.Flags().StringVarP(&info.AccessKey, "access-key", "a", "", "access key")
	cmd.Flags().StringVarP(&info.SecretKey, "secret-key", "s", "", "secret key")
	cmd.Flags().StringVarP(&info.Method, "method", "m", "GET", "http request method")

	return cmd
}

var UploadTokenCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.UploadTokenInfo{}
	cmd := &cobra.Command{
		Use:   "upload <PutPolicyJsonFile>",
		Short: "Create upload token using put policy",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.TokenType
			if len(args) > 0 {
				info.PutPolicyFilePath = args[0]
			}
			operations.CreateUploadToken(cfg, info)
		},
	}

	cmd.Flags().StringVarP(&info.AccessKey, "access-key", "a", "", "access key")
	cmd.Flags().StringVarP(&info.SecretKey, "secret-key", "s", "", "secret key")

	return cmd
}

var dirCacheCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.DirCacheInfo{}
	cmd := &cobra.Command{
		Use:   "dircache <DirCacheRootPath>",
		Short: "Cache the directory structure of a file path",
		Long:  "Cache the directory structure of a file path to a file, \nif <DirCacheResultFile> not specified, cache to stdout",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.DirCacheType
			if len(args) > 0 {
				info.Dir = args[0]
			}
			operations.DirCache(cfg, info)
		},
	}
	cmd.Flags().StringVarP(&info.SaveToFile, "outfile", "o", "", "output filepath")
	return cmd
}

var funcCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.FuncCallInfo{}
	cmd := &cobra.Command{
		Use:   "func <ParamsJson> <FuncTemplate>",
		Short: "qshell custom function call",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.FuncType
			if len(args) > 0 {
				info.ParamsJson = args[0]
			}
			if len(args) > 1 {
				info.FuncTemplate = args[1]
			}
			operations.FuncCall(cfg, info)
		},
	}
	return cmd
}

func init() {
	registerLoader(toolsCmdLoader)
}

func toolsCmdLoader(superCmd *cobra.Command, cfg *iqshell.Config) {
	superCmd.AddCommand(
		rpcEncodeCmdBuilder(cfg),
		rpcDecodeCmdBuilder(cfg),
		base64EncodeCmdBuilder(cfg),
		base64DecodeCmdBuilder(cfg),
		ts2dCmdCmdBuilder(cfg),
		tms2dCmdBuilder(cfg),
		tns2dCmdBuilder(cfg),
		d2tsCmdBuilder(cfg),
		urlEncodeCmdBuilder(cfg),
		urlDecodeCmdBuilder(cfg),
		etagCmdBuilder(cfg),
		unzipCmdBuilder(cfg),
		reqIdCmdBuilder(cfg),
		IpCmdBuilder(cfg),
		TokenCmdBuilder(cfg),
		dirCacheCmdBuilder(cfg),
		funcCmdBuilder(cfg),
	)
}
