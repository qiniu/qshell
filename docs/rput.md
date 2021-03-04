# 简介

`rput`命令使用七牛支持的分片上传的方式来上传一个文件，一般文件大小较大的情况下，可以使用分片上传来有效地保证文件上传的成功。

参考文档：
[分片上传 V1](https://developer.qiniu.com/kodo/7443/shard-to-upload)
[分片上传 V2)](https://developer.qiniu.com/kodo/6364/multipartupload-interface)

# 格式

```
qshell rput [--overwrite] [--version2] [--mimetype <MimeType>] [--callback-urls <CallbackUrls>] [--callback-host <CallbackHost>] [--storage <StorageType> ] <Bucket> <Key> <LocalFile>
```

其中 `Overwrite`，`MimeType`，`StorageType` (0 -> 标准存储， 1 - 低频存储)参数可根据需要指定一个或者多个，参数顺序随意，程序会自动识别。

# 鉴权

需要在使用了`account`设置了`AccessKey`, `SecretKey`和`Name`的情况下使用。

# 参数

| 参数名称     | 描述                                             | 可选参数 |
|--------------|--------------------------------------------------|----------|
| Bucket       | 七牛空间名称，可以为公开空间或私有空间           | N        |
| Key          | 文件保存在七牛空间的名称                         | N        |
| LocalFile    | 本地文件的路径                                   | N        |
| Overwrite    | 是否覆盖空间已有文件，默认为`false`              | Y        |
| MimeType     | 指定文件的MimeType                               | Y        |
| StorageType  | 文件存储类型，默认为`0`(标准存储） `1`为低频存储 | Y        |
| CallbackUrls | 上传回调地址，可以指定多个地址， 以逗号分开      | Y        |
| CallbackHost     | 上传回调HOST, 必须和CallbackUrls一起指定 | Y        |
| version2 | 使用分片上传 API V2 进行上传，默认为`false`, 使用 V1 上传 | Y |


# 示例

1.上传本地文件`/Users/jemy/Documents/qiniu.mp4`到空间`if-pbl`里面。

```
$ qshell rput if-pbl qiniu.mp4 /Users/jemy/Documents/qiniu.mp4
```

输出：

```
Uploading /Users/jemy/Documents/qiniu.mp4 => if-pbl : qiniu.mp4 ...
Progress: 100.00%
Put file /Users/jemy/Documents/qiniu.mp4 => if-pbl : qiniu.mp4 success!
Hash: lhsawSRA9-0L8b0s-cXmojaMhGqn
Fsize: 25538648 ( 24.36 MB )
MimeType: video/mp4
Last time: 16.82 s, Average Speed: 1517.9 KB/s
```

2.上传本地文件`/Users/jemy/Documents/qiniu.mp4`到空间`if-pbl`里面，带前缀`2015/01/18/`，并且指定`MimeType`参数为`video/mp4`

```
$ qshell rput if-pbl 2015/01/18/qiniu.mp4 /Users/jemy/Documents/qiniu.mp4 --mimetype video/mp4
```
输出：
```
Uploading /Users/jemy/Documents/qiniu.mp4 => if-pbl : 2015/01/18/qiniu.mp4 ...
Progress: 100%
Put file /Users/jemy/Documents/qiniu.mp4 => if-pbl : 2015/01/18/qiniu.mp4 success!
Hash: lhsawSRA9-0L8b0s-cXmojaMhGqn
Fsize: 25538648 ( 24.36 MB )
MimeType: video/mp4
Last time: 17.55 s, Average Speed: 1454.9 KB/s
```

3.覆盖上传 `qiniu.mp4' 到空间`if-pbl`

```
$ qshell rput if-pbl 2015/01/18/qiniu.jpg /Users/jemy/Documents/qiniu.jpg --overwrite
```

输出：

```
Uploading /Users/jemy/Documents/qiniu.jpg => if-pbl : 2015/01/18/qiniu.jpg ...
Progress: 100%
Put file /Users/jemy/Documents/qiniu.jpg => if-pbl : 2015/01/18/qiniu.jpg success!
Hash: Ftgm-CkWePC9fzMBTRNmPMhGBcSV
Fsize: 39335 ( 38.41 KB )
MimeType: image/jpeg
Last time: 0.40 s, Average Speed: 98.2 KB/s
```


4. 使用低频存储

```
$ qshell rput if-pbl 2015/01/18/qiniu.jpg /Users/jemy/Documents/qiniu.jpg --storage 1
```

输出：

```
Uploading /Users/jemy/Documents/qiniu.jpg => if-pbl : 2015/01/18/qiniu.jpg ...
Progress: 100%
Put file /Users/jemy/Documents/qiniu.jpg => if-pbl : 2015/01/18/qiniu.jpg success!
Hash: Ftgm-CkWePC9fzMBTRNmPMhGBcSV
Fsize: 39335 ( 38.41 KB )
MimeType: image/jpeg
Last time: 0.40 s, Average Speed: 98.2 KB/s
```
