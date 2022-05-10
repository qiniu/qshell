# 简介
`b64encode` 命令用来将一段字符串以 `Base64编码` 或 `URL安全的Base64编码` 格式进行编码。

# 格式
```
qshell b64encode [-s|--s] <DataToEncode>
```

# 参数
- DataToDecode：待编码字符串。【必选】

# 选项
- -s/--safe：标志开启 urlsafe 的 base64 编码。【可选】

# 示例
```
$ qshell b64encode 'hello world'
aGVsbG8gd29ybGQ=
```
