# 简介
`rename` 命令可以对一个空间中的文件进行重命名。
注意如果目标文件已存在空间中的时候，默认情况下，`rename` 会失败，报错 `614 file exists`，如果一定要强制覆盖目标文件，可以使用选项 `--overwrite` 。

参考文档：[资源移动／重命名 (move)](http://developer.qiniu.com/code/v6/api/kodo-api/rs/move.html)

# 格式
```
qshell rename [--overwrite] <Bucket> <SrcKey> <DestKey>
```

# 帮助文档
可以在命令行输入如下命令获取帮助文档：
```
// 简单描述
$ qshell rename -h 

// 详细文档（此文档）
$ qshell rename --doc
```

# 鉴权
需要使用 `qshell account` 或者 `qshell user add` 命令设置鉴权信息 `AccessKey`, `SecretKey` 和 `Name`。

# 参数
- Bucket: 空间名 【必须】
- SrcKey: 原文件名称 【必须】
- DestKey: 目标文件名称 【必须】

# 示例
1 将空间 `if-pbl` 中的 `qiniu.jpg` 重命名为 `qiniu_new.jpg`
```
qshell rename if-pbl qiniu.jpg qiniu_new.jpg
```

2 将空间 `if-pbl` 中的 `qiniu.jpg` 重命名为 `2015/01/19/qiniu.jpg`
```
qshell rename if-pbl qiniu.jpg 2015/01/19/qiniu.jpg
```

4 强制将空间 `if-pbl` 中的 `qiniu.jpg` 重名名为 `2015/01/19/qiniu.jpg`
```
qshell rename --overwrite if-pbl qiniu.jpg 2015/01/19/qiniu.jpg
```
执行命令之后，此时空间 `if-pbl` 里面的 `qiniu.jpg` 文件内容覆盖空间里面原名为 `2015/01/19/qiniu.jpg`的文件，`2015/01/19/qiniu.jpg` 文件原有内容完全被`qiniu.jpg` 文件覆盖，即空间 `if-pbl` 里面的 `qiniu.jpg` 文件此后已不存在，最后剩下 `2015/01/19/qiniu.jpg` 文件，文件内容是 `qiniu.jpg` 文件的内容。可以简单理解为鸠占鹊巢。