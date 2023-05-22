# 简介
`tns2d` 因为七牛的 [stat接口](http://developer.qiniu.com/docs/v6/api/reference/rs/stat.html) 返回的 `putTime` 字段的单位是 `100纳秒`，有时候我们需要把它转出来看看。`tns2d` 这个命令就是这个作用。可以把 `putTime` 的值直接作为参数，得到日期结果。

# 格式
```
qshell tns2d <TimestampIn100NanoSeconds>
```

# 帮助文档
可以在命令行输入如下命令获取帮助文档：
```
// 简单描述
$ qshell tns2d -h 

// 详细文档（此文档）
$ qshell tns2d --doc
```

# 鉴权
无

# 参数
- TimestampIn100NanoSeconds：以 100 纳秒为单位的 Unix 时间戳。 【必选】

# 示例
```
$ qshell tns2d 13603956734587420
2013-02-09 15:41:13.458742 +0800 CST
```