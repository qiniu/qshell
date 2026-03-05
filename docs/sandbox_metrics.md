# 简介
`sandbox metrics`（别名 `mt`）查看沙箱的资源使用指标，包括 CPU、内存和磁盘使用情况。支持持续跟踪模式。输出采用 e2b 内联格式，标签带彩色高亮。

# 格式
```
qshell sandbox metrics <sandboxID> [--format <pretty|json>] [-f]
qshell sbx mt <sandboxID> [--format <pretty|json>] [-f]
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
- `-f, --follow`：持续跟踪指标输出，直到沙箱关闭

# 输出格式
Pretty 模式采用单行内联格式（匹配 e2b CLI）：
```
[2024-01-01T00:00:00Z]  CPU:  12.5% / 2 Cores | Memory:  256 MiB / 512 MiB | Disk:  100 MiB / 1024 MiB
```
其中 `CPU:`、`Memory:`、`Disk:` 标签以青色高亮显示。

# 示例
1. 查看沙箱指标
```
$ qshell sandbox metrics sb-xxxxxxxxxxxx
$ qshell sbx mt sb-xxxxxxxxxxxx
```

2. 持续跟踪指标
```
$ qshell sandbox metrics sb-xxxxxxxxxxxx -f
$ qshell sbx mt sb-xxxxxxxxxxxx -f
```

3. JSON 格式输出
```
$ qshell sandbox metrics sb-xxxxxxxxxxxx --format json
```
