# 简介
`batchforbidden` 指令用来批量禁止或启用七牛云存储的文件；禁用后文件不可访问。

# 格式
```
qshell batchforbidden [--force] [--success-list <SuccessFileName>] [--failure-list <FailureFileName>] [--sep <Separator>]  [--worker <WorkerCount>] [-i <KeyFileTypeMapFile>] [--reverse] <Bucket>
```

# 鉴权
需要在使用了 `account` 设置了 `AccessKey` 和 `SecretKey` 的情况下使用。

# 参数
- Bucket：需要验证文件所在空间名称，可以为公开空间或者私有空间【必选】

# 选项
- i/--input-file：接受一个文件参数，文件内容每行包含待检查文件的 Key 等信息。每行多个元素名之间用分割符分隔（默认 tab 制表符）； 如果需要自定义分割符，可以使用 `-F` 或 `--sep` 选项指定自定义的分隔符；也可以直接使用 list 接口结果保存的文件。如果没有通过该选项指定该文件参数， 从标准输入读取内容。 具体格式如下：（【可选】）
```
<Key> // <Key>: 七牛云存储的 Key
```
- -s/--success-list：该选项指定一个文件，qshell 会把操作成功的文件行导入到该文件；默认不导出。【可选】
- -e/--failure-list：该选项指定一个文件， qshell 会把操作失败的文件行加上错误状态码，错误的原因导入该文件；默认不导出。【可选】
- -F/--sep：该选项可以自定义每行输入内容中字段之间的分隔符（文件输入或标准输入，参考 -i 选项说明）；默认为 tab 制表符。【可选】
- -c/--worker：该选项可以定义 Batch 任务并发数；默认为 1。【可选】
- --enable-record：记录任务执行状态，当下次执行命令时会检测任务执行的状态并跳过已执行的任务。 【可选】
- --record-redo-while-error：依赖于 --enable-record，当检测任务状态时（命令重新执行时，所有任务会从头到尾重新执行；任务执行前会先检测当前任务是否已经执行），如果任务已执行且失败，则再执行一次；默认此选项不生效，当任务执行失败不重新执行。 【可选】
- -r/--reverse: 启用指定文件时指定。【可选】

# 示例
1. 禁用 if-pbl 空间下的 hello01.json 和 hello02.json 两个文件
```
// 创建禁用信息文件 forbidden_list.txt ，文件信息如下：
hello01.json
hello02.json 

// 调用命令
qshell batchforbidden if-pbl -i ./forbidden_list.txt
```

2. 启用 if-pbl 空间下的 hello01.json 和 hello02.json 两个文件
```
// 创建启用信息文件 forbidden_list.txt ，文件信息如下：
hello01.json
hello02.json 

// 调用命令
qshell batchforbidden if-pbl -r -i ./forbidden_list.txt
```