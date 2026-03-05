# 简介
`sandbox logs`（别名 `lg`）查看沙箱的日志。支持按日志级别和 logger 过滤，支持持续跟踪模式。日志输出带有彩色级别标签，内部字段（如 traceID、sandboxID）自动剥离。

# 格式
```
qshell sandbox logs <sandboxID> [--level <level>] [--limit <N>] [--format <pretty|json>] [-f] [--loggers <loggers>]
qshell sbx lg <sandboxID> [--level <level>] [--limit <N>] [--format <pretty|json>] [-f] [--loggers <loggers>]
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
- `--level`：按日志级别过滤（DEBUG, INFO, WARN, ERROR）。默认为 INFO。更高级别的日志也会显示
- `--limit`：返回的最大日志条数
- `--format`：输出格式，pretty（默认）或 json
- `-f, --follow`：持续跟踪日志输出，直到沙箱关闭
- `--loggers`：按 logger 名称前缀过滤（逗号分隔）

# 输出特性
- **彩色级别标签**：DEBUG（白色）、INFO（绿色）、WARN（黄色）、ERROR（红色）
- **字段剥离**：自动隐藏内部字段（traceID、instanceID、teamID、source、service、envID、sandboxID、source_type）
- **Logger 名称清理**：自动去除 "Svc" 后缀
- **异步结束检测**：Follow 模式下通过后台 goroutine 检测沙箱是否结束，不阻塞日志轮询

# 示例
1. 查看沙箱日志（默认 INFO 及以上级别）
```
$ qshell sandbox logs sb-xxxxxxxxxxxx
$ qshell sbx lg sb-xxxxxxxxxxxx
```

2. 查看所有级别的日志
```
$ qshell sandbox logs sb-xxxxxxxxxxxx --level DEBUG
```

3. 只查看错误日志
```
$ qshell sandbox logs sb-xxxxxxxxxxxx --level ERROR
```

4. 持续跟踪日志
```
$ qshell sandbox logs sb-xxxxxxxxxxxx -f
$ qshell sbx lg sb-xxxxxxxxxxxx -f
```

5. 按 logger 过滤
```
$ qshell sandbox logs sb-xxxxxxxxxxxx --loggers envd,proxy
```

6. JSON 格式输出
```
$ qshell sandbox logs sb-xxxxxxxxxxxx --format json
```
