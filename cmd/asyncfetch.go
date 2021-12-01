package cmd

import (
	"fmt"
	storage2 "github.com/qiniu/qshell/v2/iqshell/storage"
	"github.com/qiniu/qshell/v2/iqshell/utils"
	"os"
	"strconv"
	"time"

	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/spf13/cobra"
)

// NewCmdAsyncFetch 返回一个cobra.Command指针
// 该命令使用七牛异步抓取的接口
func NewCmdAsyncFetch() *cobra.Command {
	options := asyncFetchOptions{}

	asyncFetch := &cobra.Command{
		Use:   "abfetch <Bucket> [-i <urlList>]",
		Short: "Async Batch fetch network resources to qiniu Bucket",
		Args:  cobra.ExactArgs(1),
		Run:   options.Run,
	}

	asyncFetch.Flags().StringVarP(&options.host, "host", "t", "", "download HOST header")
	asyncFetch.Flags().StringVarP(&options.callbackUrl, "callback-url", "a", "", "callback url")
	asyncFetch.Flags().StringVarP(&options.callbackBody, "callback-body", "b", "", "callback body")
	asyncFetch.Flags().StringVarP(&options.callbackHost, "callback-host", "T", "", "callback HOST")
	asyncFetch.Flags().IntVarP(&options.fileType, "storage-type", "g", 0, "storage type")
	asyncFetch.Flags().StringVarP(&options.inputFile, "input-file", "i", "", "input file with urls")
	asyncFetch.Flags().IntVarP(&options.threadCount, "thread-count", "c", 20, "thread count")
	asyncFetch.Flags().StringVarP(&options.successFname, "success-list", "s", "", "success fetch list")
	asyncFetch.Flags().StringVarP(&options.failureFname, "failure-list", "e", "", "error fetch list")

	return asyncFetch
}

