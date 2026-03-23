# 简介
`sandbox transform-rule update`（别名 `up`）更新指定的转换规则。

# 格式
```
qshell sandbox transform-rule update <ruleID> [--name <name>] [--hosts <hosts>] [--headers <headers>] [--queries <queries>]
qshell sbx tr up <ruleID> [--name <name>] [--hosts <hosts>] [--headers <headers>] [--queries <queries>]
```

# 帮助文档
```
$ qshell sandbox transform-rule update -h
$ qshell sandbox transform-rule update --doc
```

# 鉴权
需要配置 `QINIU_API_KEY` 或 `E2B_API_KEY` 环境变量。

# 参数
- `ruleID`：转换规则 ID（必填）
- `--name`：新的规则名称
- `--hosts`：新的匹配域名列表，逗号分隔
- `--headers`：新的替换 Headers，逗号分隔的 key=value 对
- `--queries`：新的替换 Queries，逗号分隔的 key=value 对

至少需要提供一个更新参数。

# 示例
1. 更新规则名称
```
$ qshell sandbox transform-rule update rule-xxxxxxxxxxxx --name new-name
```

2. 更新匹配域名
```
$ qshell sandbox transform-rule update rule-xxxxxxxxxxxx --hosts api.new-domain.com
```

3. 更新多个字段
```
$ qshell sandbox transform-rule update rule-xxxxxxxxxxxx --name updated-rule --hosts api.example.com --headers "Authorization=Bearer newtoken"
```
