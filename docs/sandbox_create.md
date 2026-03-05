# 简介
`sandbox create`（别名 `cr`）创建一个新的沙箱实例并连接到其终端。当终端会话结束时，沙箱将被自动终止。

沙箱通过 keep-alive 机制保持存活，终端连接期间会自动续命，无需手动设置超时。

# 格式
```
qshell sandbox create [template]
qshell sbx cr [template]
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

# 示例
```
$ qshell sandbox create my-template
$ qshell sbx cr my-template
```
