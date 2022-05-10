# 简介
`restorear` 命令用来恢复一个归档文件，并且在 <FreezeAfterDays> 天之后再次恢复为原来的归档状态。<FreezeAfterDays> 解冻有效期 1～7 天。

归档存储文件完成解冻通常需要 1～5分钟，深度归档存储文件完成解冻需要 5～12 小时。

注：恢复仅仅是让文件可以进行下载等操作，并不会真的修改存储类型， 如果想把归档或者深度归档存储文件的存储类型转为标准存储，那么需要先将文件进行恢复(本命令)，再修改文件存储类型（chtype 命令）。

参考文档：[解冻归档/深度归档存储文件](https://developer.qiniu.com/kodo/6380/restore-archive)

# 格式
```
qshell restorear <Bucket> <Key> <FreezeAfterDays> [flags]
```

# 参数
- Bucket: 源空间名称 【必须】
- Key: 源文件名称 【必须】
- FreezeAfterDays: 恢复的有效期，单位：天。 【必须】

# 示例
```
// 把 test 空间下的 a.txt 文件恢复，有效期为 3 天
$ qshell restorear test a.txt 3
```