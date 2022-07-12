# 简介
`awsfetch` 迁移亚马逊存储空间的数据到七牛存储空间。 该命令需要用到亚马逊账户的 AccessKeyID 和 SecretKey, 创建访问密钥可以参考：[创建密钥](https://docs.aws.amazon.com/zh_cn/general/latest/gr/managing-aws-access-keys.html)。

该命令使用了七牛的 fetch 接口进行抓取， 需要可以直接访问的网络资源链接， 因此需要亚马逊存储开启公共可访问, 公共可访问开启参考: [文档](https://aws.amazon.com/cn/premiumsupport/knowledge-center/read-access-objects-s3-bucket/)。

因为该命令使用了七牛 fetch 接口，对于较大的资源(大于 100M), 有抓取超时的可能性。

该命令首先使用亚马逊的 List Objects V2 接口[文档](https://docs.aws.amazon.com/AmazonS3/latest/API/v2-RESTBucketGET.html) 获取空间中的文件， 然后七牛的 fetch 接口[文档](https://developer.qiniu.com/kodo/api/1263/fetch)进行抓取。

# 格式
```
qshell awsfetch [-p <Prefix>][-n <maxKeys>][-m <ContinuationToken>][-c <threadCount>][-u <QiniuUpHost>] -S <AwsSecretKey> -A <AwsID> [-s <SuccessList>][-e <FailureList>] <AwsBucket> <AwsRegion> <QiniuBucket>
```

# 帮助文档
可以在命令行输入如下命令获取帮助文档：
```
$ qshell awsfetch -h
```

# 参数
- AwsBucket: 亚马逊存储空间名称。【必选】
- AwsRegion: 亚马逊存储空间所在的地区。【必选】
- QiniuBucket: 七牛存储空间名称。【必选】

# 选项
- -A/--aws-id：亚马逊账户的 Access Key ID 。【必选】
- -S/--aws-secret-key：亚马逊账户的 Secret Key 。【必选】
- -p/--prefix：亚马逊存储空间要抓取资源的前缀。 【可选】
- -n/--max-keys：亚马逊接口每次返回的数据条目数量。 【可选】
- -m/--continuation-token：亚马逊接口数据每次会返回的 token, 用于下次列举。 【可选】
- -c/--thead-count：抓取的线程数, 默认为 20。 【可选】
- -u/--up-host：抓取的资源上传到七牛存储时的上传 HOST 。 【可选】
- -s/--success-list：文件抓取成功后，将文件信息导入到此文件中，每行一个文件。不配置导出。 【可选】
- -e/--failure-list：文件抓取失败后，将文件信息导入到此文件中，每行一个文件。不配置导出。 【可选】

# 亚马逊存储数据迁移到七牛存储
使用场景：
迁移亚马逊存储空间到七牛存储空间。

假如要迁移的亚马逊账户的 Access Key ID, SecretKey 为：
- AWS_ACCESS_KEY_ID = "12345"
- AWS_SECRET_KEY = "6789"

亚马逊存储空间名为：
AWS_BUCKET = "aws-bucket"

亚马逊存储空间所在地区为：
AWS_REGION = "us-west-2"

七牛存储空间名为：
QINIU_BUCKET = "qiniu-bucket"

导出失败的文件列表到 "failure.txt"

可以使用如下命令进行迁移：
```
$ qshell awsfetch -A 12345 -S 6789 -e failure.txt aws-bucket us-west-2 qiniu-bucket 
```
