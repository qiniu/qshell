# 简介
`tms2d` 该命令将一个以毫秒(ms)为单位的 Unix 时间戳转换为日期。

# 格式
```
qshell tms2d <TimestampInMilliSeconds>
```

# 帮助文档
可以在命令行输入如下命令获取帮助文档：
```
// 简单描述
$ qshell tms2d -h 

// 详细文档（此文档）
$ qshell tms2d --doc
```

# 鉴权
无

# 参数
- TimestampInMilliSeconds：以毫秒（ms）为单位的 Unix 时间戳。【必选】

# 示例
```
$ qshell tms2d 1427252311000
2015-03-25 10:58:31 +0800 CST
```