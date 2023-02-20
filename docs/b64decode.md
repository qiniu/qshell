# 简介
`b64decode` 命令用来将一段以 `Base64 编码` 或 `URL 安全的 Base64 编码` 的字符串解码。

# 格式
```
qshell b64decode [-s|--safe] <DataToDecode>
```

# 帮助文档
可以在命令行输入如下命令获取帮助文档：
```
// 简单描述
$ qshell b64decode -h 

// 详细文档（此文档）
$ qshell b64decode --doc
```

# 鉴权
无

# 参数
- DataToDecode：待解码字符串。【必选】

# 选项
- -s/--safe：标志开启 urlsafe 的 base64 编码。【可选】

# 示例
我们可以解码七牛上传凭证的第三部分，即编码后的 PutPolicy：
```
$ qshell b64decode 'eyJzY29wZSI6ImJiaW1nOnRlc3QucG5nIiwiZGVhZGxpbmUiOjE0MjcxODkxMzB9'
{"scope":"bbimg:test.png","deadline":1427189130}
```
