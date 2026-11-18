# 简介
`batchstat` 命令用来批量查询七牛空间中文件的基本信息。

# 格式
```
qshell batchstat [--force] [--show-fields <Fields>] [--success-list <SuccessFileName>] [--failure-list <FailureFileName>] [--outfile <ResultFileName>] [--worker <WorkerCount>] <Bucket> [-i KeyListFile]
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
- --show-fields：指定需要输出的字段，多个字段使用逗号分隔，字段名不区分大小写，输出顺序与指定顺序一致；未知字段或空字段名会报错。不指定或值为空时，保持原有六列输出。【可选】
- -c/--worker：该选项可以定义 Batch 任务并发数；1 路并发单次操作对象数为 250 ，如果配置为 10 并发，则 10 路并发单次操作对象数为 2500，此值需要和七牛对您的操作上限相吻合，否则会出现非预期错误，正常情况不需要调节此值，如果需要请谨慎调节；默认为 4。【可选】
- --min-worker：最小 Batch 任务并发数；当并发设置过高时，会触发超限错误，为了缓解此问题，qshell 会自动减小并发度，此值为减小的最低值。默认：1【可选】
- --worker-count-increase-period：为了尽可能快的完成操作 qshell 会周期性尝试增加并发度，此值为尝试增加并发数的周期，单位：秒，最小 10，默认 60。【可选】
- --enable-record：记录任务执行状态，当下次执行命令时会检测任务执行的状态并跳过已执行的任务。 【可选】
- --record-redo-while-error：依赖于 --enable-record；命令重新执行时，命令中所有任务会从头到尾重新执行；每个任务执行前会根据记录先查看当前任务是否已经执行，如果任务已执行且失败，则再执行一次；默认为 false，当任务执行失败则跳过不再重新执行。 【可选】

# 输出字段

默认仍按 `Key、Fsize、Hash、MimeType、PutTime、Type` 的顺序输出，无表头。指定 `--show-fields` 后，只输出所选字段，不自动追加其他字段。终端结果行与 `--outfile` 使用相同字段配置，以 Tab（`\t`）分隔。

| 字段 | 含义 | 格式及缺失值 |
|------|------|-------------|
| Key | 对象名称 | 字符串 |
| Fsize | 文件大小 | 十进制字节数 |
| Hash | 文件 Etag | 字符串 |
| MimeType | MIME 类型 | 字符串 |
| PutTime | 上传时间 | 保持原有的 100 纳秒单位 |
| Type | 存储类型 | 0 标准、1 低频、2 归档、3 深度归档、4 归档直读、5 智能分层 |
| EndUser | 上传时设置的用户标识 | 字符串，缺失时为空列 |
| MD5 | 文件 MD5 | 字符串，缺失时为空列 |
| Status | 文件禁用状态 | 数字或空列；1 表示禁用。SDK 将启用状态的 0 转为空指针，因此通常显示为空列 |
| RestoreStatus | 解冻状态 | 数字或空列；1 解冻中、2 解冻完成；接口未返回或为 0 时，SDK 返回空指针，输出空列 |
| Expiration | 过期删除时间 | Unix 秒，SDK 空指针时输出空列 |
| TransitionToIA | 转低频时间 | Unix 秒，SDK 空指针时输出空列 |
| TransitionToIntelligentTiering | 转智能分层时间 | Unix 秒，SDK 空指针时输出空列 |
| TransitionToArchiveIR | 转归档直读时间 | Unix 秒，当前 SDK 使用值类型，缺失时输出 `0` |
| TransitionToArchive | 转归档时间 | Unix 秒，SDK 空指针时输出空列 |
| TransitionToDeepArchive | 转深度归档时间 | Unix 秒，SDK 空指针时输出空列 |

空指针和空字符串统一输出空列，保留对应的 Tab 分隔符（包括行尾 Tab）；有值的数字 `0` 正常输出。解析时应按 Tab 拆分并保留空字段，避免按任意空白切分或先去除行尾空白，否则连续空列和末尾空列会丢失。

生命周期时间为空或为 `0`，仅表示接口没有提供有效时间，不能据此判断对象永不过期；部分历史对象会执行生命周期操作，但接口不返回对应时间。

`--success-list` 保持保存成功的输入记录，`--failure-list` 保持原有失败记录格式，扩展属性写入 `--outfile`。失败对象不产生属性结果行。并发输出不保证与输入顺序相同，应使用 Key 关联对象。

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
$ qshell batchstat 7qiniu --silence -y -i listFile
```

- 输出 Key、Fsize、Hash、MimeType、PutTime、Type 以 `\t` 分隔：
```
RclviFDHaQAUl3aL46jKRskUWbg=/FpwH76F3yfYmFKoPDjoSNWzeLKYp/000000.ts 92308   Fk8Uf2SHbQ4S2-cXHINuRc_rooNA    video/mp2t  15003760414606314   0
RclviFDHaQAUl3aL46jKRskUWbg=/FpwH76F3yfYmFKoPDjoSNWzeLKYp/000001.ts 91556   FpJP2nfipuLVc6QGvvcb868Rd0pO    video/mp2t  15003760414789673   0
RclviFDHaQAUl3aL46jKRskUWbg=/FpwH76F3yfYmFKoPDjoSNWzeLKYp/000002.ts 92496   FvBjZPch6cf52t2x0ZQBngqS1KTp    video/mp2t  15003760417159000   0
RclviFDHaQAUl3aL46jKRskUWbg=/FpwH76F3yfYmFKoPDjoSNWzeLKYp/000003.ts 92308   FoEgsbzdrcLuj_Fo5FeTI3w1jFHJ    video/mp2t  15003760419154144   0
RclviFDHaQAUl3aL46jKRskUWbg=/FpwH76F3yfYmFKoPDjoSNWzeLKYp/000004.ts 92308   FkYNctlf1JOGcJa-WzWgxsqcBjX6    video/mp2t  15003760422258065   0
RclviFDHaQAUl3aL46jKRskUWbg=/FpwH76F3yfYmFKoPDjoSNWzeLKYp/000005.ts 92120   Fh4Fwhu3dMUGbd3jE5OmRtfVZLv4    video/mp2t  15003760423842522   0
--------------- Batch Result ---------------
              Total:         6
            Success:         6
            Failure:         0
            Skipped:         0
           Duration:        1s
--------------------------------------------
```

# 注意
如果没有指定输入文件， 默认从标准输入读取内容。

## 自选字段并筛选过期时间

```bash
qshell batchstat mybucket -i keys.txt -y \
  --show-fields Key,Expiration,Status,MD5 --outfile result.tsv
```

结果为四列；若 MD5 为空，第四列仍保留。建议使用 `--outfile` 获取供程序解析的结果，终端输出还包含进度和汇总信息。

以下 Python 示例保留空字段，仅选择有明确过期时间且早于阈值的对象：

```python
import csv

with open("result.tsv", encoding="utf-8", newline="") as result:
    for key, expiration, status, md5 in csv.reader(result, delimiter="\t"):
        if expiration and 0 < int(expiration) < 1893456000:
            print(key)
```
