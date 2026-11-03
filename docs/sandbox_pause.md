# 简介
`sandbox pause`（别名 `ps`）暂停一个或多个运行中的沙箱。暂停后的沙箱可以通过 `sandbox resume` 恢复。

# 格式
```
qshell sandbox pause [sandboxIDs...] [-a] [-s <states>] [-m <metadata>]
qshell sbx ps [sandboxIDs...] [-a] [-s <states>] [-m <metadata>]
```

# 帮助文档
```
$ qshell sandbox pause -h
$ qshell sandbox pause --doc
```

# 鉴权
需要配置 `QINIU_API_KEY` 或 `E2B_API_KEY` 环境变量，或在当前目录 `.env` 文件中设置。

# 参数
- `sandboxIDs`：要暂停的沙箱 ID 列表
- `-a, --all`：暂停所有匹配的沙箱
- `-s, --state`：配合 -a/--all 使用，按状态过滤（逗号分隔：running, paused）。默认为 running
- `-m, --metadata`：配合 -a/--all 使用，按元数据过滤（格式：key1=value1,key2=value2）

# 示例
1. 暂停指定沙箱
```
$ qshell sandbox pause sb-xxxxxxxxxxxx
$ qshell sbx ps sb-xxxxxxxxxxxx
```

2. 暂停多个沙箱
```
$ qshell sandbox pause sb-111111111111 sb-222222222222
```

3. 暂停所有运行中的沙箱
```
$ qshell sandbox pause -a
$ qshell sbx ps -a
```

4. 按元数据过滤暂停
```
$ qshell sandbox pause -a -m user=alice
```

# 仅保存文件系统
`--memory` 默认 true，保留进程和内存状态。使用 `--memory=false` 时仅保存文件系统，恢复会冷启动，进程和内存状态丢失，且不能通过流量自动恢复：
```shell
qshell sandbox pause <sandboxID> --memory=false
```
此参数同样适用于批量暂停和 `--all`。
