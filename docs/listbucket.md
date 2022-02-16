# 简介
`listbucket` 用来获取七牛空间里面的文件列表，可以指定文件前缀获取指定的文件列表，如果不指定，则获取所有文件的列表。

获取的文件列表组织格式如下（每个字段用Tab分隔）：
```
Key\tSize\tHash\tPutTime\tMimeType\tFileType\tEndUser
```

参考文档：[资源列举 (list)](http://developer.qiniu.com/code/v6/api/kodo-api/rs/list.html)

备注：有个优化版本的命令叫 `listbucket2` 功能描述和这个命令一样，但是更加适合海量文件的空间列举。

# 格式
```
qshell listbucket [--prefix <Prefix>] <Bucket> [-o <ListBucketResultFile>]
```

1 获取空间中所有的文件列表，这种情况下，可以直接指定 `Bucket` 参数和结果保存文件参数 `ListBucketResultFile` 即可。
```
qshell listbucket <Bucket> -o <ListBucketResultFile>
```

2 获取空间所有文件，输出到屏幕上(标准输出)
 ```
 qshell listbucket <Bucket> 
 ```

3 获取空间中指定前缀的文件列表
```
qshell listbucket [--prefix <Prefix>] <Bucket> -o <ListBucketResultFile>
```

4 获取空间中指定前缀的文件列表， 输出到屏幕上
 ```
 qshell listbucket [--prefix <Prefix>] <Bucket>
 ```

# 鉴权
需要在使用了 `account` 设置了 `AccessKey`, `SecretKey` 和  `Name` 的情况下使用。

# 参数
- Bucket：空间名称，可以为私有空间或者公开空间名称 【必选】

# 选项
- --prefix：七牛空间中文件名的前缀，该参数为可选参数，如果不指定则获取空间中所有的文件列表 【可选】
- --out：获取的文件列表保存在本地的文件名，如果不指定该参数，则会把结果输出到终端，一般可用于获取小规模文件列表测试使用 【可选】

# 示例
1 获取空间`if-pbl`里面的所有文件列表：
```
qshell listbucket if-pbl -o if-pbl.list.txt
```

2 获取空间`if-pbl`里面的以`2014/10/07/`为前缀的文件列表：
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

