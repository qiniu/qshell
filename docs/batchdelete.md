# 简介
`batchdelete` 命令用来根据一个七牛空间中的文件名列表来批量删除空间中的这些文件。

# 格式
```
qshell batchdelete [--force] [--sucess-list <SuccessFileName>] [--failure-list <FailureFileName>] [--sep <Separator>]  [--worker <WorkerCount>] <Bucket> [-i <KeyListFile>]
```

# 帮助
```
qshell batchdelete -h
```

# 鉴权
需要在使用了 `account` 设置了 `AccessKey`, `SecretKey` 和 `Name` 的情况下使用。

# 参数
|   参数名 |               描述             |
|----------|--------------------------------|
|  Bucket  |空间名，可以为公开空间或私有空间|

##### i 短选项
接受一个文件参数， 文件内容每行包含 `文件名` 和 `文件 PutTime`；如果指定文件 PutTime 则当七牛云存储的文件 PutTime 和 该 PutTime 相等才会删除。
具体格式如下：
```
// 不指定文件 PutTime
<Key> // key：文件名

// 指定文件 PutTime
<Key><Sep>16445676785097143 // key：文件名，<Sep>：分割符，16445676785097143：文件 PutTime。
```
每行多个元素名之间用分割符分隔（默认 tab 制表符）； 如果需要自定义分割符，可以使用 `-F` 或 `--sep` 选项指定自定义的分隔符。
如果没有通过该选项指定该文件参数， 从标准输入读取内容。

##### force 选项
该选项控制工具的默认行为。默认情况下，对于批量操作，工具会要求使用者输入一个验证码，确认下要进行批量文件操作了，避免操作失误的发生。如果不需要这个验证码的提示过程，可以使用 `--force` 选项。

##### success-list 选项
该选项指定一个文件，qshell 会把操作成功的文件行导入到该文件；默认不导出。

##### failure-list 选项
该选项指定一个文件， qshell 会把操作失败的文件行加上错误状态码，错误的原因导入该文件；默认不导出。

##### sep 选项
该选项可以自定义输入内容（文件输入或标准输入，参考 -i 选项说明）；默认为 tab 制表符。

##### worker
该选项可以定义 Batch 任务并发数；默认为 1。

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
