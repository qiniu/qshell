# 简介
`urlencode` 该命令用来对一个字符串进行 URL 编码。

# 格式
```
qshell urlencode <DataToEncode>
```

# 帮助文档
可以在命令行输入如下命令获取帮助文档：
```
// 简单描述
$ qshell urlencode -h 

// 详细文档（此文档）
$ qshell urlencode --doc
```

# 鉴权
无

# 参数
- DataToEncode：待编码的字符串。 【必须】

# 示例
```
$ qshell urlencode 大数据时代
%E5%A4%A7%E6%95%B0%E6%8D%AE%E6%97%B6%E4%BB%A3
```