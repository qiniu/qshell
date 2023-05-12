# 简介
`stat` 指令根据七牛的公开API [stat](http://developer.qiniu.com/code/v6/api/kodo-api/rs/stat.html) 来获取空间中的一个文件的基本信息，包括文件的名称，保存的时间，hash值，文件的大小和MimeType。

参考文档：[资源元信息查询 (stat)](http://developer.qiniu.com/code/v6/api/kodo-api/rs/stat.html)

# 格式
```
qshell stat <Bucket> <Key>
```

# 帮助文档
可以在命令行输入如下命令获取帮助文档：
```
// 简单描述
$ qshell stat -h 

// 详细文档（此文档）
$ qshell stat --doc
```

# 鉴权
需要使用 `qshell account` 或者 `qshell user add` 命令设置鉴权信息 `AccessKey`, `SecretKey` 和 `Name`。

# 参数
- Bucket：空间名，可以为公开空间或者私有空间。【必须】
- Key：空间中的文件名。【必须】

# 示例
获取空间 `if-pbl` 中文件 `qiniu.png` 的基本信息
```
$ qshell stat if-pbl qiniu.png
```

输出：
```
Bucket:                  if-pbl
Key:                     qiniu.png
Etag:                    lozgLP_MAdAKZkPCXGvfd0LIDSUI
MD5:                     689b5cea4734143964a62214178f3f57
Fsize:                   5444314 -> 5.19MB
PutTime:                 16768889367943931 -> 2023-02-20 18:28:56.7943931 +0800 CST
MimeType:                text/plain
Status:                  0 -> 未禁用
Expiration:              not set
TransitionToIA:          not set
TransitionToArchive:     not set
TransitionToDeepArchive: not set
FileType:                1 -> 低频存储
```