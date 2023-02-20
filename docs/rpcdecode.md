# 简介
`rpcdecode` 命令是对通过 qiniu rpc 方式 encode 的数据进行解码。

# 格式
```
qshell rpcdecode [DataToEncode...] [flags]
```

# 帮助文档
可以在命令行输入如下命令获取帮助文档：
```
// 简单描述
$ qshell rpcdecode -h 

// 详细文档（此文档）
$ qshell rpcdecode --doc
```

# 鉴权
需要使用 `qshell account` 或者 `qshell user add` 命令设置鉴权信息 `AccessKey`, `SecretKey` 和 `Name`。

# 参数
- DataToEncode：待解码的数据，当不指定则从 stdin 读取，每读取一行即进行编码并输出编码结果。

# 示例
```
$ qshell rpcdecode "https:\!\!qiniu.com\!rpc'3Fa=1&b=1"
https://qiniu.com/rpc?a=1&b=1
```
