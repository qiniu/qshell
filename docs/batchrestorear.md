# 简介
`batchrestorear` 命令用来恢复一批归档文件，并且在 <FreezeAfterDays> 天之后再次恢复为原来的归档状态。<FreezeAfterDays> 解冻有效期 1～7 天。

归档存储文件完成解冻通常需要 1～5分钟，深度归档存储文件完成解冻需要 5～12 小时。

注：恢复仅仅是让文件可以进行下载等操作，并不会真的修改存储类型， 如果想把归档或者深度归档存储文件的存储类型转为标准存储，那么需要先将文件进行恢复(本命令)，再修改文件存储类型（chtype 命令）。

参考文档：[解冻归档/深度归档存储文件](https://developer.qiniu.com/kodo/6380/restore-archive)

# 格式
```
qshell batchrestorear <Bucket> <FreezeAfterDays> [flags]
```

# 参数
- Bucket: 源空间名称 【必须】
- FreezeAfterDays: 恢复的有效期，单位：天。 【必须】

# 选项
- -i/--input-file：接受一个文件（KeyFile）, 文件内容每行包含一个 `文件名`。每行多个元素名之间用分割符分隔（默认 tab 制表符）； 如果需要自定义分割符，可以使用 `-F` 或 `--sep` 选项指定自定义的分隔符。如果没有通过该选项指定该文件参数， 从标准输入读取内容。每行包含 `文件名`；具体格式如下：（【可选】）
```
<Key>    // <Key>：文件名
<Key>Sep><DestKey> // Key：文件名，<Sep>：分割符，DestKey：目标文件名。
```
- -y/--force：该选项控制工具的默认行为。默认情况下，对于批量操作，工具会要求使用者输入一个验证码，确认下要进行批量文件操作了，避免操作失误的发生。如果不需要这个验证码的提示过程，可以使用 `--force` 选项。【可选】
- -s/--success-list：该选项指定一个文件，qshell 会把操作成功的文件行导入到该文件；默认不导出。【可选】
- -e/--failure-list：该选项指定一个文件， qshell 会把操作失败的文件行加上错误状态码，错误的原因导入该文件；默认不导出。【可选】
- -F/--sep：该选项可以自定义输入内容（文件输入或标准输入，参考 -i 选项说明）；默认为 tab 制表符。【可选】
- -c/--worker：该选项可以定义 Batch 任务并发数；默认为 1。【可选】
- --enable-record：记录任务执行状态，当下次执行命令时会跳过已执行的任务。 【可选】
- --record-redo-while-error：依赖于 --enable-record，当检测任务状态时（命令重新执行时，所有任务会从头到尾重新执行；任务执行前会先检测当前任务是否已经执行），如果任务已执行且失败，则再执行一次；默认为 false，当任务执行失败不重新执行。 【可选】

# 示例
1 比如我们要将空间 `if-pbl` 里面的一些文件进行恢复，我们可以指定如下的 `KeyFile` 的内容：
```
2015/03/22/qiniu.png
2015/photo.jpg
2015/03/22/qiniu2.png
2015/photo2.jpg
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