# 简介
`sandbox transform-rule`（别名 `tr`）命令用于管理沙箱密钥转换规则，支持列出、查看、创建、更新和删除转换规则。

转换规则用于定义沙箱出站请求的自动拦截与替换，可在创建沙箱时通过规则 ID 引用。

# 格式
```
qshell sandbox transform-rule <子命令>
qshell sbx tr <子命令>
```

# 帮助文档
```
$ qshell sandbox transform-rule -h
$ qshell sandbox transform-rule --doc
```

# 鉴权
需要配置环境变量：
- `QINIU_API_KEY` 或 `E2B_API_KEY`：API 密钥（必填）
- `QINIU_SANDBOX_API_URL` 或 `E2B_API_URL`：API 服务地址（可选）

优先级：`QINIU_*` > `E2B_*`

# 子命令
transform-rule 的子命令有：
* list（ls）：列出所有转换规则
* get（gt）：查看转换规则详情
* create（cr）：创建转换规则
* update（up）：更新转换规则
* delete（dl）：删除转换规则

# 示例
1. 列出所有转换规则
```
qshell sandbox transform-rule list
qshell sbx tr ls
```

2. 查看转换规则详情
```
qshell sandbox transform-rule get rule-xxxxxxxxxxxx
qshell sbx tr gt rule-xxxxxxxxxxxx
```

3. 创建转换规则
```
qshell sandbox transform-rule create --name my-rule --hosts api.example.com --headers "Authorization=Bearer xxx"
qshell sbx tr cr --name my-rule --hosts api.example.com --headers "Authorization=Bearer xxx"
```

4. 更新转换规则
```
qshell sandbox transform-rule update rule-xxxxxxxxxxxx --name new-name
qshell sbx tr up rule-xxxxxxxxxxxx --name new-name
```

5. 删除转换规则
```
qshell sandbox transform-rule delete rule-xxxxxxxxxxxx -y
qshell sbx tr dl rule-xxxxxxxxxxxx -y
```
