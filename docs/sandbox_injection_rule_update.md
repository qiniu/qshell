# 简介
`sandbox injection-rule update`（别名 `up`）更新指定的注入规则。

# 格式

```bash
qshell sandbox injection-rule update <ruleID> [--name <name>] [--type <openai|anthropic|gemini|qiniu|http>] [--api-key <apiKey>] [--base-url <baseURL>] [--headers <headers>]
qshell sbx ir up <ruleID> [--name <name>] [--type <openai|anthropic|gemini|qiniu|http>] [--api-key <apiKey>] [--base-url <baseURL>] [--headers <headers>]
```

# 帮助文档

```bash
$ qshell sandbox injection-rule update -h
$ qshell sandbox injection-rule update --doc
```

# 参数

- `ruleID`：注入规则 ID
- `--name`：新的规则名称
- `--type`：新的注入类型；当需要更新注入配置时必须指定
- `--api-key`：新的 API Key；更新注入配置时必须与 `--type` 一同指定。注意：通过 CLI 传递密钥可能泄露到 Shell 历史或进程列表
- `--base-url`：新的基础 URL；更新注入配置时必须与 `--type` 一同指定
- `--headers`：新的自定义 HTTP 请求头，使用逗号分隔的 `key=value` 形式；更新时必须与 `--type http` 一同指定

说明：如果只传 `--api-key`、`--base-url` 或 `--headers` 而不传 `--type`，命令会报错，因为当前实现无法从已有规则自动推断要更新的注入类型。

# 示例

更新规则名称：

```bash
$ qshell sandbox injection-rule update rule-xxxxxxxxxxxx --name new-name
```

更新为 Gemini 注入规则：

```bash
$ qshell sandbox injection-rule update rule-xxxxxxxxxxxx --type gemini --api-key sk-gem --base-url https://gemini-proxy.example.com
```

更新为七牛 AI API 注入规则：

```bash
$ qshell sandbox injection-rule update rule-xxxxxxxxxxxx --type qiniu --api-key ak-new
```

更新自定义 HTTP 请求头：

```bash
$ qshell sandbox injection-rule update rule-xxxxxxxxxxxx --type http --base-url https://api.example.com --headers "Authorization=Bearer newtoken"
```
