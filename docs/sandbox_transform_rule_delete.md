# 简介
`sandbox transform-rule delete`（别名 `dl`）删除一个或多个转换规则。支持变参和交互式多选。

# 格式
```
qshell sandbox transform-rule delete [ruleIDs...] [-y] [-s]
qshell sbx tr dl [ruleIDs...] [-y] [-s]
```

# 帮助文档
```
$ qshell sandbox transform-rule delete -h
$ qshell sandbox transform-rule delete --doc
```

# 鉴权
需要配置 `QINIU_API_KEY` 或 `E2B_API_KEY` 环境变量。

# 参数
- `ruleIDs`：一个或多个转换规则 ID（与 `--select` 二选一）
- `-y, --yes`：跳过确认提示
- `-s, --select`：交互式选择规则进行删除

# 示例
1. 删除单个规则（需确认）
```
$ qshell sandbox transform-rule delete rule-xxxxxxxxxxxx
```

2. 直接删除（跳过确认）
```
$ qshell sandbox transform-rule delete rule-xxxxxxxxxxxx -y
$ qshell sbx tr dl rule-xxxxxxxxxxxx -y
```

3. 删除多个规则
```
$ qshell sandbox transform-rule delete rule-aaa rule-bbb -y
```

4. 交互式选择删除
```
$ qshell sandbox transform-rule delete -s
$ qshell sbx tr dl -s
```
