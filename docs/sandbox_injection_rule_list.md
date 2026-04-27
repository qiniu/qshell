# 简介
`sandbox injection-rule list`（别名 `ls`）列出所有注入规则。

# 格式

```bash
qshell sandbox injection-rule list [--format <pretty|json>]
qshell sbx ir ls [--format <pretty|json>]
```

# 帮助文档

```bash
$ qshell sandbox injection-rule list -h
$ qshell sandbox injection-rule list --doc
```

# 参数

- `--format`：输出格式，支持 `pretty` 和 `json`，默认 `pretty`

# 示例

默认表格输出：

```bash
$ qshell sandbox injection-rule list
```

JSON 输出：

```bash
$ qshell sandbox injection-rule list --format json
```
