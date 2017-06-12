# 简介

`listbucket`用来获取七牛空间里面的文件列表，可以指定文件前缀获取指定的文件列表，如果不指定，则获取所有文件的列表。

获取的文件列表组织格式如下（每个字段用Tab分隔）：

```
Key\tSize\tHash\tPutTime\tMimeType\tFileType\tEndUser
```


参考文档：[资源列举 (list)](http://developer.qiniu.com/code/v6/api/kodo-api/rs/list.html)

# 格式

```
qshell listbucket [-marker <Marker>] <Bucket> [<Prefix>] <ListBucketResultFile>
```

上面的命令中，可选的场景有三种：

（1）获取空间中所有的文件列表，这种情况下，可以直接指定 `Bucket` 参数和结果保存文件参数 `ListBucketResultFile` 即可。

```
qshell listbucket <Bucket> <ListBucketResultFile>
```

（2）获取空间中指定前缀的文件列表，这种情况下，除了指定（1）中的参数外，还需要指定 `Prefix` 参数。

```
qshell listbucket <Bucket> <Prefix> <ListBucketResultFile>
```

（3）该场景主要用在空间中文件列表较多导致大量列举操作超时或者是列举过程中网络异常导致列举操作失败的时候，这个时候列举失败的时候，程序会输出当时失败的`marker`，如果我们希望接着上一次的列举进度继续列举，那么可以在运行命令的时候，额外指定选项`marker`。

```
qshell listbucket -marker <Marker> <Bucket> <ListBucketResultFile>
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
hello.jpg	1710619	FlUqUK7zqbqm3NPwzq2q7TMZ-Ijs	14209629320769140	image/jpeg  1
hello.mp4	8495868	lns2dAHvO0qYseZFgDn3UqZlMOi-	14207312835630132	video/mp4   0
hhh	1492031	FjiRl_U0AeSsVCHXscCGObKyMy8f	14200176147531840	image/jpeg  1
jemygraw.jpg	1900176	FtmHAbztWfPEqPMv4t4vMNRYMETK	14208960018750329	application/octet-stream	1   QiniuAndroid
```

