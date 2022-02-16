# 简介
`listbucket2` 用来获取七牛空间里面的文件列表，可以指定文件前缀获取指定的文件列表，如果不指定，则获取所有文件的列表。

获取的文件列表组织格式如下（每个字段用 Tab 分隔）：
```
Key\tSize\tHash\tPutTime\tMimeType\tFileType\tEndUser
```

参考文档：[资源列举 (list)](http://developer.qiniu.com/code/v6/api/kodo-api/rs/list.html)

# 格式
```
qshell listbucket2 [-m|--marker <Marker>][--prefix <Prefix> | --suffixes <suffixes1,suffixes2>] [--start <StartDate>] [--max-retry <RetryCount>][--end <EndDate>] <Bucket> [--readable] [ [-a] -o <ListBucketResultFile>]
```

# 鉴权
需要在使用了 `account` 设置了 `AccessKey`, `SecretKey` 和  `Name` 的情况下使用。

# 参数
- Bucket： 空间名称，可以为私有空间或者公开空间名称 【必选】
  
# 选项
- --prefix： 七牛空间中文件名的前缀，该参数为可选参数，如果不指定则获取空间中所有的文件列表 【可选】
- --out： 获取的文件列表保存在本地的文件名，如果不指定该参数，则会把结果输出到终端，一般可用于获取小规模文件列表测试使用 【可选】
- --start： 列举整个空间，然后从中筛选出文件上传日期在 <StartDate> 之后的文件；格式：yyyy-mm-dd-hh-MM-ss eg:2022-01-10-08-30-20 【可选】
- --end： 列举整个空间， 然后从中筛选出文件上传日期在<EndDate>之前的文件；格式：yyyy-mm-dd-hh-MM-ss eg:2022-01-10-08-30-20 【可选】
- --max-retry： 列举整个空间文件出错以后，最大的尝试次数；超过最大尝试次数以后，程序退出，打印出marker 【可选】
- --suffixes： 列举整个空间文件， 然后从中筛选出文件后缀为在 [suffixes1, suffixes2, ...] 中的文件 【可选】
- --append： 开启选项 --out 的 append 模式， 如果本地保存文件列表的文件已经存在，如果希望像该文件添加内容，使用该选项, 必须和 --out 选项一起使用 【可选】
- --readable： 开启文件大小的可读性选项， 会以合适的 KB, MB, GB 等显示 【可选】
- --marker： marker 标记列举过程中的位置， 如果列举的过程中网络断开，会返回一个 marker, 可以指定该 marker 参数继续列举 【可选】


# 常用场景
1 获取空间中所有的文件列表，这种情况下，可以直接指定 `Bucket` 参数和结果保存文件参数 `ListBucketResultFile` 即可。
```
qshell listbucket2 <Bucket> -o <ListBucketResultFile>
```
 
2 如果本地文件 `ListBucketResultFile` 已经存在，有上一次列举的内容，如果希望把新的列表添加到该文件中，需要使用选项 -a 开启 -o 选项的 append 模式
 ```
 qshell listbucket2 <Bucket> -a -o <ListBucketResultFile>
 ```

3 获取空间所有文件，输出到屏幕上(标准输出)
 ```
 qshell listbucket2 <Bucket> 
 ```

4 获取空间中指定前缀的文件列表
```
qshell listbucket2 [--prefix <Prefix>] <Bucket> -o <ListBucketResultFile>
```

5 获取空间中指定前缀的文件列表， 输出到屏幕上
 ```
 qshell listbucket2 [--prefix <Prefix>] <Bucket>
 ```
 
6 获取 `2018-10-30` 到 `2018-11-02` 上传的文件
 ```
 qshell listbucket2 --start 2018-10-30 --end 2018-11-03 <Bucket>
 ```
注意startDate 和 endDate 是这种半开半闭区间[startDate, endDate)

7 获取后缀为 mp4, html 的文件
 ```
 qshell listbucket2 --suffixes mp4,html <Bucket>
 ```
 
8 通常列举的文件的大小都是以字节显示，如果想以人工可读的方式 B, KB, MB 等显示，可以使用 -r 或者 --readable 选项
 ```
 qshell listbucket2 -r <Bucket>
 ```

9 marker 的使用; 假如要列举的 bucket 名字为 "test-marker", marker 为"eyJjIjowLCJrIjoiMDkzOWM1ODU4ZmI1NGZiNzk3NTJmNjVkN2U4MWY4MmVfMTUzNTM3NzI2MDMxNV8xNTM1MzgwMjYyNDYxXzgzMjgyODAzOC0wMDAwMS5tcDQifQ=", 如果要接着这个 marker 位置继续列举，可以使用如下命令
 ```
 $ qshell listbucket2 -m eyJjIjowLCJrIjoiMDkzOWM1ODU4ZmI1NGZiNzk3NTJmNjVkN2U4MWY4MmVfMTUzNTM3NzI2MDMxNV8xNTM1MzgwMjYyNDYxXzgzMjgyODAzOC0wMDAwMS5tcDQifQ= test-marker
 ```


# 示例
1 获取空间 `if-pbl` 里面的所有文件列表：
```
qshell listbucket2 if-pbl -o if-pbl.list.txt
```

2 获取空间 `if-pbl` 里面的以 `2014/10/07/` 为前缀的文件列表：
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

2. 为什么加了 startDate 或者 endDate 或者 suffixes 之后， 列举空间很慢，很长时间终端没有数据打印出来?
只有 prefix 选项是在后台提供的，其他的选项是方便用户筛选在命令行加的选项，所以实际上是列举整个空间，然后在这些文件中一个一个筛选符合条件的文件。
因此，如果您有 100 亿个文件，相当于把这 100 亿个文件先列举出来，然后逐一筛选。
