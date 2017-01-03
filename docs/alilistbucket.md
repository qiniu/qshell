# 简介

`alilistbucket` 用来获取阿里云OSS空间中的文件列表，然后可以使用工具处理文件列表，构建出所有资源的外链，通过 [qfetch](https://github.com/qiniu/qfetch) 迁移到七牛。

# 格式

```
qshell alilistbucket <DataCenter> <Bucket> <AccessKeyId> <AccessKeySecret> [Prefix] <ListBucketResultFile>
```

# 参数

|参数名|描述|可选参数|
|---------|------------|------|
|DataCenter|阿里云OSS空间所在的数据中心|N|
|Bucket|阿里云OSS空间名称，可以为公开空间或私有空间|N|
|AccessKeyId|阿里云账号对应的AccessKeyId [获取](https://ak-console.aliyun.com/#/accesskey)|N|
|AccessKeySecret|阿里云账号对应的AccessKeySecret [获取](https://ak-console.aliyun.com/#/accesskey)|N|
|Prefix|阿里云OSS空间中文件的前缀|Y|
|ListBucketResultFile|文件列表保存的文件名称，可以为绝对路径或者相对路径|N|

阿里OSS公网数据中心

|地点|域名|
|-------|-------|
|杭州|oss-cn-hangzhou.aliyuncs.com|
|青岛|oss-cn-qingdao.aliyuncs.com|
|香港|oss-cn-hongkong.aliyuncs.com|
|北京|oss-cn-beijing.aliyuncs.com|
|深圳|oss-cn-shenzhen.aliyuncs.com|

阿里OSS内网数据中心

|地点|域名|
|-------|-------|
|杭州|oss-cn-hangzhou-internal.aliyuncs.com|
|青岛|oss-cn-qingdao-internal.aliyuncs.com|
|香港|oss-cn-hongkong-internal.aliyuncs.com|
|北京|oss-cn-beijing-internal.aliyuncs.com|
|深圳|oss-cn-shenzhen-internal.aliyuncs.com|

# 示例

1.获取阿里云OSS空间`qdisk-hz`里面的所有文件列表：

```
qshell alilistbucket oss-cn-hangzhou.aliyuncs.com qdisk-hz poeDElTwLc2w0iFJ pPlaT3umFa1lcXTwp7N5nVQt9av1yg qdisk-hz.list.txt
```

2.获取阿里云OSS空间`qdisk-hz`里面的以`2015/01/18`为前缀的文件列表：

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