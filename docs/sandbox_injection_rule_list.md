`sandbox injection-rule list`（别名 `ls`）列出所有注入规则。

## 命令格式

```bash
qshell sandbox injection-rule list [--format <pretty|json>]
```

## 查看帮助

```bash
$ qshell sandbox injection-rule list -h
$ qshell sandbox injection-rule list --doc
```

## 参数说明

- `--format`：输出格式，支持 `pretty` 和 `json`，默认 `pretty`

## 使用示例

默认表格输出：

```bash
$ qshell sandbox injection-rule list
```

JSON 输出：

```bash
$ qshell sandbox injection-rule list --format json
```
