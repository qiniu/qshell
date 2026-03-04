# 简介
`sandbox template` 命令用于管理沙箱模板，支持列出、查看、删除、构建模板以及查看构建状态。

# 格式
```
qshell sandbox template <子命令>
```

# 帮助文档
```
$ qshell sandbox template -h
$ qshell sandbox template --doc
```

# 鉴权
需要配置环境变量：
- `QINIU_API_KEY` 或 `E2B_API_KEY`：API 密钥（必填）
- `QINIU_SANDBOX_API_URL` 或 `E2B_API_URL`：API 服务地址（可选）

优先级：`QINIU_*` > `E2B_*`

# 子命令
template 的子命令有：
* list：列出模板
* get：查看模板详情
* delete：删除模板
* build：创建并构建模板
* builds：查看模板构建状态

# 示例
1. 列出所有模板
```
qshell sandbox template list
```

2. 查看模板详情
```
qshell sandbox template get tmpl-xxxxxxxxxxxx
```

3. 删除模板
```
qshell sandbox template delete tmpl-xxxxxxxxxxxx -y
```

4. 构建模板
```
qshell sandbox template build --name my-template --from-image ubuntu:22.04 --wait
```

5. 查看构建状态
```
qshell sandbox template builds tmpl-xxxxxxxxxxxx build-xxxxxxxxxxxx
```
