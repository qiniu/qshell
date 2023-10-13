# 简介
`batchstat` 命令用来批量查询七牛空间中文件的基本信息。

# 格式
```
qshell batchstat [--force] [--success-list <SuccessFileName>] [--failure-list <FailureFileName>] [--sep <Separator>]  [--worker <WorkerCount>] <Bucket> <-i KeyListFile>
```

# 帮助文档
可以在命令行输入如下命令获取帮助文档：
```
// 简单描述
$ qshell batchstat -h 

// 详细文档（此文档）
$ qshell batchstat --doc
```

# 鉴权
需要使用 `qshell account` 或者 `qshell user add` 命令设置鉴权信息 `AccessKey`, `SecretKey` 和 `Name`。

# 参数
- Bucket：空间名，可以为公开空间或私有空间。【必选】

# 选项
- -i/--input-file：指定一个文件，文件为要 stat 的文件列表, 每行包含一个文件 Key。如果没有通过该选项指定该文件参数， 从标准输入读取内容。每行具体格式如下：（【可选】）
```
<Key> // <Key>：文件名
```
- -y/--force：该选项控制工具的默认行为。默认情况下，对于批量操作，工具会要求使用者输入一个验证码，确认下要进行批量文件操作了，避免操作失误的发生。如果不需要这个验证码的提示过程可以使用此选项。【可选】
- -s/--success-list：该选项指定一个文件，程序会把操作成功的资源信息导入到该文件；默认不导出。【可选】
- -e/--failure-list：该选项指定一个文件，程序会把操作失败的资源信息加上错误信息导入该文件；默认不导出。【可选】
- -o/--outfile：该选项指定一个文件，把 stat 结果导入到此文件中。注：输出的内容顺序和 input file 内容的顺序会有不同【可选】
- -c/--worker：该选项可以定义 Batch 任务并发数；1 路并发单次操作对象数为 250 ，如果配置为 10 并发，则 10 路并发单次操作对象数为 2500，此值需要和七牛对您的操作上限相吻合，否则会出现非预期错误，正常情况不需要调节此值，如果需要请谨慎调节；默认为 4。【可选】
- --min-worker：最小 Batch 任务并发数；当并发设置过高时，会触发超限错误，为了缓解此问题，qshell 会自动减小并发度，此值为减小的最低值。默认：1【可选】
- --worker-count-increase-period：为了尽可能快的完成操作 qshell 会周期性尝试增加并发度，此值为尝试增加并发数的周期，单位：秒，最小 10，默认 60。【可选】
- --enable-record：记录任务执行状态，当下次执行命令时会检测任务执行的状态并跳过已执行的任务。 【可选】
- --record-redo-while-error：依赖于 --enable-record；命令重新执行时，命令中所有任务会从头到尾重新执行；每个任务执行前会根据记录先查看当前任务是否已经执行，如果任务已执行且失败，则再执行一次；默认为 false，当任务执行失败则跳过不再重新执行。 【可选】

# 示例
- 我们将查询空间 `7qiniu` 中的一些文件的基本信息，待查询文件列表 `listFile` 的内容为：
```
RclviFDHaQAUl3aL46jKRskUWbg=/FpwH76F3yfYmFKoPDjoSNWzeLKYp/000000.ts
RclviFDHaQAUl3aL46jKRskUWbg=/FpwH76F3yfYmFKoPDjoSNWzeLKYp/000001.ts
RclviFDHaQAUl3aL46jKRskUWbg=/FpwH76F3yfYmFKoPDjoSNWzeLKYp/000002.ts
RclviFDHaQAUl3aL46jKRskUWbg=/FpwH76F3yfYmFKoPDjoSNWzeLKYp/000003.ts
RclviFDHaQAUl3aL46jKRskUWbg=/FpwH76F3yfYmFKoPDjoSNWzeLKYp/000004.ts
RclviFDHaQAUl3aL46jKRskUWbg=/FpwH76F3yfYmFKoPDjoSNWzeLKYp/000005.ts

```

- 使用如下命令进行批量查询
```
$ qshell batchstat 7qiniu -i listFile
```

- 输出 Key、Fsize、Hash、MimeType、PutTime 以 `\t` 分隔：
```
RclviFDHaQAUl3aL46jKRskUWbg=/FpwH76F3yfYmFKoPDjoSNWzeLKYp/000000.ts 92308   Fk8Uf2SHbQ4S2-cXHINuRc_rooNA    video/mp2t  15003760414606314
RclviFDHaQAUl3aL46jKRskUWbg=/FpwH76F3yfYmFKoPDjoSNWzeLKYp/000001.ts 91556   FpJP2nfipuLVc6QGvvcb868Rd0pO    video/mp2t  15003760414789673
RclviFDHaQAUl3aL46jKRskUWbg=/FpwH76F3yfYmFKoPDjoSNWzeLKYp/000002.ts 92496   FvBjZPch6cf52t2x0ZQBngqS1KTp    video/mp2t  15003760417159000
RclviFDHaQAUl3aL46jKRskUWbg=/FpwH76F3yfYmFKoPDjoSNWzeLKYp/000003.ts 92308   FoEgsbzdrcLuj_Fo5FeTI3w1jFHJ    video/mp2t  15003760419154144
RclviFDHaQAUl3aL46jKRskUWbg=/FpwH76F3yfYmFKoPDjoSNWzeLKYp/000004.ts 92308   FkYNctlf1JOGcJa-WzWgxsqcBjX6    video/mp2t  15003760422258065
RclviFDHaQAUl3aL46jKRskUWbg=/FpwH76F3yfYmFKoPDjoSNWzeLKYp/000005.ts 92120   Fh4Fwhu3dMUGbd3jE5OmRtfVZLv4    video/mp2t  15003760423842522
```

# 注意
如果没有指定输入文件， 默认从标准输入读取内容。
