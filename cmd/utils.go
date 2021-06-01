package cmd

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/spf13/cobra"
)

const (
	// ASCII英文字母
	ALPHA_LIST = "abcdefghijklmnopqrstuvwxyz"
)

const (
	KB = 1024
	MB = 1024 * KB
	GB = 1024 * MB
	TB = 1024 * GB
)

var (
	urlSafe  bool
	unzipDir string
)

var (
	rpcECmd = &cobra.Command{
		Use:        "rpcencode <DataToEncode1> [<DataToEncode2> [...]]",
		Short:      "rpcencode of qiniu",
		Args:       cobra.MinimumNArgs(1),
		SuggestFor: []string{"rpc"},
		Run:        RpcEncode,
	}
	rpcDCmd = &cobra.Command{
		Use:        "rpcdecode [DataToEncode...]",
		Short:      "rpcdecode of qiniu",
		SuggestFor: []string{"rpc"},
		Run:        RpcDecode,
	}
	b64ECmd = &cobra.Command{
		Use:        "b64encode <DataToEncode>",
		Short:      "Base64 Encode, default not url safe",
		Long:       "Base64 encode of data, url safe base64 is not turn on by default. Use -safe flag to turn it on",
		Args:       cobra.MinimumNArgs(1),
		SuggestFor: []string{"b64"},
		Run:        Base64Encode,
	}
	b64DCmd = &cobra.Command{
		Use:        "b64decode <DataToDecode>",
		Short:      "Base64 Decode, default nor url safe",
		Long:       "Base64 Decode of data, urlsafe base64 is not turn on by default. Use -safe flag to turn it on",
		Args:       cobra.MinimumNArgs(1),
		SuggestFor: []string{"b64"},
		Run:        Base64Decode,
	}
	ts2dCmd = &cobra.Command{
		Use:   "ts2d <TimestampInSeconds>",
		Short: "Convert timestamp in seconds to a date (TZ: Local)",
		Args:  cobra.ExactArgs(1),
		Run:   Timestamp2Date,
	}
	tms2dCmd = &cobra.Command{
		Use:   "tms2d <TimestampInMilliSeconds>",
		Short: "Convert timestamp in milliseconds to a date (TZ: Local)",
		Args:  cobra.ExactArgs(1),
		Run:   TimestampMilli2Date,
	}
	tns2dCmd = &cobra.Command{
		Use:   "tns2d <TimestampInNanoSeconds>",
		Short: "Convert timestamp in Nanoseconds to a date (TZ: Local)",
		Args:  cobra.ExactArgs(1),
		Run:   TimestampNano2Date,
	}
	d2tsCmd = &cobra.Command{
		Use:   "d2ts <SecondsToNow>",
		Short: "Create a timestamp in seconds using seconds to now",
		Args:  cobra.ExactArgs(1),
		Run:   Date2Timestamp,
	}
	urlEcmd = &cobra.Command{
		Use:   "urlencode <DataToEncode>",
		Short: "Url Encode",
		Args:  cobra.ExactArgs(1),
		Run:   Urlencode,
	}
	urlDcmd = &cobra.Command{
		Use:   "urldecode <DataToDecode>",
		Short: "Url Decode",
		Args:  cobra.ExactArgs(1),
		Run:   Urldecode,
	}
	etagCmd = &cobra.Command{
		Use:   "qetag <LocalFilePath>",
		Short: "Calculate the hash of local file using the algorithm of qiniu qetag",
		Args:  cobra.ExactArgs(1),
		Run:   Qetag,
	}
	unzipCmd = &cobra.Command{
		Use:   "unzip <QiniuZipFilePath> [<UnzipToDir>",
		Short: "Unzip the archive file created by the qiniu mkzip API",
		Long:  "Unzip to current workding directory by default, use -dir to specify a directory",
		Args:  cobra.ExactArgs(1),
		Run:   Unzip,
	}
	reqidCmd = &cobra.Command{
		Use:   "reqid <ReqIdToDecode>",
		Short: "Decode qiniu reqid",
		Args:  cobra.ExactArgs(1),
		Run:   ReqId,
	}
)

