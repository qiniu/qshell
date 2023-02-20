# 简介
`batchdelete` 命令用来根据一个七牛空间中的文件名列表来批量删除空间中的这些文件。

# 格式
```
qshell batchdelete [--force] [--success-list <SuccessFileName>] [--failure-list <FailureFileName>] [--sep <Separator>]  [--worker <WorkerCount>] <Bucket> [-i <KeyListFile>]
```

# 帮助文档
可以在命令行输入如下命令获取帮助文档：
```
// 简单描述
$ qshell batchdelete -h 

// 详细文档（此文档）
$ qshell batchdelete --doc
```

# 鉴权
需要使用 `qshell account` 或者 `qshell user add` 命令设置鉴权信息 `AccessKey`, `SecretKey` 和 `Name`。

# 参数
- Bucket：空间名，可以为公开空间或私有空间。【必须】

# 选项
- -i/--input-file：指定一个文件， 文件内容每行包含 `文件名` 和 `文件 PutTime`；如果指定文件 PutTime 则当七牛云存储的文件 PutTime 和 该 PutTime 相等才会删除。每行多个元素名之间用分割符分隔（默认 tab 制表符）； 如果需要自定义分割符，可以使用 `-F` 或 `--sep` 选项指定自定义的分隔符。如果没有通过该选项指定该文件参数， 从标准输入读取内容。每行具体格式如下：（【可选】）
```
// 不指定文件 PutTime
<Key> // key：文件名

// 指定文件 PutTime
Key<Sep>PutTime // 第一行为标题，用于表述每行信息的样式，其他行样式如下
<Key><Sep><PutTime> // key：文件名，<Sep>：分割符；<PutTime>：文件上传时间，单位：100*ns，eg:16445676785097143。
```
- -y/--force：该选项控制工具的默认行为。默认情况下，对于批量操作，工具会要求使用者输入一个验证码，确认下要进行批量文件操作了，避免操作失误的发生。如果不需要这个验证码的提示过程可以使用此选项。【可选】
- -s/--success-list：该选项指定一个文件，程序会把操作成功的资源信息导入到该文件；默认不导出。【可选】
- -e/--failure-list：该选项指定一个文件，程序会把操作失败的资源信息加上错误信息导入该文件；默认不导出。【可选】
- -F/--sep：该选项可以自定义每行输入内容中字段之间的分隔符（文件输入或标准输入，参考 -i 选项说明）；默认为 tab 制表符。【可选】
- -c/--worker：该选项可以定义 Batch 任务并发数；1 路并发单次操作对象数为 250 ，如果配置为 10 并发，则 10 路并发单次操作对象数为 2500，此值需要和七牛对您的操作上限相吻合，否则会出现非预期错误，正常情况不需要调节此值，如果需要请谨慎调节；默认为 4。【可选】
- --min-worker：最小 Batch 任务并发数；当并发设置过高时，会触发超限错误，为了缓解此问题，qshell 会自动减小并发度，此值为减小的最低值。默认：1【可选】
- --worker-count-increase-period：为了尽可能快的完成操作 qshell 会周期性尝试增加并发度，此值为尝试增加并发数的周期，单位：秒，最小 10，默认 60。【可选】
- --enable-record：记录任务执行状态，当下次执行命令时会检测任务执行的状态并跳过已执行的任务。 【可选】
- --record-redo-while-error：依赖于 --enable-record；命令重新执行时，命令中所有任务会从头到尾重新执行；每个任务执行前会根据记录先查看当前任务是否已经执行，如果任务已执行且失败，则再执行一次；默认为 false，当任务执行失败则跳过不再重新执行。 【可选】

# 示例
1 删除空间 `if-pbl` 下的某些文件，指定要删除的文件列表 `todelete.txt` 进行删除，其内容如下：
```
a.jpg
test/b.jpg
```
执行如下命令：
```
$ qshell batchdelete if-pbl -i todelete.txt
```

2 删除空间 `if-pbl` 中的所有文件：
```
// 先列举空间：if-pbl，列举结果保存在文件 if-pbl.list.txt 中
$ qshell listbucket if-pbl -o if-pbl.list.txt
// 再根据 if-pbl.list.txt 进行删除
$ qshell batchdelete --force if-pbl -i if-pbl.list.txt
```

3 如果希望导出成功和失败的文件列表
```
$ qshell batchdelete if-pbl -i if-pbl.list.txt --success-list success.txt --failure-list failure.txt
```

4 对于要删除的文件名字包含了空格的情况， 那么可以指定自定义的分隔符对文件每行进行分割, 假如使用 \t 进行分割
```
$ qshell batchdelete -F '\t' if-pbl -i todelete.txt
```
