# 简介
`sandbox logs` 查看沙箱的日志。

# 格式
```
qshell sandbox logs <sandboxID> [--level <level>] [--limit <N>] [--format <pretty|json>]
```

# 帮助文档
```
$ qshell sandbox logs -h
$ qshell sandbox logs --doc
```

# 鉴权
需要配置 `QINIU_API_KEY` 或 `E2B_API_KEY` 环境变量。

# 参数
- `sandboxID`：沙箱 ID（必填）
- `--level`：按日志级别过滤（INFO, WARN, ERROR, DEBUG）
- `--limit`：返回的最大日志条数
- `--format`：输出格式，pretty（默认）或 json

# 示例
1. 查看沙箱日志
```
$ qshell sandbox logs sb-xxxxxxxxxxxx
```

2. 过滤 INFO 级别日志
```
$ qshell sandbox logs sb-xxxxxxxxxxxx --level INFO
```

3. JSON 格式输出
```
$ qshell sandbox logs sb-xxxxxxxxxxxx --format json
```
