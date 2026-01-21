package operations

import (
	"bufio"
	"context"
	"encoding/json"
	"math/rand"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/qiniu/go-sdk/v7/storagev2/apis"
	"github.com/qiniu/go-sdk/v7/storagev2/apis/create_share"
	"github.com/qiniu/go-sdk/v7/storagev2/apis/verify_share"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
)

type CreateShareInfo struct {
	apis.CreateShareRequest

	OutputPath string
}

func (info *CreateShareInfo) Check() *data.CodeError {
	if len(info.Bucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}
	if len(info.Prefix) == 0 {
		return alert.CannotEmptyError("Prefix", "")
	}
	return nil
}

const randomExtractCodeChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func CreateShare(cfg *iqshell.Config, info CreateShareInfo) {
	if err := createShare(cfg, &info); err != nil {
		data.SetCmdStatusError()
		log.ErrorF("Create share Failed, [%s:%s], Error: %v", info.Bucket, info.Prefix, err)
	}
}

func createShare(cfg *iqshell.Config, info *CreateShareInfo) error {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: info,
	}); !shouldContinue {
		return nil
	}

	storagev2Client, codeErr := bucket.GetStorageV2()
	if codeErr != nil {
		return codeErr
	}

	if info.DurationSeconds == 0 {
		info.DurationSeconds = 15 * 60
	}
	if info.ExtractCode == "" {
		bytes := []byte(randomExtractCodeChars)
		extractCodeBytes := make([]byte, 0, 6)
		for i := 0; i < 6; i++ {
			n := rand.Intn(len(bytes))
			extractCodeBytes = append(extractCodeBytes, bytes[n])
		}
		info.ExtractCode = string(extractCodeBytes)
	}

	response, err := storagev2Client.CreateShare(context.Background(), &info.CreateShareRequest, nil)
	if err != nil {
		return err
	}

	body, err := newCreateShareReponseBody(cfg, info, response)
	if err != nil {
		return err
	}

	if info.OutputPath != "" {
		createdFile, err := os.Create(info.OutputPath)
		if err != nil {
			return err
		}
		defer createdFile.Close()
		if err := json.NewEncoder(createdFile).Encode(&body); err != nil {
			return err
		}
	} else {
		log.AlertF("Link:\n%s", body.Link)
		log.AlertF("Extract Code:\n%s", body.ExtractCode)
		log.AlertF("Expire:\n%s", body.WillExpireAt.Local().Format(time.DateTime+" -0700"))
	}
	return nil
}

type createShareReponseBody struct {
	Link         string    `json:"link"`
	ExtractCode  string    `json:"extract_code"`
	WillExpireAt time.Time `json:"will_expire_at"`
}

func newCreateShareReponseBody(cfg *iqshell.Config, info *CreateShareInfo, response *create_share.Response) (*createShareReponseBody, error) {
	expires, err := time.Parse(time.RFC3339, response.Expires)
	if err != nil {
		return nil, err
	}
	urlString := cfg.CmdCfg.GetPortalHost()
	if !strings.Contains(urlString, "://") {
		if cfg.CmdCfg.IsUseHttps() {
			urlString = "https://" + urlString
		} else {
			urlString = "http://" + urlString
		}
	}
	urlString += "/kodo-shares/verify"
	parsedUrl, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}
	query := parsedUrl.Query()
	query.Set("id", response.Id)
	query.Set("token", response.Token)
	parsedUrl.RawQuery = query.Encode()
	return &createShareReponseBody{Link: parsedUrl.String(), ExtractCode: info.ExtractCode, WillExpireAt: expires}, nil
}

type ListShareInfo struct {
	LinkURL     string
	ExtractCode string
	Prefix      string
	Limit       int64
	Marker      string
}

func (info *ListShareInfo) Check() *data.CodeError {
	if len(info.LinkURL) == 0 {
		return alert.CannotEmptyError("Link", "")
	}
	if len(info.ExtractCode) == 0 {
		return alert.CannotEmptyError("ExtractCode", "")
	}
	if info.Limit < 0 {
		return alert.Error("Limit should not be negative", "")
	}
	return nil
}

func ListShare(cfg *iqshell.Config, info ListShareInfo) {
	if err := listShare(cfg, &info); err != nil {
		data.SetCmdStatusError()
		log.ErrorF("List share Failed, [%s], Error: %v", info.LinkURL, err)
	}
}

