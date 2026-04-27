# 简介
`sandbox injection-rule create`（别名 `cr`）创建一个新的注入规则。

# 格式

```bash
qshell sandbox injection-rule create --name <name> --type <openai|anthropic|gemini|qiniu|http> [--api-key <apiKey>] [--base-url <baseURL>] [--headers <headers>]
qshell sbx ir cr --name <name> --type <openai|anthropic|gemini|qiniu|http> [--api-key <apiKey>] [--base-url <baseURL>] [--headers <headers>]
```

# 帮助文档

```bash
$ qshell sandbox injection-rule create -h
$ qshell sandbox injection-rule create --doc
```

# 参数

- `--name`：规则名称，必填，同一用户下唯一
- `--type`：注入类型，必填，支持 `openai`、`anthropic`、`gemini`、`qiniu`、`http`
- `--api-key`：`openai`、`anthropic`、`gemini`、`qiniu` 类型使用的 API Key。注意：通过 CLI 传递密钥可能泄露到 Shell 历史或进程列表
- `--base-url`：覆盖默认目标地址，或 `http` 类型的目标基础 URL；`qiniu` 默认为 `api.qnaigc.com`
- `--headers`：`http` 类型的请求头，使用逗号分隔的 `key=value` 形式

# 示例

创建 OpenAI 注入规则：

```bash
$ qshell sandbox injection-rule create --name openai-default --type openai --api-key sk-xxx
```

创建 Anthropic 注入规则并指定代理地址：

```bash
$ qshell sandbox injection-rule create --name anthropic-proxy --type anthropic --api-key sk-ant --base-url https://anthropic-proxy.example.com
```

创建自定义 HTTP 注入规则：

```bash
$ qshell sandbox injection-rule create --name api-auth --type http --base-url https://api.example.com --headers "Authorization=Bearer token123,X-Env=prod"
```

创建七牛 AI API 注入规则：

```bash
$ qshell sandbox injection-rule create --name qiniu-ai --type qiniu --api-key ak-xxx
```
