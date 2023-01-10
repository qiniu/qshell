# 简介
`chgm` 指令用来为空间中的一个文件修改 MimeType。

参考文档：[资源元信息修改 (chgm)](http://developer.qiniu.com/code/v6/api/kodo-api/rs/chgm.html)

# 格式
```
qshell chgm <Bucket> <Key> <NewMimeType>
```

# 帮助文档
可以在命令行输入如下命令获取帮助文档：
```
// 简单描述
$ qshell chgm -h 

// 详细文档（此文档）
$ qshell chgm --doc
```

# 鉴权
需要使用 `qshell account` 或者 `qshell user add` 命令设置鉴权信息 `AccessKey`, `SecretKey` 和 `Name`。

# 参数
- Bucket：空间名，可以为公开空间或私有空间。【必须】
- Key：空间中的文件名。【必须】
- NewMimeType：给文件指定的新的 MimeType 。【必须】

# 示例
修改 `if-pbl` 空间中 `qiniu.png` 图片的MimeType为 `image/jpeg`
```
$ qshell chgm if-pbl qiniu.png image/jpeg
```

修改完成，我们检查一下文件的 MimeType：
```
$ qshell stat if-pbl qiniu.png
```

输出
```
Bucket:             if-pbl
Key:                qiniu.png
Hash:               FrUHIqhkDDd77-AtiDcOwi94YIeM
Fsize:              5331
PutTime:            14285516077733591
MimeType:           image/jpeg
```
我们发现，文件的 MimeType 已经被修改为 `image/jpeg`。