func listShare(cfg *iqshell.Config, info *ListShareInfo) error {
	if !strings.HasPrefix(info.LinkURL, "http://") && !strings.HasPrefix(info.LinkURL, "https://") {
		linkUrl, extractCode, err := readLinkURLFromPath(info.LinkURL)
		if err != nil {
			return err
		}
		info.LinkURL = linkUrl
		if info.ExtractCode == "" {
			info.ExtractCode = extractCode
		}
	}
	if info.ExtractCode == "" {
		info.ExtractCode = promptExtractCode()
	}

	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: info,
	}); !shouldContinue {
		return nil
	}

	storagev2Client, codeErr := bucket.GetStorageV2()
	if codeErr != nil {
		return codeErr
	}

	parsedUrl, err := url.Parse(info.LinkURL)
	if err != nil {
		return err
	}

	response, err := storagev2Client.VerifyShare(context.Background(), &verify_share.Request{
		ShareId:     parsedUrl.Query().Get("id"),
		Token:       parsedUrl.Query().Get("token"),
		ExtractCode: info.ExtractCode,
	}, nil)
	if err != nil {
		return err
	}

	s3Svc, err := getS3Service(response, cfg)
	if err != nil {
		return err
	}
	prefix := info.Prefix
	if prefix == "" {
		prefix = response.Prefix
	}
	input := s3.ListObjectsV2Input{
		Bucket: aws.String(response.BucketId),
	}
	if prefix != "" {
		input.Prefix = aws.String(prefix)
	}
	if info.Marker != "" {
		input.ContinuationToken = aws.String(info.Marker)
	}

	var (
		stats      = listedStats{prefix: info.Prefix}
		restListed = info.Limit
		listOutput *s3.ListObjectsV2Output
	)
	for restListed > 0 || info.Limit == 0 {
		if restListed > 0 && restListed < 1000 {
			input.MaxKeys = aws.Int64(restListed)
		}
		listOutput, err = s3Svc.ListObjectsV2(&input)
		if err != nil {
			return err
		}
		for _, s3Object := range listOutput.Contents {
			if restListed > 0 {
				restListed -= 1
			}
			if strings.HasSuffix(aws.StringValue(s3Object.Key), "/") && aws.Int64Value(s3Object.Size) == 0 {
				printListedS3Directory(s3Object)
				stats.directoryNumbers += 1
			} else {
				printListedS3Object(s3Object)
				stats.objectNumbers += 1
				stats.totalSize += aws.Int64Value(s3Object.Size)
			}
		}
		if listOutput.NextContinuationToken == nil {
			break
		}
		input.ContinuationToken = listOutput.NextContinuationToken
	}
	if listOutput.NextContinuationToken != nil {
		log.AlertF("Marker: %s", aws.StringValue(listOutput.NextContinuationToken))
	}
	printListedStats(&stats)
	return nil
}

func printListedS3Object(object *s3.Object) {
	log.AlertF("%s\t%d\t%s\t%s", aws.StringValue(object.Key), aws.Int64Value(object.Size), aws.StringValue(object.StorageClass), aws.TimeValue(object.LastModified))
}

func printListedS3Directory(object *s3.Object) {
	log.AlertF("%s", aws.StringValue(object.Key))
}

type listedStats struct {
	prefix           string
	totalSize        int64
	directoryNumbers int64
	objectNumbers    int64
}

func printListedStats(info *listedStats) {
	if info.prefix == "" {
		log.AlertF("Total size: %s", utils.FormatFileSize(info.totalSize))
	} else {
		log.AlertF("Total size of prefix [%s]: %s", info.prefix, utils.FormatFileSize(info.totalSize))
	}
	log.AlertF("Folder number: %d", info.directoryNumbers)
	log.AlertF("File number: %d", info.objectNumbers)
}

type CopyShareInfo struct {
	FromPath    string
	ToPath      string
	LinkURL     string
	ExtractCode string
	Recursive   bool
}

func (info *CopyShareInfo) Check() *data.CodeError {
	if len(info.ToPath) == 0 {
		return alert.CannotEmptyError("ToURL", "")
	}
	if len(info.LinkURL) == 0 {
		return alert.CannotEmptyError("Link", "")
	}
	if len(info.ExtractCode) == 0 {
		return alert.CannotEmptyError("ExtractCode", "")
	}
	return nil
}

func CopyShare(cfg *iqshell.Config, info CopyShareInfo) {
	if err := copyShare(cfg, &info); err != nil {
		data.SetCmdStatusError()
		log.ErrorF("Copy share Failed, [%s], Error: %v", info.LinkURL, err)
	}
}

