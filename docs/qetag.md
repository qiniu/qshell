# 简介
`qtag` 算法是七牛用来计算文件 `hash` 的自定义算法。这个命令用来根据 `qetag` 算法快速计算文件 `hash`。

参考文档：[qetag](https://github.com/qiniu/qetag)

# 格式
```
qshell qetag <LocalFilePath>
```

# 帮助文档
可以在命令行输入如下命令获取帮助文档：
```
// 简单描述
$ qshell qetag -h 

// 详细文档（此文档）
$ qshell qetag --doc
```

# 参数
- LocalFilePath：本地文件路径

# 示例
```
$ qshell qetag yyy.jpg
Fu9LtwRE8Q_iy4ITrFOGYvqfbifZ
```
