# 简介
`sandbox transform-rule list`（别名 `ls`）列出所有转换规则。

# 格式
```
qshell sandbox transform-rule list [--format <pretty|json>]
qshell sbx tr ls [--format <pretty|json>]
```

# 帮助文档
```
$ qshell sandbox transform-rule list -h
$ qshell sandbox transform-rule list --doc
```

# 鉴权
需要配置 `QINIU_API_KEY` 或 `E2B_API_KEY` 环境变量。

# 参数
- `--format`：输出格式，pretty（默认）或 json

# 示例
1. 列出所有转换规则
```
$ qshell sandbox transform-rule list
$ qshell sbx tr ls
```

2. JSON 格式输出
```
$ qshell sandbox transform-rule list --format json
```
