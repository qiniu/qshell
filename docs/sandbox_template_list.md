# 简介
`sandbox template list`（别名 `ls`）列出所有沙箱模板。

# 格式
```
qshell sandbox template list [--format <pretty|json>]
qshell sbx tpl ls [--format <pretty|json>]
```

# 帮助文档
```
$ qshell sandbox template list -h
$ qshell sandbox template list --doc
```

# 鉴权
需要配置 `QINIU_API_KEY` 或 `E2B_API_KEY` 环境变量。

# 参数
- `--format`：输出格式，pretty（默认）或 json

# 示例
1. 列出所有模板
```
$ qshell sandbox template list
$ qshell sbx tpl ls
```

2. JSON 格式输出
```
$ qshell sandbox template list --format json
```
