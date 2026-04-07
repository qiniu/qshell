`sandbox injection-rule create`（别名 `cr`）创建一个新的注入规则。

## 命令格式

```bash
qshell sandbox injection-rule create --name <name> --type <openai|anthropic|gemini|http> [--api-key <apiKey>] [--base-url <baseURL>] [--headers <headers>]
```

## 查看帮助

```bash
$ qshell sandbox injection-rule create -h
$ qshell sandbox injection-rule create --doc
```

## 参数说明

- `--name`：规则名称，必填，同一用户下唯一
- `--type`：注入类型，必填，支持 `openai`、`anthropic`、`gemini`、`http`
- `--api-key`：`openai`、`anthropic`、`gemini` 类型使用的 API Key
- `--base-url`：覆盖默认目标地址，或 `http` 类型的目标基础 URL
- `--headers`：`http` 类型的请求头，使用逗号分隔的 `key=value` 形式

## 使用示例

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
