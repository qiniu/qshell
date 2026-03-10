# 简介
`sandbox resume`（别名 `rs`）恢复一个或多个已暂停的沙箱。

# 格式
```
qshell sandbox resume [sandboxIDs...] [-a] [-m <metadata>]
qshell sbx rs [sandboxIDs...] [-a] [-m <metadata>]
```

# 帮助文档
```
$ qshell sandbox resume -h
$ qshell sandbox resume --doc
```

# 鉴权
需要配置 `QINIU_API_KEY` 或 `E2B_API_KEY` 环境变量，或在当前目录 `.env` 文件中设置。

# 参数
- `sandboxIDs`：要恢复的沙箱 ID 列表
- `-a, --all`：恢复所有已暂停的沙箱
- `-m, --metadata`：配合 -a/--all 使用，按元数据过滤（格式：key1=value1,key2=value2）

# 示例
1. 恢复指定沙箱
```
$ qshell sandbox resume sb-xxxxxxxxxxxx
$ qshell sbx rs sb-xxxxxxxxxxxx
```

2. 恢复多个沙箱
```
$ qshell sandbox resume sb-111111111111 sb-222222222222
```

3. 恢复所有已暂停的沙箱
```
$ qshell sandbox resume -a
$ qshell sbx rs -a
```

4. 按元数据过滤恢复
```
$ qshell sandbox resume -a -m env=staging
```
