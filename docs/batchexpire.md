# 简介
`batchexpire` 命令用来为空间中的文件设置过期时间。该操作发生在同一个空间中。（将文件设置为从现在开始xx天后自动删除的状态）

# 格式
```
qshell batchexpire [--force] [--sucess-list <SuccessFileName>] [--failure-list <FailureFileName>] [--sep <Separator>]  [--worker <WorkerCount>] <Bucket> <-i KeyDeleteAfterDaysMapFile>
```

# 帮助 
```
qshell batchexpire -h
```

# 鉴权
需要在使用了 `account` 设置了 `AccessKey`, `SecretKey` 和 `Name` 的情况下使用。

# 参数
- Bucket：空间名，可以为公开空间或私有空间。【必须】

# 选项
- -i：接受一个文件参数， 文件内容每行包含 `文件名` 和 `过期天数`，过期天数仅用数字表示即可。每行多个元素名之间用分割符分隔（默认 tab 制表符）； 如果需要自定义分割符，可以使用 `-F` 或 `--sep` 选项指定自定义的分隔符。如果没有通过该选项指定该文件参数， 从标准输入读取内容。每行具体格式如下：（【可选】）
```
<Key><Sep>1 // <Key>：文件名，<Sep>：分割符，1：过期天数。
```
- --force：该选项控制工具的默认行为。默认情况下，对于批量操作，工具会要求使用者输入一个验证码，确认下要进行批量文件操作了，避免操作失误的发生。如果不需要这个验证码的提示过程，可以使用 `--force` 选项。【可选】
- --success-list：该选项指定一个文件，qshell 会把操作成功的文件行导入到该文件；默认不导出。【可选】
- --failure-list：该选项指定一个文件， qshell 会把操作失败的文件行加上错误状态码，错误的原因导入该文件；默认不导出。【可选】
- --sep：该选项可以自定义输入内容（文件输入或标准输入，参考 -i 选项说明）；默认为 tab 制表符。【可选】
- --worker：该选项可以定义 Batch 任务并发数；默认为 1。【可选】

# 示例
1 比如我们要将空间 `if-pbl` 里面的一些文件改为3天后过期，我们可以指定如下的`KeyFileTypeMapFile` 的内容：
```
2015/03/22/qiniu.png	3
2015/photo.jpg	3
```

上面，我们将 `2015/03/22/qiniu.png` 文件设置为3天后过期了，诸如此类。
把这个内容保存到文件 `toexpire.txt` 中，然后使用如下的命令对 `toexpire.txt` 中的所有文件设置过期时间。
```
$ qshell batchexpire if-pbl -i toexpire.txt
```

2 如果不希望上面的重命名过程出现验证码提示，可以使用 `--force` 选项：
```
$ qshell batchexpire --force if-pbl -i toexpire.txt
```

# 注意
如果没有指定输入文件的话，会从标准输入读取内容。
