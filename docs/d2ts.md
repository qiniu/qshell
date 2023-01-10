# 简介
`d2ts` 命令用来生成一个 SecondsToNow 秒后的 Unix 时间戳（单位秒）。

# 格式
```
qshell d2ts <SecondsToNow>
```

# 帮助文档
可以在命令行输入如下命令获取帮助文档：
```
// 简单描述
$ qshell d2ts -h 

// 详细文档（此文档）
$ qshell d2ts --doc
```

# 鉴权
无

# 参数
- SecondsToNow: 指定的秒数

# 示例
```
$ qshell d2ts 3600
1427252311
```