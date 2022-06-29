# 简介
`batchrename` 命令用来为空间中的文件进行重命名。该操作发生在同一个空间中。

当然，如果这个时候源文件和目标文件名相同，那么复制会失败（这个操作其实没有意义），如果你实在想这样做，可以用 `--overwrite` 选项。

# 格式
```
qshell batchrename [--force] [--success-list <SuccessFileName>] [--failure-list <FailureFileName>] [--sep <Separator>]  [--worker <WorkerCount>] <Bucket> [-i <OldNewKeyMapFile>]
```

# 鉴权
需要在使用了 `account` 设置了 `AccessKey`, `SecretKey` 和 `Name` 的情况下使用。

# 参数
- Bucket：空间名，可以为公开空间或私有空间。【必选】

# 选项
- -i/--input-file：接受一个文件参数, 文件中每行包含 `原文件名` 和 `目标文件名`；注意这里 `目标文件名` 不可以和 `原文件名` 相同，否则对于这个文件来说就是重命名失败。每行多个元素名之间用分割符分隔（默认 tab 制表符）； 如果需要自定义分割符，可以使用 `-F` 或 `--sep` 选项指定自定义的分隔符。如果没有通过该选项指定该文件参数， 从标准输入读取内容。文件每行格式如下：（【可选】）
```
<OldKey><Sep><NewKey> // <OldKey>：原文件名，<Sep>：分割符，<NewKey>：新文件名。
```
- -y/--force：该选项控制工具的默认行为。默认情况下，对于批量操作，工具会要求使用者输入一个验证码，确认下要进行批量文件操作了，避免操作失误的发生。如果不需要这个验证码的提示过程，可以使用 `--force` 选项。【可选】
- -s/--success-list：该选项指定一个文件，qshell 会把操作成功的文件行导入到该文件；默认不导出。【可选】
- -e/--failure-list：该选项指定一个文件， qshell 会把操作失败的文件行加上错误状态码，错误的原因导入该文件；默认不导出。【可选】
- -F/--sep：该选项可以自定义输入内容（文件输入或标准输入，参考 -i 选项说明）；默认为 tab 制表符。【可选】
- -c/--worker：该选项可以定义 Batch 任务并发数；1 路并发单次操作对象数为 1000 ，如果配置为 2 并发，则 2 路并发单次操作对象数为 2000，此值需要和七牛对您的操作上限相吻合，否则会出现非预期错误，正常情况不需要调节此值，如果需要请谨慎调节；默认为 1。【可选】
- --overwrite：默认情况下，如果批量重命名的文件列表中存在目标空间已有同名文件的情况，针对该文件的重命名会失败，如果希望能够强制覆盖目标文件，那么可以使用 `--overwrite` 选项。【可选】
- --enable-record：记录任务执行状态，当下次执行命令时会跳过已执行的任务。 【可选】
- --record-redo-while-error：依赖于 --enable-record，当检测任务状态时（命令重新执行时，所有任务会从头到尾重新执行；任务执行前会先检测当前任务是否已经执行），如果任务已执行且失败，则再执行一次；默认为 false，当任务执行失败不重新执行。 【可选】

# 示例
1 比如我们要将空间 `if-pbl` 里面的一些文件进行重命名，我们可以指定如下的 `OldNewKeyMapFile` 的内容：
```
2015/03/22/qiniu.png	test/qiniu.png
2015/photo.jpg	test/photo.jpg
```

上面，我们将 `2015/03/22/qiniu.png` 重命名为 `test/qiniu.png`，诸如此类。
把这个内容保存到文件`torename.txt`中，然后使用如下的命令将所有的文件进行重命名。
```
$ qshell batchrename if-pbl -i torename.txt
```

2 如果不希望上面的重命名过程出现验证码提示，可以使用 `--force` 选项：
```
$ qshell batchrename --force if-pbl -i torename.txt
```

3 对于重新命名的过程中，希望导入成功失败的文件，可以这样导出
```
$ qshell batchrename if-pbl -i torename.txt --success-list success.txt --failure-list failure.txt
```

如果都重命名成功，success.txt的内容为：
```
2015/03/22/qiniu.png	test/qiniu.png
2015/photo.jpg	test/photo.jpg
```

# 注意 
如果没有指定输入文件的话， 会从标准输入读取内容
