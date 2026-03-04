# 简介
`sandbox` 命令用于管理沙箱实例和模板，支持创建、连接、终止沙箱以及查看沙箱日志和指标。

# 格式
```
qshell sandbox <子命令>
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
* list：列出沙箱
* create：创建沙箱并连接终端
* connect：连接到已有沙箱终端
* kill：终止沙箱
* logs：查看沙箱日志
* metrics：查看沙箱资源指标
* template：管理沙箱模板

# 示例
1. 列出所有运行中的沙箱
```
qshell sandbox list --state running
```

2. 创建沙箱
```
qshell sandbox create my-template
```

3. 连接到沙箱
```
qshell sandbox connect sb-xxxxxxxxxxxx
```

4. 终止沙箱
```
qshell sandbox kill sb-xxxxxxxxxxxx
```
