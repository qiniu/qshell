# 简介

`listbucket`用来获取七牛空间里面的文件列表，可以指定文件前缀获取指定的文件列表，如果不指定，则获取所有文件的列表。

获取的文件列表组织格式如下：

```
Key\tSize\tHash\tPutTime\tMimeType\tEndUser
```

# 格式

```
qshell listbucket <Bucket> [<Prefix>] <ListBucketResultFile>
```

# 鉴权

需要在使用了`account`设置了`AccessKey`和`SecretKey`的情况下使用。

# 参数

|参数名|描述|可选参数|
|------|------|----|
|Bucket|空间名称，可以为私有空间或者公开空间名称|N|
|Prefix|七牛空间中文件名的前缀，该参数为可选参数，如果不指定则获取空间中所有的文件列表|Y|
|ListBucketResultFile|获取的文件列表保存在本地的文件名，如果该参数指定为`stdout`，则会把结果输出到终端，一般可用于获取小规模文件列表测试使用|N|

# 示例

1.获取空间`if-pbl`里面的所有文件列表：

```
qshell listbucket if-pbl if-pbl.list.txt
```

2.获取空间`if-pbl`里面的以`2014/10/07/`为前缀的文件列表：

```
qshell listbucket if-pbl '2014/10/07/' if-pbl.prefix.list.txt
```

结果：

```
hello.jpg	1710619	FlUqUK7zqbqm3NPwzq2q7TMZ-Ijs	14209629320769140	image/jpeg
hello.mp4	8495868	lns2dAHvO0qYseZFgDn3UqZlMOi-	14207312835630132	video/mp4
hhh	1492031	FjiRl_U0AeSsVCHXscCGObKyMy8f	14200176147531840	image/jpeg
jemygraw.jpg	1900176	FtmHAbztWfPEqPMv4t4vMNRYMETK	14208960018750329	application/octet-stream	QiniuAndroid
```

