# 简介
`rput` 命令使用七牛支持的分片上传的方式来上传一个文件，一般文件大小较大的情况下，可以使用分片上传来有效地保证文件上传的成功。

参考文档：
- [分片上传 V1](https://developer.qiniu.com/kodo/7443/shard-to-upload)
- [分片上传 V2](https://developer.qiniu.com/kodo/6364/multipartupload-interface)

# 格式
```
qshell rput [--overwrite] [--v2] [--mimetype <MimeType>] [--callback-urls <CallbackUrls>] [--callback-host <CallbackHost>] [--file-type <FileType> ] <Bucket> <Key> <LocalFile>
```

其中 `Overwrite`，`MimeType`，`FileType` (0: 标准存储， 1: 低频存储， 2: 归档存储， 3: 深度归档存储， 4: 归档直读存储)参数可根据需要指定一个或者多个，参数顺序随意，程序会自动识别。

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
- --accelerate：启用上传加速。【可选】
- --overwrite：是否覆盖空间已有文件，默认为 `false`。 【可选】
- -t/--mimetype：指定文件的 MimeType 。【可选】
- --file-type：文件存储类型；0: 标准存储， 1: 低频存储， 2: 归档存储， 3: 深度归档存储， 4: 归档直读存储；默认为`0`(标准存储）。 【可选】
- --resumable-api-v2：使用分片上传 API V2 进行上传，默认为 `false`, 使用 V1 上传。【可选】
- --resumable-api-v2-part-size：使用分片上传 API V2 进行上传时的分片大小，默认为 4M 。【可选】
- --sequential-read-file: 文件读为顺序读，不涉及跳读；开启后，上传中的分片数据会被加载至内存。此选项可能会增加挂载网络文件系统的文件上传速度。默认是：false。【可选】
- -l/--callback-urls：上传回调地址，可以指定多个地址，以逗号分开。【可选】
- -T/--callback-host：上传回调HOST, 必须和 CallbackUrls 一起指定。 【可选】
-    --callback-body：上传成功后，七牛云向业务服务器发送 Content-Type: application/x-www-form-urlencoded 的 POST 请求。业务服务器可以通过直接读取请求的 query 来获得该字段，支持魔法变量和自定义变量。callbackBody 要求是合法的 url query string。例如key=$(key)&hash=$(etag)&w=$(imageInfo.width)&h=$(imageInfo.height)。如果callbackBodyType指定为application/json，则callbackBody应为json格式，例如:{“key”:"$(key)",“hash”:"$(etag)",“w”:"$(imageInfo.width)",“h”:"$(imageInfo.height)"}。【可选】
-    --callback-body-type：上传成功后，七牛云向业务服务器发送回调通知 callbackBody 的 Content-Type。默认为 application/x-www-form-urlencoded，也可设置为 application/json。【可选】
-    --end-user：上传成功后，七牛云向业务服务器发送回调通知 callbackBody 的 Content-Type。默认为 application/x-www-form-urlencoded，也可设置为 application/json。【可选】
-    --persistent-ops：资源上传成功后触发执行的预转持久化处理指令列表。fileType=2或3（上传归档存储或深度归档存储文件）时，不支持使用该参数。支持魔法变量和自定义变量。每个指令是一个 API 规格字符串，多个指令用;分隔。【可选】
-    --persistent-notify-url：接收持久化处理结果通知的 URL。必须是公网上可以正常进行 POST 请求并能成功响应的有效 URL。该 URL 获取的内容和持久化处理状态查询的处理结果一致。发送 body 格式是 Content-Type 为 application/json 的 POST 请求，需要按照读取流的形式读取请求的 body 才能获取。【可选】
-    --persistent-pipeline：转码队列名。资源上传成功后，触发转码时指定独立的队列进行转码。为空则表示使用公用队列，处理速度比较慢。建议使用专用队列。【可选】
-    --detect-mime：开启 MimeType 侦测功能，并按照下述规则进行侦测；如不能侦测出正确的值，会默认使用 application/octet-stream 。【可选】
```
    1. 设为 1 值，则忽略上传端传递的文件 MimeType 信息，并按如下顺序侦测 MimeType 值：
        1) 侦测内容；
        2) 检查文件扩展名；
        3) 检查 Key 扩展名。
    2. 默认设为 0 值，如上传端指定了 MimeType（application/octet-stream 除外）则直接使用该值，否则按如下顺序侦测 MimeType 值：
        1) 检查文件扩展名；
        2) 检查 Key 扩展名；
        3) 侦测内容。
    3. 设为 -1 值，无论上传端指定了何值直接使用该值。
```
-    --traffic-limit：上传请求单链接速度限制，控制客户端带宽占用。限速值取值范围为 819200 ~ 838860800，单位为 bit/s。【可选】

# 示例
1 上传本地文件 `/Users/jemy/Documents/qiniu.mp4` 到空间 `if-pbl` 里面。
```
// 使用使用分片上传 API V1
$ qshell rput if-pbl qiniu.mp4 /Users/jemy/Documents/qiniu.mp4

// 使用使用分片上传 API V2
$ qshell rput if-pbl qiniu.mp4 /Users/jemy/Documents/qiniu.mp4 --resumable-api-v2
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
