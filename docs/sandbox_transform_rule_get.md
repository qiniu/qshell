# 简介
`sandbox transform-rule get`（别名 `gt`）查看转换规则的详细信息。

# 格式
```
qshell sandbox transform-rule get <ruleID>
qshell sbx tr gt <ruleID>
```

# 帮助文档
```
$ qshell sandbox transform-rule get -h
$ qshell sandbox transform-rule get --doc
```

# 鉴权
需要配置 `QINIU_API_KEY` 或 `E2B_API_KEY` 环境变量。

# 参数
- `ruleID`：转换规则 ID（必填）

# 示例
```
$ qshell sandbox transform-rule get rule-xxxxxxxxxxxx
$ qshell sbx tr gt rule-xxxxxxxxxxxx
```
