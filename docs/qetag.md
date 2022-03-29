# 简介
`qtag` 算法是七牛用来计算文件hash的自定义算法。这个命令用来根据 `qetag` 算法快速计算文件 `hash`。

参考文档：[qetag](https://github.com/qiniu/qetag)

# 格式
```
qshell qetag <LocalFilePath>
```

# 参数
- LocalFilePath：本地文件路径

# 示例
```
$ qshell qetag yyy.jpg
Fu9LtwRE8Q_iy4ITrFOGYvqfbifZ
```