func copyShare(cfg *iqshell.Config, info *CopyShareInfo) error {
	if !strings.HasPrefix(info.LinkURL, "http://") && !strings.HasPrefix(info.LinkURL, "https://") {
		linkUrl, extractCode, err := readLinkURLFromPath(info.LinkURL)
		if err != nil {
			return err
		}
		info.LinkURL = linkUrl
		if info.ExtractCode == "" {
			info.ExtractCode = extractCode
		}
	}
	if info.ExtractCode == "" {
		info.ExtractCode = promptExtractCode()
	}

	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: info,
	}); !shouldContinue {
		return nil
	}

	storagev2Client, codeErr := bucket.GetStorageV2()
	if codeErr != nil {
		return codeErr
	}

	parsedLinkUrl, err := url.Parse(info.LinkURL)
	if err != nil {
		return err
	}

	response, err := storagev2Client.VerifyShare(context.Background(), &verify_share.Request{
		ShareId:     parsedLinkUrl.Query().Get("id"),
		Token:       parsedLinkUrl.Query().Get("token"),
		ExtractCode: info.ExtractCode,
	}, nil)
	if err != nil {
		return err
	}

	s3Svc, err := getS3Service(response, cfg)
	if err != nil {
		return err
	}

	s3session, err := session.NewSession(&s3Svc.Config)
	if err != nil {
		return err
	}
	s3Downloader := s3manager.NewDownloader(s3session)

	fromPrefix := info.FromPath
	if fromPrefix == "" {
		fromPrefix = response.Prefix
	}

	toPath := strings.TrimPrefix(info.ToPath, "file://")
	if strings.Contains(toPath, "://") {
		return err
	}

	if info.Recursive && strings.HasSuffix(fromPrefix, "/") {
		fromParentDictionaryPrefix := fromPrefix[:strings.LastIndex(strings.TrimSuffix(fromPrefix, "/"), "/")+1]
		input := s3.ListObjectsV2Input{
			Bucket: aws.String(response.BucketId),
		}
		if fromPrefix != "" {
			input.Prefix = aws.String(fromPrefix)
		}
		for {
			listOutput, err := s3Svc.ListObjectsV2(&input)
			if err != nil {
				return err
			}
			for _, s3Object := range listOutput.Contents {
				relativePath := strings.TrimPrefix(aws.StringValue(s3Object.Key), fromParentDictionaryPrefix)
				if filepath.Separator != '/' {
					relativePath = strings.Replace(relativePath, "/", string(filepath.Separator), -1)
				}
				downloadPath := toPath
				if relativePath != "" {
					downloadPath = filepath.Join(toPath, relativePath)
				}
				if strings.HasSuffix(aws.StringValue(s3Object.Key), "/") && aws.Int64Value(s3Object.Size) == 0 {
					err = os.MkdirAll(downloadPath, 0o700)
				} else {
					if err = os.MkdirAll(filepath.Dir(downloadPath), 0o700); err != nil {
						return err
					}
					err = s3DownloadObjectToPath(s3Downloader, response.BucketId, aws.StringValue(s3Object.Key), downloadPath)
				}
				if err != nil {
					return err
				}
			}
			if listOutput.NextContinuationToken == nil {
				break
			}
			input.ContinuationToken = listOutput.NextContinuationToken
		}
	} else {
		if err = s3StatObject(s3Svc, response.BucketId, fromPrefix); err != nil {
			return err
		}
		downloadPath := toPath
		onlyMkDir := strings.HasSuffix(fromPrefix, "/")
		fromPrefix = strings.TrimSuffix(fromPrefix, "/")
		downloadPath = filepath.Join(downloadPath, fromPrefix[(strings.LastIndex(fromPrefix, "/")+1):])
		if onlyMkDir {
			if err = os.MkdirAll(downloadPath, 0o700); err != nil {
				return err
			}
		} else {
			if err = os.MkdirAll(filepath.Dir(downloadPath), 0o700); err != nil {
				return err
			}
			if err = s3DownloadObjectToPath(s3Downloader, response.BucketId, fromPrefix, downloadPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func readLinkURLFromPath(path string) (string, string, error) {
	path = strings.TrimPrefix(path, "file://")
	b, err := os.ReadFile(path)
	if err != nil {
		return "", "", err
	}
	var body createShareReponseBody
	if err = json.Unmarshal(b, &body); err != nil {
		return "", "", err
	}
	return body.Link, body.ExtractCode, nil
}

func promptExtractCode() string {
	log.AlertF("Input Extract Code:")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func s3StatObject(s3Service *s3.S3, fromBucketId, key string) error {
	_, err := s3Service.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(fromBucketId),
		Key:    aws.String(key),
	})
	return err
}

func s3DownloadObjectToPath(downloader *s3manager.Downloader, fromBucketId, key, downloadPath string) error {
	file, err := os.OpenFile(downloadPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(fromBucketId),
		Key:    aws.String(key),
	})
	return err
}

func getAwsConfig(response *verify_share.Response, cfg *iqshell.Config) *aws.Config {
	config := aws.NewConfig()
	if !cfg.CmdCfg.IsUseHttps() {
		config.WithDisableSSL(true)
	}
	config.WithEndpoint(response.Endpoint)
	config.WithRegion(response.Region)
	config.WithCredentials(credentials.NewStaticCredentials(response.FederatedAk, response.FederatedSk, response.SessionToken))
	if cfg.DebugEnable {
		config.WithLogLevel(aws.LogDebug)
	} else if cfg.DDebugEnable {
		config.WithLogLevel(aws.LogDebugWithHTTPBody)
	}
	return config
}

func getS3Service(response *verify_share.Response, cfg *iqshell.Config) (*s3.S3, error) {
	s3session, err := session.NewSession(getAwsConfig(response, cfg))
	if err != nil {
		return nil, err
	}
	s3service := s3.New(s3session)
	return s3service, nil
}
