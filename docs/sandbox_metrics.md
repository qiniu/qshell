# 简介
`sandbox metrics` 查看沙箱的资源使用指标，包括 CPU、内存和磁盘使用情况。

# 格式
```
qshell sandbox metrics <sandboxID> [--format <pretty|json>]
```

# 帮助文档
```
$ qshell sandbox metrics -h
$ qshell sandbox metrics --doc
```

# 鉴权
需要配置 `QINIU_API_KEY` 或 `E2B_API_KEY` 环境变量。

# 参数
- `sandboxID`：沙箱 ID（必填）
- `--format`：输出格式，pretty（默认）或 json

# 示例
1. 查看沙箱指标
```
$ qshell sandbox metrics sb-xxxxxxxxxxxx
```

2. JSON 格式输出
```
$ qshell sandbox metrics sb-xxxxxxxxxxxx --format json
```
