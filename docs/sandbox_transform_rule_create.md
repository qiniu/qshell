# 简介
`sandbox transform-rule create`（别名 `cr`）创建一个新的转换规则。

# 格式
```
qshell sandbox transform-rule create --name <name> [--hosts <hosts>] [--headers <headers>] [--queries <queries>]
qshell sbx tr cr --name <name> [--hosts <hosts>] [--headers <headers>] [--queries <queries>]
```

# 帮助文档
```
$ qshell sandbox transform-rule create -h
$ qshell sandbox transform-rule create --doc
```

# 鉴权
需要配置 `QINIU_API_KEY` 或 `E2B_API_KEY` 环境变量。

# 参数
- `--name`：规则名称（必填），同一用户下唯一
- `--hosts`：匹配条件，逗号分隔的域名列表（如 `api.example.com,cdn.example.com`）
- `--headers`：替换的 HTTP Headers，逗号分隔的 key=value 对（如 `Authorization=Bearer xxx,X-Custom=val`）
- `--queries`：替换的 URL Query 参数，逗号分隔的 key=value 对（如 `token=abc,version=2`）

# 示例
1. 创建基本规则
```
$ qshell sandbox transform-rule create --name my-rule --hosts api.example.com
```

2. 创建带 Headers 替换的规则
```
$ qshell sandbox transform-rule create --name api-auth --hosts api.example.com --headers "Authorization=Bearer token123"
```

3. 创建带 Headers 和 Queries 替换的规则
```
$ qshell sandbox transform-rule create --name full-rule --hosts api.example.com,cdn.example.com --headers "Authorization=Bearer xxx" --queries "token=abc"
```
