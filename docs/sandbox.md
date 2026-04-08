# 简介
`sandbox`（别名 `sbx`）命令用于管理沙箱实例和模板，支持创建、连接、终止沙箱以及查看沙箱日志和指标。

# 格式
```
qshell sandbox <子命令>
qshell sbx <子命令>
```

# 帮助文档
可以在命令行输入如下命令获取帮助文档：
```
// 简单描述
$ qshell sandbox -h

// 详细文档（此文档）
$ qshell sandbox --doc
```

# 鉴权
需要配置环境变量：
- `QINIU_API_KEY` 或 `E2B_API_KEY`：API 密钥（必填）
- `QINIU_SANDBOX_API_URL` 或 `E2B_API_URL`：API 服务地址（可选）

优先级：`QINIU_*` > `E2B_*`

# 子命令
sandbox 的子命令有：
* list（ls）：列出沙箱
* create（cr）：创建沙箱并连接终端
* injection-rule（ir）：管理沙箱注入规则
* connect（cn）：连接到已有沙箱终端
* kill（kl）：终止沙箱
* logs（lg）：查看沙箱日志
* metrics（mt）：查看沙箱资源指标
* template（tpl）：管理沙箱模板

# 示例
1. 列出所有运行中的沙箱
```
qshell sandbox list --state running
qshell sbx ls -s running
```

2. 创建沙箱
```
qshell sandbox create my-template
qshell sbx cr my-template
```

3. 创建沙箱时附加已存在的注入规则
```
qshell sandbox create my-template --injection-rule rule-openai --injection-rule rule-http
qshell sbx cr my-template --injection-rule rule-openai --injection-rule rule-http
```

4. 创建沙箱时附加内联注入配置
```
qshell sandbox create my-template --inline-injection 'type=openai,api-key=sk-xxx'
qshell sbx cr my-template --inline-injection 'type=http,base-url=https://api.example.com,headers=Authorization=Bearer token,X-Env=prod'
```

5. 管理注入规则
```
qshell sandbox injection-rule list
qshell sbx ir create --name openai-default --type openai --api-key sk-xxx
```

6. 连接到沙箱
```
qshell sandbox connect sb-xxxxxxxxxxxx
qshell sbx cn sb-xxxxxxxxxxxx
```

7. 终止沙箱
```
qshell sandbox kill sb-xxxxxxxxxxxx
qshell sbx kl sb-xxxxxxxxxxxx
```
