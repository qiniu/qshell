# 简介
`fput` 命令用来以 `multipart/form-data` 的表单方式上传一个文件。适合于中小型文件的上传，一般建议如果文件大小超过 100MB 的话，都使用分片上传。

参考文档：[直传文件 (upload)](http://developer.qiniu.com/code/v6/api/kodo-api/up/upload.html)

# 格式
```
qshell fput [--overwrite] [--callback-urls <CallbackUrls>] [--callback-host <CallbackHost>] [--file-type <FileType>] [--mimetype <MimeType>] <Bucket> <Key> <LocalFile>
```

其中 `Overwrite`，`MimeType`，`FileType` 参数可根据需要指定一个或者多个，参数顺序随意，程序会自动识别。

# 帮助文档
可以在命令行输入如下命令获取帮助文档：
```
// 简单描述
$ qshell fput -h 

// 详细文档（此文档）
$ qshell fput --doc
```

# 鉴权
需要使用 `qshell account` 或者 `qshell user add` 命令设置鉴权信息 `AccessKey`, `SecretKey` 和 `Name`。

# 参数
- Bucket：空间名，可以为公开空间或私有空间【必选】
- Key：文件保存在七牛空间的名称 【必选】
- LocalFile：本地文件的路径【必选】
  
# 选项
-    --overwrite：是否覆盖空间已有文件，默认为 `false`。 【可选】
- -t/--mimetype：指定文件的 MimeType。 【可选】
-    --file-type：文件存储类型，默认为 `0`（标准存储），`1` 为低频存储，`2` 为归档存储，`3` 为深度归档存储，`4` 为归档直读存储。 【可选】
- -u/--up-host: 指定上传域名。 【可选】
- -l/--callback-urls：上传回调地址， 可以指定多个地址，以逗号分隔。 【可选】
- -T/--callback-host：上传回调的 HOST, 必须和 CallbackUrls 一起指定。 【可选】

# 示例
1 上传本地文件 `/Users/jemy/Documents/qiniu.jpg` 到空间 `if-pbl` 里面。
```
$ qshell fput if-pbl qiniu.jpg /Users/jemy/Documents/qiniu.jpg
```

2 上传本地文件 `/Users/jemy/Documents/qiniu.jpg` 到空间 `if-pbl` 里面，带前缀 `2015/01/18/`，并且指定 `MimeType` 参数为 `image/jpg`
```
$ qshell fput if-pbl 2015/01/18/qiniu.jpg /Users/jemy/Documents/qiniu.jpg --mimetype image/jpg
```

3 覆盖上传 `qiniu.mp4` 到空间 `if-pbl`
```
$ qshell fput if-pbl 2015/01/18/qiniu.jpg /Users/jemy/Documents/qiniu.jpg --overwrite
```

5 使用低频存储
```
$ qshell fput if-pbl 2015/01/18/qiniu.jpg /Users/jemy/Documents/qiniu.jpg --file-type 1
```
