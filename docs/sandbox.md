# 简介
`sandbox`（别名 `sbx`）命令用于管理沙箱实例、模板、注入规则和已挂载资源，支持创建、连接、执行命令、暂停、恢复、终止沙箱，以及查看沙箱日志、指标和资源信息。

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
sandbox 命令同时支持两套凭据，按子命令所需自动选择：

| 子命令 | 鉴权方式 |
| --- | --- |
| `sandbox` 实例操作（list/create/connect/exec/...）、`resource`、`template` 系列 | API Key |
| `injection-rule` 系列（list/create/get/update/delete） | AK/SK |

## API Key
配置以下任一环境变量（也支持当前目录下的 `.env` 文件），优先级 `QINIU_*` > `E2B_*`：
- `QINIU_API_KEY` 或 `E2B_API_KEY`：API 密钥
- `QINIU_SANDBOX_API_URL` 或 `E2B_API_URL`：API 服务地址（可选）

## AK/SK
按以下优先级解析，**qshell 账号优先于环境变量**：
1. `qshell user` 配置的当前账号（推荐方式）
2. 环境变量 `QINIU_ACCESS_KEY` / `QINIU_SECRET_KEY`

> 同时配置 API Key 与 AK/SK 时互不冲突，每个端点会按 SDK 既定方式选择对应鉴权。当某子命令所需的凭据缺失时，会给出明确报错。

# 子命令
sandbox 的子命令有：
* list（ls）：列出沙箱
* create（cr）：创建沙箱并连接终端
* connect（cn）：连接到已有沙箱终端
* kill（kl）：终止沙箱
* pause（ps）：暂停沙箱
* resume（rs）：恢复已暂停的沙箱
* exec（ex）：在沙箱中执行命令
* logs（lg）：查看沙箱日志
* metrics（mt）：查看沙箱资源指标
* resource（resources）：管理沙箱已挂载资源
* template（tpl）：管理沙箱模板
* injection-rule（ir）：管理沙箱注入规则

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
qshell sandbox create my-template \
  --inline-injection 'type=openai,api-key=sk-xxx' \
  --inline-injection 'type=http,base-url=https://api.example.com,headers=Authorization=Bearer token;X-Env=prod'
```

5. 管理注入规则
```
qshell sandbox injection-rule list
qshell sbx ir cr --name openai-default --type openai --api-key sk-xxx
```

6. 连接到沙箱
```
qshell sandbox connect sb-xxxxxxxxxxxx
qshell sbx cn sb-xxxxxxxxxxxx
```

7. 在沙箱中执行命令
```
qshell sandbox exec sb-xxxxxxxxxxxx -- ls -la
qshell sbx ex sb-xxxxxxxxxxxx -- ls -la
```

8. 暂停并恢复沙箱
```
qshell sandbox pause sb-xxxxxxxxxxxx
qshell sandbox resume sb-xxxxxxxxxxxx
```

9. 终止沙箱
```
qshell sandbox kill sb-xxxxxxxxxxxx
qshell sbx kl sb-xxxxxxxxxxxx
```

10. 查询沙箱已挂载资源
```
qshell sandbox resource list sb-xxxxxxxxxxxx
qshell sbx resource ls sb-xxxxxxxxxxxx --format json
```

11. 更新 Git 仓库资源令牌
```
qshell sandbox resource update sb-xxxxxxxxxxxx res-xxxxxxxxxxxx --token ghp-xxx
qshell sbx resource up sb-xxxxxxxxxxxx res-xxxxxxxxxxxx -t ghp-xxx
```

12. 构建模板
```
qshell sandbox template build --name my-template --from-image ubuntu:22.04 --wait
qshell sbx tpl bd --name my-template --from-image ubuntu:22.04 --wait

# 创建模板时指定 20 GiB 构建磁盘（--disk-size 单位为 MiB）
qshell sandbox template build --name my-template --from-image ubuntu:22.04 --disk-size 20480 --wait
qshell sbx tpl bd --name my-template --from-image ubuntu:22.04 --disk-size 20480 --wait

# 使用当前目录的 qshell.sandbox.toml 构建模板
qshell sandbox template build --wait
qshell sbx tpl bd --wait
```
