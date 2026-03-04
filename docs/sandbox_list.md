# 简介
`sandbox list` 列出沙箱实例，支持按状态和元数据过滤。

# 格式
```
qshell sandbox list [--state <states>] [--metadata <key=value>] [--limit <N>] [--format <pretty|json>]
```

# 帮助文档
```
$ qshell sandbox list -h
$ qshell sandbox list --doc
```

# 鉴权
需要配置 `QINIU_API_KEY` 或 `E2B_API_KEY` 环境变量。

# 参数
- `--state`：按状态过滤，逗号分隔（可选值：running, paused）
- `--metadata`：按元数据过滤（格式：key=value）
- `--limit`：返回的最大数量
- `--format`：输出格式，pretty（默认）或 json

# 示例
1. 列出所有沙箱
```
$ qshell sandbox list
```

2. 列出运行中的沙箱
```
$ qshell sandbox list --state running
```

3. 以 JSON 格式输出
```
$ qshell sandbox list --format json
```
