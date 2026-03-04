# 简介
`sandbox kill` 终止一个或多个沙箱实例。

# 格式
```
qshell sandbox kill [sandboxIDs...] [--all] [--state <states>] [--metadata <key=value>]
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
- `--all`：终止所有匹配的沙箱
- `--state`：配合 --all 使用，按状态过滤（逗号分隔：running, paused）
- `--metadata`：配合 --all 使用，按元数据过滤（格式：key=value）

# 示例
1. 终止指定沙箱
```
$ qshell sandbox kill sb-xxxxxxxxxxxx
```

2. 终止多个沙箱
```
$ qshell sandbox kill sb-111111111111 sb-222222222222
```

3. 终止所有运行中的沙箱
```
$ qshell sandbox kill --all --state running
```
