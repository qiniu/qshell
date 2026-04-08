# 简介
`sandbox create`（别名 `cr`）创建一个新的沙箱实例并连接到其终端。当终端会话结束时，沙箱将被自动终止。

使用 `--detach` 模式时，创建沙箱后不连接终端，沙箱保持存活直到超时。

沙箱通过 keep-alive 机制保持存活，终端连接期间会自动续命，无需手动设置超时。

# 格式
```
qshell sandbox create <template> [-t <seconds>] [-d] [-m <metadata>] [--injection-rule <ruleID>...] [--inline-injection <spec>...]
qshell sbx cr <template> [-t <seconds>] [-d] [-m <metadata>] [--injection-rule <ruleID>...] [--inline-injection <spec>...]
```

# 帮助文档
```
$ qshell sandbox create -h
$ qshell sandbox create --doc
```

# 鉴权
需要配置 `QINIU_API_KEY` 或 `E2B_API_KEY` 环境变量，或在当前目录 `.env` 文件中设置。

# 参数
- `template`：模板 ID（必填）
- `-t, --timeout`：沙箱超时时间（秒）
- `-d, --detach`：创建沙箱但不连接终端，沙箱保持存活直到超时
- `-m, --metadata`：元数据键值对（格式：key1=value1,key2=value2）
- `--injection-rule`：创建沙箱时附加的注入规则 ID，可多次指定
- `--inline-injection`：创建沙箱时附加的内联注入配置，可多次指定，格式为 `type=<type>,api-key=<key>,base-url=<url>,headers=<k1=v1,k2=v2>`

内联注入说明：
- `type` 支持 `openai`、`anthropic`、`gemini`、`http`
- `api-key` 可用于 `openai`、`anthropic`、`gemini`
- `base-url` 可用于覆盖默认目标地址；`type=http` 时必填
- `headers` 仅用于 `type=http`，多个请求头使用逗号分隔，例如 `headers=Authorization=Bearer token,X-Env=prod`
- `headers` 建议放在内联配置的最后一个字段，避免和其他顶层字段分隔符冲突

# 示例
1. 创建沙箱
```
$ qshell sandbox create my-template
$ qshell sbx cr my-template
```

2. 设置超时时间
```
$ qshell sandbox create my-template --timeout 300
$ qshell sbx cr my-template -t 300
```

3. 分离模式创建（不连接终端，沙箱存活 5 分钟）
```
$ qshell sandbox create my-template -t 300 -d
$ qshell sbx cr my-template -t 300 -d
```

4. 添加元数据
```
$ qshell sandbox create my-template -m env=dev,team=backend
$ qshell sbx cr my-template -m env=dev,team=backend
```

5. 创建时附加注入规则
```
$ qshell sandbox create my-template --injection-rule rule-openai --injection-rule rule-http
$ qshell sbx cr my-template --injection-rule rule-openai --injection-rule rule-http
```

6. 创建时附加内联注入配置
```
$ qshell sandbox create my-template --inline-injection 'type=openai,api-key=sk-xxx' --inline-injection 'type=http,base-url=https://api.example.com,headers=Authorization=Bearer token,X-Env=prod'
$ qshell sbx cr my-template --inline-injection 'type=gemini,api-key=sk-gem'
```