func init() {
	b64ECmd.Flags().BoolVarP(&urlSafe, "safe", "s", false, "use urlsafe base64 encode")
	b64DCmd.Flags().BoolVarP(&urlSafe, "safe", "s", false, "use urlsafe base64 decode")
	unzipCmd.Flags().StringVar(&unzipDir, "dir", "", "unzip directory")

	RootCmd.AddCommand(rpcECmd, rpcDCmd, b64ECmd, b64DCmd, ts2dCmd, tms2dCmd, tns2dCmd,
		d2tsCmd, urlEcmd, urlDcmd, etagCmd, unzipCmd, reqidCmd)
}

// 转化文件大小到人工可读的字符串，以相应的单位显示
func FormatFsize(fsize int64) (result string) {
	if fsize > TB {
		result = fmt.Sprintf("%.2f TB", float64(fsize)/float64(TB))
	} else if fsize > GB {
		result = fmt.Sprintf("%.2f GB", float64(fsize)/float64(GB))
	} else if fsize > MB {
		result = fmt.Sprintf("%.2f MB", float64(fsize)/float64(MB))
	} else if fsize > KB {
		result = fmt.Sprintf("%.2f KB", float64(fsize)/float64(KB))
	} else {
		result = fmt.Sprintf("%d B", fsize)
	}

	return
}

func RpcDecode(cmd *cobra.Command, params []string) {
	if len(params) > 0 {
		for _, param := range params {
			decodedStr, _ := iqshell.Decode(param)
			fmt.Println(decodedStr)
		}
	} else {
		bScanner := bufio.NewScanner(os.Stdin)
		for bScanner.Scan() {
			toDecode := bScanner.Text()
			decodedStr, _ := iqshell.Decode(string(toDecode))
			fmt.Println(decodedStr)
		}
	}
}

func RpcEncode(cmd *cobra.Command, params []string) {
	for _, param := range params {
		encodedStr := iqshell.Encode(param)
		fmt.Println(encodedStr)
	}
}

// base64编码数据
func Base64Encode(cmd *cobra.Command, params []string) {
	dataToEncode := params[0]
	dataEncoded := ""
	if urlSafe {
		dataEncoded = base64.URLEncoding.EncodeToString([]byte(dataToEncode))
	} else {
		dataEncoded = base64.StdEncoding.EncodeToString([]byte(dataToEncode))
	}
	fmt.Println(dataEncoded)
}

// 解码base64编码的数据
func Base64Decode(cmd *cobra.Command, params []string) {
	var err error

	dataToDecode := params[0]
	var dataDecoded []byte
	if urlSafe {
		dataDecoded, err = base64.URLEncoding.DecodeString(dataToDecode)
		if err != nil {
			logs.Error("Failed to decode `", dataToDecode, "' in url safe mode.")
			return
		}
	} else {
		dataDecoded, err = base64.StdEncoding.DecodeString(dataToDecode)
		if err != nil {
			logs.Error("Failed to decode `", dataToDecode, "' in standard mode.")
			return
		}
	}
	fmt.Println(string(dataDecoded))
}

// 转化unix时间戳为可读的字符串
func Timestamp2Date(cmd *cobra.Command, params []string) {
	ts, err := strconv.ParseInt(params[0], 10, 64)
	if err != nil {
		logs.Error("Invalid timestamp value,", params[0])
		return
	}
	t := time.Unix(ts, 0)
	fmt.Println(t.String())
}

// 转化纳秒时间戳到人工可读的字符串
func TimestampNano2Date(cmd *cobra.Command, params []string) {
	tns, err := strconv.ParseInt(params[0], 10, 64)
	if err != nil {
		logs.Error("Invalid nano timestamp value,", params[0])
		return
	}
	t := time.Unix(0, tns*100)
	fmt.Println(t.String())
}

