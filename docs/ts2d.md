# 简介
`ts2d` 命令将一个以秒（s）为单位的 Unix 时间戳转换为日期。

# 格式
```
qshell ts2d <TimestampInSeconds>
```

# 帮助文档
可以在命令行输入如下命令获取帮助文档：
```
// 简单描述
$ qshell ts2d -h 

// 详细文档（此文档）
$ qshell ts2d --doc
```

# 鉴权
无

# 参数
- TimestampInSeconds：以秒（s）为单位的 Unix 时间戳。 【必须】

# 示例
```
$ qshell ts2d 1427252311
2015-03-25 10:58:31 +0800 CST
```