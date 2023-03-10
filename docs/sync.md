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
- -k/--key：该资源保存在空间中的 key，不配置时使用资源 Url 中文件名作为存储的 key。 【可选】
- -u/--uphost：上传入口的 IP 地址，一般在大文件的情况下，可以指定上传入口的 IP 来减少 DNS 环节，提升同步速度。 【可选】
- --file-type：文件存储类型，默认为 `0` (标准存储），`1` 为低频存储，`2` 为归档存储，`3` 为深度归档存储，【可选】
- --resumable-api-v2：使用分片 v2 进行上传；默认使用 v1。 【可选】
- --resumable-api-v2-part-size：使用分片上传 API V2 进行上传时的分片大小，默认为 4M 。【可选】
- --overwrite：是否覆盖空间已有文件，默认为 `false`。 【可选】

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