// 转化毫秒时间戳到人工可读的字符串
func TimestampMilli2Date(cmd *cobra.Command, params []string) {
	tms, err := strconv.ParseInt(params[0], 10, 64)
	if err != nil {
		logs.Error("Invalid milli timestamp value,", params[0])
		return
	}
	t := time.Unix(tms/1000, 0)
	fmt.Println(t.String())
}

// 转化时间字符串到unix时间戳
func Date2Timestamp(cmd *cobra.Command, params []string) {
	secs, err := strconv.ParseInt(params[0], 10, 64)
	if err != nil {
		logs.Error("Invalid seconds to now,", params[0])
		return
	}
	t := time.Now()
	t = t.Add(time.Second * time.Duration(secs))
	fmt.Println(t.Unix())
}

func Urlencode(cmd *cobra.Command, params []string) {
	dataToEncode := params[0]
	_, err := url.Parse(dataToEncode)
	if err != nil {
		dataEncoded := url.PathEscape(dataToEncode)
		fmt.Println(dataEncoded)
	} else {
		dataEncoded := url.PathEscape(dataToEncode)
		fmt.Println(dataEncoded)
	}
}

func Urldecode(cmd *cobra.Command, params []string) {
	dataToDecode := params[0]
	dataDecoded, err := url.PathUnescape(dataToDecode)
	if err != nil {
		logs.Error("Failed to unescape data `", dataToDecode, "'")
	} else {
		fmt.Println(dataDecoded)
	}
}

// 计算文件的hash值，使用七牛的etag算法
func Qetag(cmd *cobra.Command, params []string) {
	localFilePath := params[0]
	qetag, err := iqshell.GetEtag(localFilePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(qetag)
}

// 解压使用mkzip压缩的文件
func Unzip(cmd *cobra.Command, params []string) {
	zipFilePath := params[0]
	var err error
	if unzipDir == "" {
		unzipDir, err = os.Getwd()
		if err != nil {
			logs.Error("Get current work directory failed due to error", err)
			return
		}
	} else {
		if _, statErr := os.Stat(unzipDir); statErr != nil {
			logs.Error("Specified <UnzipToDir> is not a valid directory")
			return
		}
	}
	unzipErr := iqshell.Unzip(zipFilePath, unzipDir)
	if unzipErr != nil {
		logs.Error("Unzip file failed due to error", unzipErr)
	}
}

// 解析reqid， 打印人工可读的字符串
func ReqId(cmd *cobra.Command, params []string) {
	reqId := params[0]
	decodedBytes, err := base64.URLEncoding.DecodeString(reqId)
	if err != nil || len(decodedBytes) < 4 {
		logs.Error("Invalid reqid", reqId, err)
		return
	}

	newBytes := decodedBytes[4:]
	newBytesLen := len(newBytes)
	newStr := ""
	for i := newBytesLen - 1; i >= 0; i-- {
		newStr += fmt.Sprintf("%02X", newBytes[i])
	}
	unixNano, err := strconv.ParseInt(newStr, 16, 64)
	if err != nil {
		logs.Error("Invalid reqid", reqId, err)
		return
	}
	dstDate := time.Unix(0, unixNano)
	fmt.Println(fmt.Sprintf("%04d-%02d-%02d/%02d-%02d", dstDate.Year(), dstDate.Month(), dstDate.Day(),
		dstDate.Hour(), dstDate.Minute()))
}

// 生成随机的字符串
func CreateRandString(num int) (rcode string) {
	if num <= 0 || num > len(ALPHA_LIST) {
		rcode = ""
		return
	}

	buffer := make([]byte, num)
	_, err := rand.Read(buffer)
	if err != nil {
		rcode = ""
		return
	}

	for _, b := range buffer {
		index := int(b) / len(ALPHA_LIST)
		rcode += string(ALPHA_LIST[index])
	}

	return
}

func ParseLine(line, sep string) []string {
	if strings.TrimSpace(sep) == "" {
		return strings.Fields(line)
	}
	return strings.Split(line, sep)
}
