# 简介
`batchfetch` 命令用来批量抓取远程地址到七牛存储空间。

# 格式
```
qshell batchfetch [--force] [--success-list <SuccessFileName>] [--failure-list <FailureFileName>] [--sep <Separator>]  [--worker <WorkerCount>] <Bucket>
```

# 帮助文档
可以在命令行输入如下命令获取帮助文档：
```
// 简单描述
$ qshell batchfetch -h 

// 详细文档（此文档）
$ qshell batchfetch --doc
```

# 鉴权
需要使用 `qshell account` 或者 `qshell user add` 命令设置鉴权信息 `AccessKey`, `SecretKey` 和 `Name`。

# 参数
- Bucket：空间名，可以为公开空间或私有空间。 【必选】

# 选项
- i/--input-file：指定一个文件，文件内容每行包含待 fetch 文件的 Url 和保存的 Key, Key 可省略。每行多个元素名之间用分割符分隔（默认 tab 制表符）； 如果需要自定义分割符，可以使用 `-F` 或 `--sep` 选项指定自定义的分隔符。如果没有通过该选项指定该文件参数， 从标准输入读取内容。 具体格式如下：（【可选】）
```
// 不指定指定存储文件名
<Url>            // <Url>: 文件 url，eg:http://img.abc.com/0/000/484/0000484193.fid 保存的文件名为：0/000/484/0000484193.fid

// 指定指定存储文件名
<Url><Sep><Key> // <Url>: 文件 url，<Sep>：分割符，<Key>：文件名
```
- -y/--force：该选项控制工具的默认行为。默认情况下，对于批量操作，工具会要求使用者输入一个验证码，确认下要进行批量文件操作了，避免操作失误的发生。如果不需要这个验证码的提示过程可以使用此选项。【可选】
- -s/--success-list：该选项指定一个文件，程序会把操作成功的资源信息导入到该文件；默认不导出。【可选】
- -e/--failure-list：该选项指定一个文件，程序会把操作失败的资源信息加上错误信息导入该文件；默认不导出。【可选】
- -F/--sep：该选项可以自定义每行输入内容中字段之间的分隔符（文件输入或标准输入，参考 -i 选项说明）；默认为 tab 制表符。【可选】
- -c/--worker：该选项可以定义 Batch 任务并发数；默认为 1。【可选】
- --enable-record：记录任务执行状态，当下次执行命令时会检测任务执行的状态并跳过已执行的任务。 【可选】
- --record-redo-while-error：依赖于 --enable-record；命令重新执行时，命令中所有任务会从头到尾重新执行；每个任务执行前会根据记录先查看当前任务是否已经执行，如果任务已执行且失败，则再执行一次；默认为 false，当任务执行失败则跳过不再重新执行。 【可选】

# 使用示例
假如我们的 `AccessKey="test-ak"`, `SecretKey="test-sk"`, 我给自己账号起了个名字 `Name="myself"`

第一步:
检查 qshell 本地数据库有没有该账号，如果有该账号，会打印出来该账号的信息
```
$ qshell user lookup myself
```

如果有该账号，可以使用
```
$ qshell user cu myself
```

切换到该账号, 如果您配置了自动补全（配置方法参考 README.md)， 在命令行输入
```
$ qshell user cu <TAB>
```
会自动补全本地数据库的账户名字。

如果没有该账号，需要使用 `qshell account` 添加账号到 `qshell` 的本地数据库, 其中 `<Your AccountName>` 可以自定义, 改名字的作用只是用来在本地数据库中唯一表示账户名称。
```
$ qshell account <Your AccessKey> <Your SecretKey> <Your AccountName>
```

第二步:
使用 `batchfetch` 命令操作, 假如我要操作的 `bucket="test-bucket"`, 要预取的文件地址列表保存在文件 `batchfetchurls.txt`：
```
$ qshell batchfetch test-bucket -i batchfetchurls.txt
```

如果想导出 `fetch` 成功，失败的列表分别到文件 `fetch_success.txt`, `fetch_failure.txt`，可以使用如下命令:
```
$ qshell batchfetch test-bucket -i batchfetchurls.txt --success-list fetch_success.txt --failure-list fetch_failure.txt
```
