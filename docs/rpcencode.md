# 简介
`rpcencode` 命令是通过 qiniu rpc 方式对数据进行编码。

# 格式
```
qshell rpcencode <DataToEncode1> [<DataToEncode2> [...]] [flags]
```

# 参数
- Data：待编码的数据 【必须】

# 示例
```
$ qshell rpcencode "https://qiniu.com/rpc?a=1&b=1"
https:!!qiniu.com!rpc'3Fa=1&b=1
```