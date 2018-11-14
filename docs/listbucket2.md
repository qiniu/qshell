# 简介

`listbucket2`用来获取七牛空间里面的文件列表，可以指定文件前缀获取指定的文件列表，如果不指定，则获取所有文件的列表。

获取的文件列表组织格式如下（每个字段用Tab分隔）：

```
Key\tSize\tHash\tPutTime\tMimeType\tFileType\tEndUser
```


参考文档：[资源列举 (list)](http://developer.qiniu.com/code/v6/api/kodo-api/rs/list.html)


# 格式

```
qshell listbucket2 [--prefix <Prefix> | --suffixes <suffixes1,suffixes2>] [--start <StartDate>] [--max-retry <RetryCount>][--end <EndDate>] <Bucket> [-o <ListBucketResultFile>]
```

选项max-retry， 默认列举出错的重试次数为20次，如果希望可以列举完文件列表，不限错误次数，可以设置max-retry为负数。

（1）获取空间中所有的文件列表，这种情况下，可以直接指定 `Bucket` 参数和结果保存文件参数 `ListBucketResultFile` 即可。

```
qshell listbucket2 <Bucket> -o <ListBucketResultFile>
```

 (2) 获取空间所有文件，输出到屏幕上(标准输出)

 ```
 qshell listbucket2 <Bucket> 
 ```

（3）获取空间中指定前缀的文件列表

```
qshell listbucket2 [--prefix <Prefix>] <Bucket> -o <ListBucketResultFile>
```

 (4) 获取空间中指定前缀的文件列表， 输出到屏幕上
 
 ```
 qshell listbucket2 [--prefix <Prefix>] <Bucket>
 ```
 
 (5) 获取2018-10-30到2018-11-02上传的文件
 ```
 qshell listbucket2 --start 2018-10-30 --end 2018-11-02 <Bucket>
 ```
 
 (6) 获取后缀为mp4, html的文件
 
 ```
 qshell listbucket2 --suffixes mp4,html <Bucket>
 ```


# 鉴权

需要在使用了`account`设置了`AccessKey`, `SecretKey` 和  `Name` 的情况下使用。

# 参数

|参数名|描述|可选参数|
|------|------|----|
|Bucket|空间名称，可以为私有空间或者公开空间名称|N|
|Prefix|七牛空间中文件名的前缀，该参数为可选参数，如果不指定则获取空间中所有的文件列表|Y|
|ListBucketResultFile|获取的文件列表保存在本地的文件名，如果不指定该参数，则会把结果输出到终端，一般可用于获取小规模文件列表测试使用|Y|

# 示例

1.获取空间`if-pbl`里面的所有文件列表：

```
qshell listbucket2 if-pbl -o if-pbl.list.txt
```

2.获取空间`if-pbl`里面的以`2014/10/07/`为前缀的文件列表：

```
qshell listbucket if-pbl --prefix '2014/10/07/' -o if-pbl.prefix.list.txt
```

结果：

```
hello.jpg	1710619	FlUqUK7zqbqm3NPwzq2q7TMZ-Ijs	14209629320769140	image/jpeg  1
hello.mp4	8495868	lns2dAHvO0qYseZFgDn3UqZlMOi-	14207312835630132	video/mp4   0
hhh	1492031	FjiRl_U0AeSsVCHXscCGObKyMy8f	14200176147531840	image/jpeg  1
jemygraw.jpg	1900176	FtmHAbztWfPEqPMv4t4vMNRYMETK	14208960018750329	application/octet-stream	1   QiniuAndroid
```
