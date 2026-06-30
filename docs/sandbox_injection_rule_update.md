# 简介
`sandbox injection-rule update`（别名 `up`）更新指定的注入规则。

# 格式

```bash
qshell sandbox injection-rule update <ruleID> [--name <name>] [--type <openai|anthropic|gemini|qiniu|github|http>] [--api-key <apiKey>] [--base-url <baseURL>] [--headers <headers>] [--if-headers <headers>] [--if-queries <queries>]
qshell sbx ir up <ruleID> [--name <name>] [--type <openai|anthropic|gemini|qiniu|github|http>] [--api-key <apiKey>] [--base-url <baseURL>] [--headers <headers>] [--if-headers <headers>] [--if-queries <queries>]
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
- `--api-key`：新的 API Key；更新注入配置时必须与 `--type` 一同指定，且 `type=github` 时必填。注意：通过 CLI 传递密钥可能泄露到 Shell 历史或进程列表
- `--base-url`：新的基础 URL；更新注入配置时必须与 `--type` 一同指定；`github` 类型指定时 host 必须为 `github.com` 或 `api.github.com`
- `--headers`：新的自定义 HTTP 请求头，使用逗号分隔的 `key=value` 形式；更新时必须与 `--type http` 一同指定
- `--if-headers`：新的请求 Header 匹配条件，使用逗号分隔的 `key=value` 形式；仅当请求中已存在这些 Header 且值精确匹配时才注入
- `--if-queries`：新的请求 query 匹配条件，使用逗号分隔的 `key=value` 形式；仅当请求 URL 中已存在这些 query 参数且值精确匹配时才注入

说明：如果只传 `--api-key`、`--base-url`、`--headers`、`--if-headers` 或 `--if-queries` 而不传 `--type`，命令会报错，因为当前实现无法从已有规则自动推断要更新的注入类型。

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

更新自定义 HTTP 匹配条件：

```bash
$ qshell sandbox injection-rule update rule-xxxxxxxxxxxx --type http --base-url https://api.example.com --headers "Authorization=Bearer newtoken" --if-headers "X-Scope=demo" --if-queries "inject=true"
```

更新为 GitHub 凭证注入（token 通过 `--api-key` 传入）：

```bash
$ qshell sandbox injection-rule update rule-xxxxxxxxxxxx --type github --api-key ghp-new
```

限制 GitHub 凭证注入的匹配路径：

```bash
$ qshell sandbox injection-rule update rule-xxxxxxxxxxxx --type github --api-key ghp-new --base-url https://api.github.com/repos/qiniu/*
```
