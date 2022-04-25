# 简介
`batchchtype` 命令用来为空间中的文件设置存储类型。该操作发生在同一个空间中。（将文件设置为 **深度归档存储** 或者 **归档存储** 或者 **低频存储** 或者 **普通存储**，默认：文件为**普通存储**）

# 格式
```
qshell batchchtype  [--force] [--success-list <SuccessFileName>] [--failure-list <FailureFileName>] [--sep <Separator>]  [--worker <WorkerCount>] <Bucket> [-i <KeyFileTypeMapFile>]
```

# 鉴权
需要在使用了 `account` 设置了 `AccessKey`, `SecretKey` 和 `Name ` 的情况下使用。

# 帮助 
```
qshell batchchtype -h
```

# 参数
- Bucket：空间名，可以为公开空间或私有空间。【必选】

# 选项
- -i/--input-file：接受一个文件, 文件内容每行包含 `原文件名` 和 `存储类型`，存储类型用数字表示，0 为普通存储，1 为低频存储，2 为归档存储，3 为深度归档存储。每行多个元素名之间用分割符分隔（默认 tab 制表符）； 如果需要自定义分割符，可以使用 `-F` 或 `--sep` 选项指定自定义的分隔符。如果没有通过该选项指定该文件参数， 从标准输入读取内容。每行包含 `文件名` 和 `存储类型`；具体格式如下：（【可选】）
```
<Key><Sep>1     // <Key>：文件名，<Sep>：分割符，1：低频存储。
```
- -y/--force：该选项控制工具的默认行为。默认情况下，对于批量操作，工具会要求使用者输入一个验证码，确认下要进行批量文件操作了，避免操作失误的发生。如果不需要这个验证码的提示过程，可以使用 `--force` 选项。【可选】
- -s/--success-list：该选项指定一个文件，qshell 会把操作成功的文件行导入到该文件；默认不导出。【可选】
- -e/--failure-list：该选项指定一个文件， qshell 会把操作失败的文件行加上错误状态码，错误的原因导入该文件；默认不导出。【可选】
- -F/--sep：该选项可以自定义输入内容（文件输入或标准输入，参考 -i 选项说明）；默认为 tab 制表符。【可选】
- -c/--worker：该选项可以定义 Batch 任务并发数；默认为 1。【可选】
- --enable-record：记录任务执行状态，当下次执行命令时会跳过已执行的任务。 【可选】
- --record-redo-while-error：依赖于 --enable-record，当检测任务状态时（命令重新执行时，所有任务会从头到尾重新执行；任务执行前会先检测当前任务是否已经执行），如果任务已执行且失败，则再执行一次；默认为 false，当任务执行失败不重新执行。 【可选】

# 示例
1 比如我们要将空间 `if-pbl` 里面的一些文件改为低频存储，我们可以指定如下的`KeyFileTypeMapFile` 的内容：
```
2015/03/22/qiniu.png	1
2015/photo.jpg	1
2015/03/22/qiniu2.png	0
2015/photo2.jpg	2
```

上面，我们将 `2015/03/22/qiniu.png` 文件设置为低频存储了，诸如此类。
把这个内容保存到文件 `tochangetype.txt` 中，然后使用如下的命令将 `tochangetype.txt` 中所有的文件进行存储类型改变。

```
$ qshell batchchtype if-pbl -i tochangetype.txt
```

2 如果不希望上面的重命名过程出现验证码提示，可以使用 `-force` 选项：
```
$ qshell batchchtype --force if-pbl -i tochangetype.txt
```

# 注意
如果没有指定输入文件的话，默认会从标准输入读取同样格式的内容
