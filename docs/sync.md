# 简介
`sync` 指令用来弥补 `fetch` 指令的不足之处。`fetch` 指令适合于中小文件的抓取，根据实际经验，基本上适合 `50MB` 以下的文件抓取。但是很多场合，大的文件，比如 1GB，100GB 的文件想要直接从服务器迁移过来，就不能使用 `fetch` 功能，这个时候可以使用 `sync` 指令。

`sync` 指令的基本原理是使用 `Range` 方式默认按照 `4MB` 一个块从资源服务器获取数据，然后使用七牛支持的分片上传功能直接传到七牛存储空间中。

另外 `sync` 指令在执行过程中，并不用担心网络中断导致的同步中断，因为采用了分片上传的机制，我们会把每一个成功上传的块的位置记录下来，当下次网络恢复的时候，只需要运行原始命令即可从断点处恢复。

注：如果 url 不支持 Range 则不可以 sync。

# 格式
```
qshell sync <SrcResUrl> <Bucket> [-k <Key>]
```

# 帮助文档
可以在命令行输入如下命令获取帮助文档：
```
// 简单描述
$ qshell sync -h

// 详细文档（此文档）
$ qshell sync --doc
```

# 鉴权
需要使用 `qshell account` 或者 `qshell user add` 命令设置鉴权信息 `AccessKey`, `SecretKey` 和 `Name`。

# 参数
- SrcResUrl：互联网上资源的链接，必须是可访问的链接。 【必选】
- Bucket：空间名，可以为公开空间或者私有空间。 【必选】

# 选项
- --accelerate：启用上传加速。【可选】
- -k/--key：该资源保存在空间中的 key，不配置时使用资源 Url 中文件名作为存储的 key。 【可选】
- -u/--uphost：上传入口的 IP 地址，一般在大文件的情况下，可以指定上传入口的 IP 来减少 DNS 环节，提升同步速度。 【可选】
- --file-type：文件存储类型，0:标准存储 1:低频存储 2:归档存储 3:深度归档 4:归档直读存储 5:智能分层存储；默认为 0【可选】
- --resumable-api-v2：使用分片 v2 进行上传；默认使用 v1。 【可选】
- --resumable-api-v2-part-size：使用分片上传 API V2 进行上传时的分片大小，默认为 4M 。【可选】
- --overwrite：是否覆盖空间已有文件，默认为 `false`。 【可选】
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


##### 备注：
上传入口的域名对应的 IP 地址一般情况下是不变的，减少 DNS 的查询环节，可以提升同步速度和稳定性。
上传入口的域名对应的 IP 地址可以通过如下的命令来获取解析的结果：
```
华东机房
$ dig up.qiniu.com

华北机房
$ dig up-z1.qiniu.com

华南机房
$ dig up-z2.qiniu.com

北美机房
$ dig up-na0.qiniu.com

东南亚机房
$ dig up-as0.qiniu.com
```

# 示例
使用分片 v2 抓取一个资源并以指定的文件名保存在七牛的空间里面：
```
$ qshell sync http://if-pbl.qiniudn.com/test_big_movie.mp4 if-pbl test.mp4 --resumable-api-v2
```
