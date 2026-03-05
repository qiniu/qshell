# 简介
`sandbox template`（别名 `tpl`）命令用于管理沙箱模板，支持列出、查看、删除、构建、发布模板以及初始化模板项目。

# 格式
```
qshell sandbox template <子命令>
qshell sbx tpl <子命令>
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
* list（ls）：列出模板
* get（gt）：查看模板详情
* delete（dl）：删除模板
* build（bd）：创建并构建模板
* builds（bds）：查看模板构建状态
* publish（pb）：发布模板（设为公开）
* unpublish（upb）：取消发布模板（设为私有）
* init（it）：初始化模板项目脚手架

# 示例
1. 列出所有模板
```
qshell sandbox template list
qshell sbx tpl ls
```

2. 查看模板详情
```
qshell sandbox template get tmpl-xxxxxxxxxxxx
qshell sbx tpl gt tmpl-xxxxxxxxxxxx
```

3. 删除模板
```
qshell sandbox template delete tmpl-xxxxxxxxxxxx -y
qshell sbx tpl dl tmpl-xxxxxxxxxxxx -y
```

4. 构建模板
```
qshell sandbox template build --name my-template --from-image ubuntu:22.04 --wait
qshell sbx tpl bd --name my-template --from-image ubuntu:22.04 --wait
```

5. 查看构建状态
```
qshell sandbox template builds tmpl-xxxxxxxxxxxx build-xxxxxxxxxxxx
qshell sbx tpl bds tmpl-xxxxxxxxxxxx build-xxxxxxxxxxxx
```

6. 发布/取消发布模板
```
qshell sandbox template publish tmpl-xxxxxxxxxxxx -y
qshell sandbox template unpublish tmpl-xxxxxxxxxxxx -y
```

7. 初始化模板项目
```
qshell sandbox template init
qshell sandbox template init --name my-template --language go
```
