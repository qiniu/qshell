# 简介
`alilistbucket` 用来列举阿里云 OSS 空间中的文件列表。

# 格式
```
qshell alilistbucket <DataCenter> <Bucket> <AccessKeyId> <AccessKeySecret> [Prefix] <ListBucketResultFile>
```

# 参数
- DataCenter：阿里云 OSS 空间所在的数据中心域名。【必选】
- Bucket：阿里云 OSS 空间名称，可以为公开空间或私有空间。【必选】
- AccessKeyId：阿里云账号对应的 AccessKeyId [获取](https://ak-console.aliyun.com/#/accesskey) 。【必选】
- AccessKeySecret：阿里云账号对应的 AccessKeySecret [获取](https://ak-console.aliyun.com/#/accesskey) 。【必选】
- Prefix：阿里云 OSS 空间中文件的前缀。【可选】
- ListBucketResultFile：文件列表保存的文件名称，可以为绝对路径或者相对路径。【必选】

#### 阿里 OSS 公网数据中心域名
- 杭州：oss-cn-hangzhou.aliyuncs.com
- 青岛：oss-cn-qingdao.aliyuncs.com 
- 香港：oss-cn-hongkong.aliyuncs.com
- 北京：oss-cn-beijing.aliyuncs.com 
- 深圳：oss-cn-shenzhen.aliyuncs.com

#### 阿里 OSS 内网数据中心域名
- 杭州：oss-cn-hangzhou-internal.aliyuncs.com
- 青岛：oss-cn-qingdao-internal.aliyuncs.com 
- 香港：oss-cn-hongkong-internal.aliyuncs.com
- 北京：oss-cn-beijing-internal.aliyuncs.com 
- 深圳：oss-cn-shenzhen-internal.aliyuncs.com

# 示例
1.获取阿里云 OSS 空间 `qdisk-hz` 里面的所有文件列表：
```
qshell alilistbucket oss-cn-hangzhou.aliyuncs.com qdisk-hz poeDElTwLc2w0iFJ pPlaT3umFa1lcXTwp7N5nVQt9av1yg qdisk-hz.list.txt
```

2.获取阿里云 OSS 空间 `qdisk-hz` 里面的以 `2015/01/18` 为前缀的文件列表：
```
qshell alilistbucket oss-cn-hangzhou.aliyuncs.com qdisk-hz poeDElTwLc2w0iFJ pPlaT3umFa1lcXTwp7N5nVQt9av1yg "2015/01/18" qdisk-hz.prefix.list.txt
```

获取的文件内容组织方式为：
```
Key\tSize\tPutTime
```

比如：
```
bucket_domain.png	287727	14215611840000000
```