// NewCmdAsyncCheck 用来查询异步抓取的结果
func NewCmdAsyncCheck() *cobra.Command {

	asyncCheck := &cobra.Command{
		Use:   "acheck <Bucket> <ID>",
		Short: "Check Async fetch status",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, positionalArgs []string) {
			bm := storage2.GetBucketManager()

			ret, err := bm.CheckAsyncFetchStatus(positionalArgs[0], positionalArgs[1])
			if err != nil {
				fmt.Fprintf(os.Stderr, "CheckAsyncFetchStatus: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(ret)
		},
	}
	return asyncCheck
}

type asyncFetchOptions struct {
	// 从指定URL下载时指定的HOST
	host string

	// 设置了该值，抓取的过程使用文件md5值进行校验, 校验失败不存在七牛空间
	md5 string

	// 设置了该值， 抓取的过程中使用etag进行校验，失败不保存在存储空间中
	etag string

	// 抓取成功的回调地址
	callbackUrl string

	callbackBody string

	callbackBodyType string

	// 回调时使用的HOST
	callbackHost string

	// 文件存储类型， 0 标准存储， 1 低频存储
	fileType int

	// 输入访问地址列表
	inputFile string

	fetchConfig
}

type asyncItem struct {
	id       string
	url      string
	key      string
	size     uint64
	bucket   string
	duration int
	waiter   int

	start time.Time
}

func (i *asyncItem) degrade() {
	i.duration = i.duration / 2
	if i.duration <= 0 {
		i.duration = 3
	}
}

func (i *asyncItem) estimatDuration() {
	if i.duration == 0 {
		switch {
		case i.size >= 500*MB:
			i.duration = 40
		case i.size > 200*MB && i.size < 500*MB:
			i.duration = 30
		case i.size > 100*MB && i.size <= 200*MB:
			i.duration = 20
		case i.size <= 10*MB:
			i.duration = 3
		case i.size <= 10*MB:
			i.duration = 6
		case i.size <= 100*MB:
			i.duration = 10
		default:
			i.duration = 3
		}
	}

}

func (i *asyncItem) timeEnough() bool {

	now := time.Now()

	i.estimatDuration()
	if now.Sub(i.start) > time.Duration(i.duration)*time.Second {
		return true
	}
	return false
}

func (ao *asyncFetchOptions) Run(cmd *cobra.Command, positionalArgs []string) {
	bucket := positionalArgs[0]

	var lc chan string
	var err error

	if ao.inputFile != "" {
		lc, err = getLines(ao.inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "get lines from file: %s: %v\n", ao.inputFile, err)
			os.Exit(1)
		}
	} else {
		lc = getLinesFromReader(os.Stdin)
	}
	ao.initFileExporter()
	ao.initBucketManager()

	limitc := make(chan struct{}, ao.threadCount)
	queuec := make(chan *asyncItem, 1000)
	donec := make(chan struct{})

	go func() {

		for item := range queuec {
			counter := 0
			for counter < 3 {
				if item.timeEnough() {

					ret, cErr := ao.bm.CheckAsyncFetchStatus(item.bucket, item.id)
					if cErr != nil {
						fmt.Fprintf(os.Stderr, "CheckAsyncFetchStatus: %v\n", cErr)
					} else if ret.Wait == -1 { // 视频抓取过一次，有可能成功了，有可能失败了
						counter += 1
						_, err := ao.bm.Stat(item.bucket, item.key)
						if err != nil {
							fmt.Fprintf(os.Stderr, "Stat: %s: %v\n", item.key, err)
						} else {
							ao.fileExporter.WriteToSuccessWriter(fmt.Sprintf("%s\t%s\n", item.url, item.key))
							fmt.Printf("fetch %s => %s:%s success\n", item.url, item.bucket, item.key)
							break
						}
					}
					item.degrade()
				}
				time.Sleep(3 * time.Second)
			}
			if counter >= 3 {
				ao.fileExporter.WriteToFailedWriter(fmt.Sprintf("%s\t%d\t%s\n", item.url, item.size, item.key))
				fmt.Fprintf(os.Stderr, "fetch %s => %s:%s failed\n", item.url, item.bucket, item.key)
			}
		}

		donec <- struct{}{}
	}()

	var size uint64
	var pErr error
	for line := range lc {
		limitc <- struct{}{}

		fields := ParseLine(line, "")
		if len(fields) <= 0 {
			continue
		}
		url := fields[0]
		if len(fields) >= 2 {
			size, pErr = strconv.ParseUint(fields[1], 10, 64)
			if pErr != nil {
				ao.fileExporter.WriteToFailedWriter(fmt.Sprintf("%s: %v\n", line, pErr))
				continue
			}
		} else {
			size = 0
		}
		saveKey, pError := utils.KeyFromUrl(url)
		if pError != nil {
			ao.fileExporter.WriteToFailedWriter(fmt.Sprintf("%s: %v\n", line, pError))
			continue
		}
		params := storage.AsyncFetchParam{
			Url:              line,
			Host:             ao.host,
			Bucket:           bucket,
			Key:              saveKey,
			CallbackURL:      ao.callbackUrl,
			CallbackBody:     ao.callbackBody,
			CallbackBodyType: ao.callbackBodyType,
			FileType:         ao.fileType,
		}
		go func(params storage.AsyncFetchParam) {

			ret, aerr := ao.bm.AsyncFetch(params)
			if aerr != nil {
				ao.fileExporter.WriteToFailedWriter(fmt.Sprintf("%s: %v\n", params.Url, aerr))
				<-limitc
				return
			}
			queuec <- &asyncItem{
				id:     ret.Id,
				size:   size,
				waiter: ret.Wait,
				key:    params.Key,
				url:    params.Url,
				bucket: params.Bucket,
				start:  time.Now(),
			}

			<-limitc
		}(params)
	}

	for i := 0; i < ao.threadCount; i++ {
		limitc <- struct{}{}
	}
	close(queuec)

	<-donec
}

func init() {
	RootCmd.AddCommand(NewCmdAsyncFetch())
	RootCmd.AddCommand(NewCmdAsyncCheck())
}
