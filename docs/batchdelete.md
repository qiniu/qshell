# 简介
`batchdelete` 命令用来根据一个七牛空间中的文件名列表来批量删除空间中的这些文件。

# 格式
```
qshell batchdelete [--force] [--success-list <SuccessFileName>] [--failure-list <FailureFileName>] [--sep <Separator>]  [--worker <WorkerCount>] <Bucket> [-i <KeyListFile>]
```

# 帮助
```
qshell batchdelete -h
```

# 鉴权
需要在使用了 `account` 设置了 `AccessKey`, `SecretKey` 和 `Name` 的情况下使用。

# 参数
- Bucket：空间名，可以为公开空间或私有空间。【必须】

# 选项
- -i/--input-file：接受一个文件参数， 文件内容每行包含 `文件名` 和 `文件 PutTime`；如果指定文件 PutTime 则当七牛云存储的文件 PutTime 和 该 PutTime 相等才会删除。每行多个元素名之间用分割符分隔（默认 tab 制表符）； 如果需要自定义分割符，可以使用 `-F` 或 `--sep` 选项指定自定义的分隔符。如果没有通过该选项指定该文件参数， 从标准输入读取内容。每行具体格式如下：（【可选】）
```
// 不指定文件 PutTime
<Key> // key：文件名

// 指定文件 PutTime
Key<Sep>PutTime // 第一行标记每行字段信息，其他行如下
<Key><Sep>16445676785097143 // key：文件名，<Sep>：分割符，16445676785097143：文件 PutTime。
```
- -y/--force：该选项控制工具的默认行为。默认情况下，对于批量操作，工具会要求使用者输入一个验证码，确认下要进行批量文件操作了，避免操作失误的发生。如果不需要这个验证码的提示过程，可以使用 `--force` 选项。【可选】
- -s/--success-list：该选项指定一个文件，qshell 会把操作成功的文件行导入到该文件；默认不导出。【可选】
- -e/--failure-list：该选项指定一个文件， qshell 会把操作失败的文件行加上错误状态码，错误的原因导入该文件；默认不导出。【可选】
- -F/--sep：该选项可以自定义输入内容（文件输入或标准输入，参考 -i 选项说明）；默认为 tab 制表符。【可选】
- -c/--worker：该选项可以定义 Batch 任务并发数；默认为 1。【可选】
- --enable-record：记录任务执行状态，当下次执行命令时会跳过已执行的任务。 【可选】
- --record-redo-while-error：依赖于 --enable-record，当检测任务状态时（命令重新执行时，所有任务会从头到尾重新执行；任务执行前会先检测当前任务是否已经执行），如果任务已执行且失败，则再执行一次；默认为 false，当任务执行失败不重新执行。 【可选】

# 示例
1 指定要删除的文件列表 `todelete.txt` 进行删除：
```
a.jpg
test/b.jpg
```

```
$ qshell batchdelete if-pbl -i todelete.txt
```

2 删除空间 `if-pbl` 中的所有文件：
```
$ qshell listbucket if-pbl -o if-pbl.list.txt
$ qshell batchdelete --force if-pbl -i if-pbl.list.txt
```

3 如果希望导出成功和失败的文件列表
```
$ qshell batchdelete if-pbl -i if-pbl.list.txt --success-list success.txt --failure-list failure.txt
```

4 对于要删除的文件名字包含了空格的情况， 那么可以指定自定义的分隔符对文件每行进行分割, 假如使用\t进行分割
```
$ qshell batchdelete -F '\t' if-pbl -i todelete.txt
```
