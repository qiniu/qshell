package aws

type FetchInfo struct {
	AwsBucketInfo ListBucketInfo
}

//func Fetch(info FetchInfo) {
//	if info.AwsBucketInfo.AccessKey == "" || info.AwsBucketInfo.SecretKey == "" {
//		log.Error(alert.CannotEmpty("AWS ID and SecretKey", ""))
//		os.Exit(data.STATUS_ERROR)
//	}
//
//	if info.AwsBucketInfo.MaxKeys <= 0 || info.AwsBucketInfo.MaxKeys > 1000 {
//		log.Warning("max key:%d out of range {}, change to 1000", info.AwsBucketInfo.MaxKeys)
//		info.AwsBucketInfo.MaxKeys = 1000
//	}
//
//	// check AWS region
//	if info.AwsBucketInfo.Region == "" {
//		log.Error(alert.CannotEmpty("AWS region", ""))
//		os.Exit(data.STATUS_ERROR)
//	}
//
//	if o.threadCount <= 0 || o.threadCount >= 1000 {
//		o.threadCount = 20
//	}
//
//	o.initBucketManager()
//	o.initFileExporter()
//	o.initUpHost(positionalArgs[2])
//
//	// check AWS region
//	if positionalArgs[1] == "" {
//		fmt.Fprintf(os.Stderr, "AWS region cannot be empty\n")
//		os.Exit(1)
//	}
//
//	awsBucket, region, qiniuBucket := positionalArgs[0], positionalArgs[1], positionalArgs[2]
//
//	// AWS related code
//	s3session := session.New()
//	s3session.Config.WithRegion(region)
//	s3session.Config.WithCredentials(credentials.NewStaticCredentials(o.id, o.secretKey, ""))
//
//	svc := s3.New(s3session)
//	input := &s3.ListObjectsV2Input{
//		Bucket:    aws.String(awsBucket),
//		Prefix:    aws.String(o.prefix),
//		Delimiter: aws.String(o.delimiter),
//		MaxKeys:   aws.Int64(o.maxKeys),
//	}
//	if o.ctoken != "" {
//		input.ContinuationToken = aws.String(o.ctoken)
//	}
//	itemc := make(chan *data.FetchItem)
//	donec := make(chan struct{})
//
//	go fetchChannel(itemc, donec, &o.fetchConfig)
//
//	for {
//		result, err := svc.ListObjectsV2(input)
//
//		if err != nil {
//			if aerr, ok := err.(awserr.Error); ok {
//				switch aerr.Code() {
//				case s3.ErrCodeNoSuchBucket:
//					fmt.Fprintln(os.Stderr, s3.ErrCodeNoSuchBucket, aerr.Error())
//				default:
//					fmt.Fprintln(os.Stderr, aerr.Error())
//				}
//			} else {
//				// Print the error, cast err to awserr.Error to get the Code and
//				// Message from an error.
//				fmt.Fprintln(os.Stderr, err.Error())
//			}
//			close(itemc)
//			fmt.Fprintf(os.Stderr, "ContinuationToken: %v\n", input.ContinuationToken)
//			os.Exit(1)
//		}
//		for _, obj := range result.Contents {
//			if strings.HasSuffix(*obj.Key, "/") && *obj.Size == 0 { // 跳过目录
//				continue
//			}
//
//			item := &data.FetchItem{
//				Bucket:    qiniuBucket,
//				Key:       *obj.Key,
//				RemoteUrl: awsUrl(awsBucket, region, *obj.Key),
//			}
//			itemc <- item
//		}
//
//		if *result.IsTruncated {
//			input.ContinuationToken = result.NextContinuationToken
//		} else {
//			break
//		}
//	}
//	close(itemc)
//
//	<-donec
//}
//
//func awsUrl(awsBucket, region, key string) string {
//	return fmt.Sprintf("https://%s.s3-%s.amazonaws.com/%s", awsBucket, region, key)
//}
