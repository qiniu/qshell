# 简介
`awslist` 使用亚马逊的 List Objects V2 接口[文档](https://docs.aws.amazon.com/AmazonS3/latest/API/v2-RESTBucketGET.html)， 获取空间中的文件， 打印到标准输出。

该命令的数据格式为：
```
<文件名>\t<文件大小>\t<Etag>\t<最后修改时间>    //注： \t 为 Tab 键
```

当程序列举的过程中遇到错误，比如网络断开等， 会把当前的 ContinuationToken 打印到标准错误输出上, 可以使用 shell 重定向把标准输出到一个文件， 这样可以方便地找到 continuationToken 继续列举。

# 格式
```
qshell awslist [-p <Prefix>][-n <maxKeys>][-m <ContinuationToken>] -S <AwsSecretKey> -A <AwsID> <AwsBucket> <AwsRegion>
```

# 帮助文档
可以在命令行输入如下命令获取帮助文档：
```
$ qshell awslist -h
```

# 参数
- AwsBucket: 亚马逊存储空间名称
- AwsRegion: 亚马逊存储空间所在的地区

# 选项
- -A/--aws-id：亚马逊账户的Access Key ID 。【必选】
- -S/--aws-secret-key：亚马逊账户的Secret Key 。【必选】
- -p/--prefix：亚马逊存储空间要抓取资源的前缀。【可选】
- -n/-max-keys：亚马逊接口每次返回的数据条目数量。【可选】
- -m/--continuation-token：亚马逊接口数据每次会返回的token, 用于下次列举。【可选】

# 列举
使用场景：
列举亚马逊存储空间中所有的文件

假如要迁移的亚马逊账户的 Access Key ID, SecretKey 为：
- AWS_ACCESS_KEY_ID = "12345"
- AWS_SECRET_KEY = "6789"

亚马逊存储空间名为：
AWS_BUCKET = "aws-bucket"

亚马逊空间所在地区为：
AWS_REGION = "us-west-2"

可以使用如下命令进行列举：
```
$ qshell awslist -A 12345 -S 6789 aws-bucket us-west-2
```
