# 简介

`listbucket2`用来获取七牛空间里面的文件列表，可以指定文件前缀获取指定的文件列表，如果不指定，则获取所有文件的列表。

获取的文件列表组织格式如下（每个字段用Tab分隔）：

```
Key\tSize\tHash\tPutTime\tMimeType\tFileType\tEndUser
```


参考文档：[资源列举 (list)](http://developer.qiniu.com/code/v6/api/kodo-api/rs/list.html)


# 格式

```
qshell listbucket2 [--prefix <Prefix> | --suffixes <suffixes1,suffixes2>] [--start <StartDate>] [--max-retry <RetryCount>][--end <EndDate>] <Bucket> [ [-a] -o <ListBucketResultFile>]
```

# 鉴权

需要在使用了`account`设置了`AccessKey`, `SecretKey` 和  `Name` 的情况下使用。


# 选项和参数

| 参数名               | 描述                                                                                                           | 可选参数 |
|----------------------|----------------------------------------------------------------------------------------------------------------|----------|
| Bucket               | 空间名称，可以为私有空间或者公开空间名称                                                                       | N        |
| Prefix               | 七牛空间中文件名的前缀，该参数为可选参数，如果不指定则获取空间中所有的文件列表                                 | Y        |
| ListBucketResultFile | 获取的文件列表保存在本地的文件名，如果不指定该参数，则会把结果输出到终端，一般可用于获取小规模文件列表测试使用 | Y        |
| StartDate            | 列举整个空间，然后从中筛选出文件上传日期在<StartDate>之后的文件                                                | Y        |
| EndDate              | 列举整个空间， 然后从中筛选出文件上传日期在<EndDate>之前的文件                                                 | Y        |
| RetryCount           | 列举整个空间文件出错以后，最大的尝试次数；超过最大尝试次数以后，程序退出，打印出marker                         | Y        |
| suffixes             | 列举整个空间文件， 然后从中筛选出文件后缀为在[suffixes1, suffixes2, ...]中的文件                               | Y        |
| a                    | 开启选项o 的append模式， 如果本地保存文件列表的文件已经存在，如果希望像该文件添加内容，使用该选项, 必须和-o选项一起使用   | Y        |


# 常用使用场景介绍

（1）获取空间中所有的文件列表，这种情况下，可以直接指定 `Bucket` 参数和结果保存文件参数 `ListBucketResultFile` 即可。

```
qshell listbucket2 <Bucket> -o <ListBucketResultFile>
```
 
 (2) 如果本地文件`ListBucketResultFile`已经存在，有上一次列举的内容，如果希望把新的列表添加到该文件中，需要使用选项-a开启-o选项的append 模式
 
 ```
 qshell listbucket2 <Bucket> -a -o <ListBucketResultFile>
 ```

 (3) 获取空间所有文件，输出到屏幕上(标准输出)

 ```
 qshell listbucket2 <Bucket> 
 ```

（4）获取空间中指定前缀的文件列表

```
qshell listbucket2 [--prefix <Prefix>] <Bucket> -o <ListBucketResultFile>
```

 (5) 获取空间中指定前缀的文件列表， 输出到屏幕上
 
 ```
 qshell listbucket2 [--prefix <Prefix>] <Bucket>
 ```
 
 (6) 获取2018-10-30到2018-11-02上传的文件
 ```
 qshell listbucket2 --start 2018-10-30 --end 2018-11-02 <Bucket>
 ```
 
 (7) 获取后缀为mp4, html的文件
 
 ```
 qshell listbucket2 --suffixes mp4,html <Bucket>
 ```


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

# FAQ
1. 为什么列举空间很慢，很长时间没有打印出数据?
可能是您最近删除了大量的数据，在这种情况下接口会返回空的数据，这种数据不会打印出来，所以终端看到的感觉很慢

2. 为什么加了startDate 或者 endDate 或者 suffixes之后， 列举空间很慢，很长时间终端没有数据打印出来?
只有prefix选项是在后台提供的，其他的选项是方便用户筛选在命令行加的选项，所以实际上是列举整个空间，然后在这些文件中一个一个筛选符合条件的文件。
因此，如果您有100亿个文件，相当于把这100亿个文件先列举出来，然后逐一筛选。
