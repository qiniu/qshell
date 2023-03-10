# 简介
`rput` 命令使用七牛支持的分片上传的方式来上传一个文件，一般文件大小较大的情况下，可以使用分片上传来有效地保证文件上传的成功。

参考文档：
- [分片上传 V1](https://developer.qiniu.com/kodo/7443/shard-to-upload)
- [分片上传 V2](https://developer.qiniu.com/kodo/6364/multipartupload-interface)

# 格式
```
qshell rput [--overwrite] [--v2] [--mimetype <MimeType>] [--callback-urls <CallbackUrls>] [--callback-host <CallbackHost>] [--file-type <FileType> ] <Bucket> <Key> <LocalFile>
```

其中 `Overwrite`，`MimeType`，`FileType` (0: 标准存储， 1: 低频存储， 2: 归档存储， 3: 深度归档存储)参数可根据需要指定一个或者多个，参数顺序随意，程序会自动识别。

# 帮助文档
可以在命令行输入如下命令获取帮助文档：
```
// 简单描述
$ qshell rput -h 

// 详细文档（此文档）
$ qshell rput --doc
```

# 鉴权
需要使用 `qshell account` 或者 `qshell user add` 命令设置鉴权信息 `AccessKey`, `SecretKey` 和 `Name`。

# 参数
- Bucket：空间名，可以为公开空间或私有空间。 【必选】
- Key: 文件保存在七牛空间的名称。 【必选】
- LocalFile：本地文件的路径。 【必选】

# 选项
- --overwrite：是否覆盖空间已有文件，默认为 `false`。 【可选】
- -t/--mimetype：指定文件的 MimeType 。【可选】
- --file-type：文件存储类型；0: 标准存储， 1: 低频存储， 2: 归档存储， 3: 深度归档存储；默认为`0`(标准存储）。 【可选】
- -l/--callback-urls：上传回调地址，可以指定多个地址，以逗号分开。【可选】
- -T/--callbackHost：上传回调HOST, 必须和 CallbackUrls 一起指定。 【可选】
- --resumable-api-v2：使用分片上传 API V2 进行上传，默认为 `false`, 使用 V1 上传。【可选】
- --resumable-api-v2-part-size：使用分片上传 API V2 进行上传时的分片大小，默认为 4M 。【可选】

# 示例
1 上传本地文件 `/Users/jemy/Documents/qiniu.mp4` 到空间 `if-pbl` 里面。
```
// 使用使用分片上传 API V1
$ qshell rput if-pbl qiniu.mp4 /Users/jemy/Documents/qiniu.mp4

// 使用使用分片上传 API V2
$ qshell rput if-pbl qiniu.mp4 /Users/jemy/Documents/qiniu.mp4 --v2
```

2 上传本地文件 `/Users/jemy/Documents/qiniu.mp4` 到空间 `if-pbl` 里面，带前缀 `2015/01/18/`，并且指定 `MimeType` 参数为 `video/mp4`。
```
$ qshell rput if-pbl 2015/01/18/qiniu.mp4 /Users/jemy/Documents/qiniu.mp4 --mimetype video/mp4
```

3 覆盖上传 `qiniu.mp4` 到空间 `if-pbl`
```
$ qshell rput if-pbl 2015/01/18/qiniu.jpg /Users/jemy/Documents/qiniu.jpg --overwrite
```

4 使用低频存储
```
$ qshell rput if-pbl 2015/01/18/qiniu.jpg /Users/jemy/Documents/qiniu.jpg --file-type 1
```
