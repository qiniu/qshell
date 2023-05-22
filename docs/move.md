# 简介
`move` 命令可以将一个空间中的文件移动到另外一个空间中，也可以对同一空间中的文件重命名。移动文件仅支持在同一个帐号下面的同区域空间中移动。
注意如果目标文件已存在空间中的时候，默认情况下，`move` 会失败，报错 `614 file exists`，如果一定要强制覆盖目标文件，可以使用选项 `--overwrite` 。

参考文档：[资源移动／重命名 (move)](http://developer.qiniu.com/code/v6/api/kodo-api/rs/move.html)

# 格式
```
qshell move [--overwrite] <SrcBucket> <SrcKey> <DestBucket> [-k DestKey]
```

# 帮助文档
可以在命令行输入如下命令获取帮助文档：
```
// 简单描述
$ qshell move -h 

// 详细文档（此文档）
$ qshell move --doc
```

# 鉴权
需要使用 `qshell account` 或者 `qshell user add` 命令设置鉴权信息 `AccessKey`, `SecretKey` 和 `Name`。

# 参数
- SrcBucket: 源空间名称
- SrcKey: 源文件名称
- DestBucket: 目标空间名称

# 选项
- -k/--key: 目标文件名称(DestKey)，如果是 `DestBucket` 和 `SrcBucket` 不同的情况下，这个参数可以不填，默认和 `SrcKey` 相同。【可选】
- --overwrite: 当保存的文件已存在时，强制用新文件覆盖原文件，如果无此选项操作会失败。【可选】

# 示例
1 将空间 `if-pbl` 中的 `qiniu.jpg` 移动到 `if-pri` 中
```
qshell move if-pbl qiniu.jpg if-pri qiniu.jpg
```

2 将空间 `if-pbl` 中的 `qiniu.jpg` 重命名为 `2015/01/19/qiniu.jpg`
```
qshell move if-pbl qiniu.jpg if-pbl 2015/01/19/qiniu.jpg
```

3 将空间 `if-pbl` 中的 `qiniu.jpg` 移动到 `if-pri` 中，并命名为 `2015/01/19/qiniu.jpg`
```
qshell move if-pbl qiniu.jpg if-pri 2015/01/19/qiniu.jpg
```

4 强制覆盖 `if-pbl` 中的已有文件 `2015/01/19/qiniu.jpg`
```
qshell move --overwrite if-pbl qiniu.jpg if-pbl 2015/01/19/qiniu.jpg
```
执行命令之后，此时空间 `if-pbl` 里面的 `qiniu.jpg` 文件内容覆盖空间 `if-pbl` 里面的 `2015/01/19/qiniu.jpg`，`2015/01/19/qiniu.jpg` 文件原有内容完全被`qiniu.jpg` 文件覆盖，即空间 `if-pbl` 里面的 `qiniu.jpg` 文件此后已不存在，最后剩下 `2015/01/19/qiniu.jpg` 文件，文件内容是 `qiniu.jpg` 文件的内容。可以简单理解为鸠占鹊巢。
