# 简介
`chlifecycle` 命令用于修改已上传文件 Object 的生命周期。后续可以通过 `qshell stat` 命令查看文件修改生命周期的相关时间。

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
qshell chlifecycle [--to-ia-after-days <ToIAAfterDays>] [--to-archive-after-days <ToArchiveAfterDays>] [--to-deep-archive-after-days <ToDeepArchiveAfterDays>] [--delete-after-days <DeleteAfterDays>] <Bucket> <Key> 
```

# 帮助文档
可以在命令行输入如下命令获取帮助文档：
```
// 简单描述
$ qshell chlifecycle -h 

// 详细文档（此文档）
$ qshell chlifecycle --doc
```

# 鉴权
需要使用 `qshell account` 或者 `qshell user add` 命令设置鉴权信息 `AccessKey`, `SecretKey` 和 `Name`。

# 参数
- Bucket：空间名，可以为公开空间或私有空间。【必须】

# 选项
- --to-ia-after-days：指定文件上传后并在设置的时间后转换到 `低频存储类型`；值范围为 -1 或者大于 0，设置为 -1 表示取消已设置的转 `低频存储` 的生命周期规则，单位：天【可选】
- --to-archive-after-days：指定文件上传后并在设置的时间后转换到 `归档存储类型`；值范围为 -1 或者大于 0，设置为 -1 表示取消已设置的转 `归档存储` 的生命周期规则，单位：天【可选】
- --to-deep-archive-after-days：指定文件上传后并在设置的时间后转换到 `深度归档存储类型`；值范围为 -1 或者大于 0，设置为 -1 表示取消已设置的转 `深度归档存储` 的生命周期规则，单位：天【可选】
- --delete-after-days：指定文件上传后并在设置的时间后进行 `过期删除`，删除后不可恢复；值范围为 -1 或者大于 0，设置为 -1 表示取消已设置的 `过期删除` 的生命周期规则，单位：天【可选】


# 示例
1 比如我们要将空间 `if-pbl` 里面 `qiniu.png` 文件的生命周期改为 30 天后转低频存储，60 天后转归档存储，180 天后转深度归档存储，365 天后过期删除：
```
$ qshell chlifecycle if-pbl qiniu.png \
 --to-ia-after-days 30 \
 --to-archive-after-days 60 \
 --to-deep-archive-after-days 180 \
 --delete-after-days 365
```

2 查询修改效果：
```
$ qshell stat if-pbl qiniu.png

// 命令输出：
Bucket:                  if-pbl
Key:                     qiniu.png
FileHash:                lozgLP_MAdAKZkPCXGvfd0LIDSUI
Fsize:                   5444314 -> 5.19MB
PutTime:                 16768889367943931 -> 2023-02-20 18:28:56.7943931 +0800 CST
MimeType:                text/plain
Expiration:              1710518400 -> 2024-03-16 00:00:00 +0800 CST
TransitionToIA:          1681574400 -> 2023-04-16 00:00:00 +0800 CST
TransitionToArchive:     1684166400 -> 2023-05-16 00:00:00 +0800 CST
TransitionToDeepArchive: 1694534400 -> 2023-09-13 00:00:00 +0800 CST
FileType:                0 -> 标准存储
```
