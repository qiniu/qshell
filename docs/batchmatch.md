# 简介
`batchmatch` 指令用来批量验证指定的七牛云存储文件内容是否和本地文件内容相匹配。

使用此命令，您需要指定待验证的七牛云文件 Key，指定方式见选项 `-i/--input-file`，命令会根据指定的 Key 在 <LocalFileDir> 路径下查找与此 Key 对应的文件，然后检查两个文件内容的一致性。

注：
本地文件和七牛云存储文件对应关系必须满足一下条件：
`${LocalFilePath} = ${LocalFileDir} + ${文件分隔符} + ${七牛存储 Key}`

# 格式
```
qshell batchmatch <Bucket> <LocalFileDir> 
```

# 帮助文档
可以在命令行输入如下命令获取帮助文档：
```
// 简单描述
$ qshell batchmatch -h 

// 详细文档（此文档）
$ qshell batchmatch --doc
```

# 鉴权
需要使用 `qshell account` 或者 `qshell user add` 命令设置鉴权信息 `AccessKey`, `SecretKey` 和 `Name`。

# 参数
- Bucket：需要验证文件所在空间名称，可以为公开空间或者私有空间【必选】
- LocalFileDir：本地存储文件的路径 【必选】

# 选项
- i/--input-file：指定一个文件，文件内容每行包含待检查文件的 Key 等信息。每行多个元素名之间用分割符分隔（默认 tab 制表符）； 如果需要自定义分割符，可以使用 `-F` 或 `--sep` 选项指定自定义的分隔符；也可以直接使用 list 接口结果保存的文件。如果没有通过该选项指定该文件参数， 从标准输入读取内容。 具体格式如下：（【可选】）
```
<Key> // <Key>: 七牛云存储的 Key
```
- -s/--success-list：该选项指定一个文件，程序会把操作成功的资源信息导入到该文件；默认不导出。【可选】
- -e/--failure-list：该选项指定一个文件，程序会把操作失败的资源信息加上错误信息导入该文件；默认不导出。【可选】
- -F/--sep：该选项可以自定义每行输入内容中字段之间的分隔符（文件输入或标准输入，参考 -i 选项说明）；默认为 tab 制表符。【可选】
- -c/--worker：该选项可以定义 Batch 任务并发数；默认为 1。【可选】
- --enable-record：记录任务执行状态，当下次执行命令时会检测任务执行的状态并跳过已执行的任务。 【可选】
- --record-redo-while-error：依赖于 --enable-record；命令重新执行时，命令中所有任务会从头到尾重新执行；每个任务执行前会根据记录先查看当前任务是否已经执行，如果任务已执行且失败，则再执行一次；默认为 false，当任务执行失败则跳过不再重新执行。 【可选】

# 示例
```
检查 if-pbl 存储空间的文件和文件夹 /Users/lala/Desktop/qiniu 下的文件是否匹配。

1. 列举空间 if-pbl 下所有文件并保存到文件： /Users/lala/Desktop/Match.conf
$ qshell listbucket2 if-pbl -o /Users/lala/Desktop/Match.conf

2. 验证文件夹 /Users/lala/Desktop/qiniu 下文件是否和 /Users/lala/Desktop/Match.conf 中的信息是否匹配。
$ qshell batchmatch if-pbl /Users/lala/Desktop/qiniu -i /Users/lala/Desktop/Match.conf
```