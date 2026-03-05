# 简介
`sandbox kill`（别名 `kl`）终止一个或多个沙箱实例。使用 `--all` 时默认终止运行中的沙箱。

# 格式
```
qshell sandbox kill [sandboxIDs...] [-a] [-s <states>] [-m <metadata>]
qshell sbx kl [sandboxIDs...] [-a] [-s <states>] [-m <metadata>]
```

# 帮助文档
```
$ qshell sandbox kill -h
$ qshell sandbox kill --doc
```

# 鉴权
需要配置 `QINIU_API_KEY` 或 `E2B_API_KEY` 环境变量。

# 参数
- `sandboxIDs`：要终止的沙箱 ID 列表
- `-a, --all`：终止所有匹配的沙箱
- `-s, --state`：配合 -a/--all 使用，按状态过滤（逗号分隔：running, paused）。默认为 running
- `-m, --metadata`：配合 -a/--all 使用，按元数据过滤（格式：key1=value1,key2=value2）

# 示例
1. 终止指定沙箱
```
$ qshell sandbox kill sb-xxxxxxxxxxxx
$ qshell sbx kl sb-xxxxxxxxxxxx
```

2. 终止多个沙箱
```
$ qshell sandbox kill sb-111111111111 sb-222222222222
```

3. 终止所有运行中的沙箱
```
$ qshell sandbox kill -a
$ qshell sbx kl -a
```

4. 终止所有暂停的沙箱
```
$ qshell sandbox kill -a -s paused
```

5. 按元数据过滤终止
```
$ qshell sandbox kill -a -m user=alice
```
