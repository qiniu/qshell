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
- `-m, --metadata`：配合 -a/--all 使用，按元数据过滤（格式：key1=value1,key2=value2，键和值均不能为空，键不能重复）。显式传入无效条件时，命令以非零状态退出，不发送沙箱请求。

# 示例
`--all` 先完整枚举并去重，再执行恢复；枚举失败时不修改第一页资源。结束后输出成功和失败数量，任一请求失败返回非零退出码。

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

恢复使用 SDK 的 v2 Connect 接口，未指定超时时间时采用服务端默认 300 秒。仅保存文件系统的暂停恢复后会冷启动。
