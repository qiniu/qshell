package cmd

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/config"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/cobra"
)

type awsfetchOptions struct {
	fetchConfig
	awslistOptions
}

type awslistOptions struct {
	// aws continuation token
	ctoken string

	delimiter string

	maxKeys int64

	prefix string

	// aws id and secretKey
	id        string
	secretKey string
}

func (lo *awslistOptions) Run(cmd *cobra.Command, positionalArgs []string) {

	lo.checkOptions()

	awsBucket := positionalArgs[0]
	region := positionalArgs[1]
	// check AWS region
	if region == "" {
		fmt.Fprintf(os.Stderr, "AWS region cannot be empty\n")
		os.Exit(1)
	}

	// AWS related code
	s3session := session.New()
	s3session.Config.WithRegion(region)
	s3session.Config.WithCredentials(credentials.NewStaticCredentials(lo.id, lo.secretKey, ""))

	svc := s3.New(s3session)
	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(awsBucket),
		Prefix:  aws.String(lo.prefix),
		MaxKeys: aws.Int64(lo.maxKeys),
	}
	if lo.ctoken != "" {
		input.ContinuationToken = aws.String(lo.ctoken)
	}

	for {
		result, err := svc.ListObjectsV2(input)

		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case s3.ErrCodeNoSuchBucket:
					fmt.Fprintln(os.Stderr, s3.ErrCodeNoSuchBucket, aerr.Error())
				default:
					fmt.Fprintln(os.Stderr, aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Fprintln(os.Stderr, err.Error())
			}
			fmt.Fprintf(os.Stderr, "ContinuationToken: %v\n", input.ContinuationToken)
			os.Exit(1)
		}
		for _, obj := range result.Contents {
			if strings.HasSuffix(*obj.Key, "/") && *obj.Size == 0 { // 跳过目录
				continue
			}
			fmt.Printf("%s\t%d\t%s\t%s\n", *obj.Key, *obj.Size, *obj.ETag, *obj.LastModified)
		}

		if *result.IsTruncated {
			input.ContinuationToken = result.NextContinuationToken
		} else {
			break
		}
	}

}

func (lo *awslistOptions) checkOptions() {
	if lo.id == "" || lo.secretKey == "" {
		fmt.Fprintf(os.Stderr, "AWS ID and SecretKey cannot be empty\n")
		os.Exit(1)
	}
	if lo.maxKeys <= 0 || lo.maxKeys > 1000 {
		lo.maxKeys = 1000
	}
}

func awsUrl(awsBucket, region, key string) string {
	return fmt.Sprintf("https://%s.s3-%s.amazonaws.com/%s", awsBucket, region, key)
}

func (o *awsfetchOptions) Run(cmd *cobra.Command, positionalArgs []string) {
	o.checkOptions()

	if o.threadCount <= 0 || o.threadCount >= 1000 {
		o.threadCount = 20
	}

	o.initBucketManager()
	o.initFileExporter()
	o.initUpHost(positionalArgs[2])

	// check AWS region
	if positionalArgs[1] == "" {
		fmt.Fprintf(os.Stderr, "AWS region cannot be empty\n")
		os.Exit(1)
	}

	awsBucket, region, qiniuBucket := positionalArgs[0], positionalArgs[1], positionalArgs[2]

	// AWS related code
	s3session := session.New()
	s3session.Config.WithRegion(region)
	s3session.Config.WithCredentials(credentials.NewStaticCredentials(o.id, o.secretKey, ""))

	svc := s3.New(s3session)
	input := &s3.ListObjectsV2Input{
		Bucket:    aws.String(awsBucket),
		Prefix:    aws.String(o.prefix),
		Delimiter: aws.String(o.delimiter),
		MaxKeys:   aws.Int64(o.maxKeys),
	}
	if o.ctoken != "" {
		input.ContinuationToken = aws.String(o.ctoken)
	}
	itemc := make(chan *config.FetchItem)
	donec := make(chan struct{})

	go fetchChannel(itemc, donec, &o.fetchConfig)

	for {
		result, err := svc.ListObjectsV2(input)

		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case s3.ErrCodeNoSuchBucket:
					fmt.Fprintln(os.Stderr, s3.ErrCodeNoSuchBucket, aerr.Error())
				default:
					fmt.Fprintln(os.Stderr, aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Fprintln(os.Stderr, err.Error())
			}
			close(itemc)
			fmt.Fprintf(os.Stderr, "ContinuationToken: %v\n", input.ContinuationToken)
			os.Exit(1)
		}
		for _, obj := range result.Contents {
			if strings.HasSuffix(*obj.Key, "/") && *obj.Size == 0 { // 跳过目录
				continue
			}

			item := &config.FetchItem{
				Bucket:    qiniuBucket,
				Key:       *obj.Key,
				RemoteUrl: awsUrl(awsBucket, region, *obj.Key),
			}
			itemc <- item
		}

		if *result.IsTruncated {
			input.ContinuationToken = result.NextContinuationToken
		} else {
			break
		}
	}
	close(itemc)

	<-donec
}

