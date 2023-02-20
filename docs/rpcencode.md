# 简介
`rpcencode` 命令是通过 qiniu rpc 方式对数据进行编码。

# 格式
```
qshell rpcencode <DataToEncode1> [<DataToEncode2> [...]] [flags]
```

# 帮助文档
可以在命令行输入如下命令获取帮助文档：
```
// 简单描述
$ qshell rpcencode -h 

// 详细文档（此文档）
$ qshell rpcencode --doc
```

# 鉴权
需要使用 `qshell account` 或者 `qshell user add` 命令设置鉴权信息 `AccessKey`, `SecretKey` 和 `Name`。

# 参数
- Data：待编码的数据 【必须】

# 示例
```
$ qshell rpcencode "https://qiniu.com/rpc?a=1&b=1"
https:!!qiniu.com!rpc'3Fa=1&b=1
```