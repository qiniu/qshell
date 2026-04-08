`sandbox injection-rule delete`（别名 `dl`）删除一个或多个注入规则。支持变参和交互式多选。

## 命令格式

```bash
qshell sandbox injection-rule delete [ruleIDs...] [-y] [-s]
```

## 查看帮助

```bash
$ qshell sandbox injection-rule delete -h
$ qshell sandbox injection-rule delete --doc
```

## 参数说明

- `ruleIDs...`：一个或多个注入规则 ID
- `-y, --yes`：跳过确认
- `-s, --select`：交互式选择规则进行删除

## 使用示例

删除单个规则：

```bash
$ qshell sandbox injection-rule delete rule-xxxxxxxxxxxx
```

不加 `-y` 时，命令会交互式询问确认。

跳过确认直接删除：

```bash
$ qshell sandbox injection-rule delete rule-xxxxxxxxxxxx -y
```

删除多个规则：

```bash
$ qshell sandbox injection-rule delete rule-aaa rule-bbb -y
```

交互式选择删除：

```bash
$ qshell sandbox injection-rule delete -s
```
