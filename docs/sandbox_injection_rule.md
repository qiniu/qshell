# 简介
`sandbox injection-rule`（别名 `ir`）命令用于管理沙箱注入规则，支持列出、查看、创建、更新和删除注入规则。

# 格式

```bash
qshell sandbox injection-rule <子命令>
qshell sbx ir <子命令>
```

# 帮助文档

```bash
$ qshell sandbox injection-rule -h
$ qshell sandbox injection-rule --doc
```

# 子命令

`injection-rule` 的子命令有：

- `list`（`ls`）：列出所有注入规则
- `get`（`gt`）：查看指定注入规则详情
- `create`（`cr`）：创建新的注入规则
- `update`（`up`）：更新已有注入规则
- `delete`（`dl`）：删除一个或多个注入规则

# 示例

列出所有注入规则：

```bash
qshell sandbox injection-rule list
```

查看指定注入规则：

```bash
qshell sandbox injection-rule get rule-xxxxxxxxxxxx
```

创建 OpenAI 注入规则：

```bash
qshell sandbox injection-rule create --name openai-default --type openai --api-key sk-xxx
```

更新自定义 HTTP 注入规则：

```bash
qshell sandbox injection-rule update rule-xxxxxxxxxxxx --type http --base-url https://api.example.com --headers "Authorization=Bearer newtoken"
```

删除注入规则：

```bash
qshell sandbox injection-rule delete rule-xxxxxxxxxxxx -y
```
