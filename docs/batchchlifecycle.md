# 简介
`batchchlifecycle` 命令用于批量修改已上传文件 Object 的生命周期。后续可以通过 `qshell stat` 命令查看文件修改生命周期的相关时间。

# 注：
1. 生命周期值的范围是：-1 或者 大于 0，单位：天
   * 小于 -1: 没有任何意义，不会产生任何效果
   * 等于 -1: 取消已设置的相关生命周期
   * 等于  0: 没有任何意义，不会产生任何效果
   * 大于  0: 设置相关的生命周期
2. 生命周期时间大小规则如下（在相关生命周期时间值大于 0 时需满足）：
```
转低频存储时间 < 转归档存储时间 < 转深度归档存储时间 
```
3. 转低频存储时间、转归档存储时间、转深度归档存储时间 和 过期删除时间 至少配置一个

# 格式
```
qshell batchchlifecycle [--force] [--success-list <SuccessFileName>] [--failure-list <FailureFileName>] [--sep <Separator>]  [--worker <WorkerCount>] [--to-ia-after-days <ToIAAfterDays>] [--to-archive-after-days <ToArchiveAfterDays>] [--to-deep-archive-after-days <ToDeepArchiveAfterDays>] [--delete-after-days <DeleteAfterDays>] <Bucket> <-i KeysFile> 
```

# 帮助文档
可以在命令行输入如下命令获取帮助文档：
```
// 简单描述
$ qshell batchchlifecycle -h 

// 详细文档（此文档）
$ qshell batchchlifecycle --doc
```

# 鉴权
需要使用 `qshell account` 或者 `qshell user add` 命令设置鉴权信息 `AccessKey`, `SecretKey` 和 `Name`。

# 参数
- Bucket：空间名，可以为公开空间或私有空间。【必须】

# 选项
- -i/--input-file：指定一个文件， 文件内容每行包含 `文件名`。每行多个元素名之间用分割符分隔（默认 tab 制表符）； 如果需要自定义分割符，可以使用 `-F` 或 `--sep` 选项指定自定义的分隔符。如果没有通过该选项指定该文件参数， 从标准输入读取内容。每行具体格式如下：（【可选】）
```
// 情景一
<Key>

// 情景二
<Key><Sep><Other> // <Key>：文件名，<Sep>：分割符，<Other> 其他无效内容
```
- --to-ia-after-days：指定文件上传后并在设置的时间后转换到 `低频存储类型`；值范围为 -1 或者大于 0，设置为 -1 表示取消已设置的转 `低频存储` 的生命周期规则，单位：天【可选】
- --to-archive-after-days：指定文件上传后并在设置的时间后转换到 `归档存储类型`；值范围为 -1 或者大于 0，设置为 -1 表示取消已设置的转 `归档存储` 的生命周期规则，单位：天【可选】
- --to-deep-archive-after-days：指定文件上传后并在设置的时间后转换到 `深度归档存储类型`；值范围为 -1 或者大于 0，设置为 -1 表示取消已设置的转 `深度归档存储` 的生命周期规则，单位：天【可选】
  - --delete-after-days：指定文件上传后并在设置的时间后进行 `过期删除`，删除后不可恢复；值范围为 -1 或者大于 0，设置为 -1 表示取消已设置的 `过期删除` 的生命周期规则，单位：天【可选】
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
1 比如我们要将空间 `if-pbl` 里面一些文件的生命周期改为 30 天后转低频存储，60 天后转归档存储，180 天后转深度归档存储，365 天后过期删除；我们可以指定如下的 `KeysFile` 的内容：
```
2015/03/22/qiniu.png
2015/photo.jpg
```

把这个内容保存到文件 `lifecycle.txt` 中，然后使用如下的命令对 `lifecycle.txt` 中的所有文件设置生命周期。
```
$ qshell batchchlifecycle if-pbl -i lifecycle.txt \
 --to-ia-after-days 30 \
 --to-archive-after-days 60 \
 --to-deep-archive-after-days 180 \
 --delete-after-days 365
```

2 如果不希望上面的重命名过程出现验证码提示，可以使用 `--force` 选项：
```
$ qshell batchchlifecycle force if-pbl -i lifecycle.txt \
 --to-ia-after-days 30 \
 --to-archive-after-days 60 \
 --to-deep-archive-after-days 180 \
 --delete-after-days 365
```

# 注意
如果没有指定输入文件的话，会从标准输入读取内容。
