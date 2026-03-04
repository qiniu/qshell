# 简介
`sandbox create` 创建一个新的沙箱实例并连接到其终端。当终端会话结束时，沙箱将被自动终止。

# 格式
```
qshell sandbox create [template] [--timeout <seconds>]
```

# 帮助文档
```
$ qshell sandbox create -h
$ qshell sandbox create --doc
```

# 鉴权
需要配置 `QINIU_API_KEY` 或 `E2B_API_KEY` 环境变量。

# 参数
- `template`：模板 ID（必填）
- `--timeout`：沙箱超时时间（秒）

# 示例
1. 使用模板创建沙箱
```
$ qshell sandbox create my-template
```

2. 指定超时时间
```
$ qshell sandbox create my-template --timeout 300
```