// NewCmdAwsFetch 返回一个cobra.Command指针
func NewCmdAwsFetch() *cobra.Command {
	options := awsfetchOptions{}

	awsFetch := &cobra.Command{
		Use:   "awsfetch [-p <Prefix>] [-n <maxKeys>] [-m <ContinuationToken>] [-c <threadCount>][-u <Qiniu UpHost>] -S <AwsSecretKey> -A <AwsID> <awsBucket> <awsRegion> <qiniuBucket>",
		Short: "Copy data from AWS bucket to qiniu bucket",
		Args:  cobra.ExactArgs(3),
		Run:   options.Run,
	}

	awsFetch.Flags().StringVarP(&options.prefix, "prefix", "p", "", "list AWS bucket with this prefix if set")
	awsFetch.Flags().Int64VarP(&options.maxKeys, "max-keys", "n", 1000, "list AWS bucket with numbers of keys returned each time limited by this number if set")
	awsFetch.Flags().StringVarP(&options.ctoken, "continuation-token", "m", "", "AWS list continuation token")
	awsFetch.Flags().IntVarP(&options.threadCount, "thead-count", "c", 20, "maximum of fetch thread")
	awsFetch.Flags().StringVarP(&options.upHost, "up-host", "u", "", "Qiniu fetch up host")
	awsFetch.Flags().StringVarP(&options.secretKey, "aws-secret-key", "S", "", "AWS secret key")
	awsFetch.Flags().StringVarP(&options.id, "aws-id", "A", "", "AWS ID")
	awsFetch.Flags().StringVarP(&options.successFname, "success-list", "s", "", "success fetch key list")
	awsFetch.Flags().StringVarP(&options.failureFname, "failure-list", "e", "", "error fetch key list")

	return awsFetch
}

// NewCmdAwsList 返回一个cobra.Command指针
// 该命令列举亚马逊存储空间中的文件, 会忽略目录
func NewCmdAwsList() *cobra.Command {
	options := awslistOptions{}

	awsList := &cobra.Command{
		Use:   "awslist [-p <Prefix>] [-n <maxKeys>] [-m <ContinuationToken>] -S <AwsSecretKey> -A <AwsID> <awsBucket> <awsRegion>",
		Short: "List Objects in AWS bucket",
		Args:  cobra.ExactArgs(2),
		Run:   options.Run,
	}

	awsList.Flags().StringVarP(&options.prefix, "prefix", "p", "", "list AWS bucket with this prefix if set")
	awsList.Flags().Int64VarP(&options.maxKeys, "max-keys", "n", 1000, "list AWS bucket with numbers of keys returned each time limited by this number if set")
	awsList.Flags().StringVarP(&options.ctoken, "continuation-token", "m", "", "AWS list continuation token")
	awsList.Flags().StringVarP(&options.secretKey, "aws-secret-key", "S", "", "AWS secret key")
	awsList.Flags().StringVarP(&options.id, "aws-id", "A", "", "AWS ID")

	return awsList
}

func init() {
	RootCmd.AddCommand(NewCmdAwsFetch())
	RootCmd.AddCommand(NewCmdAwsList())
}
