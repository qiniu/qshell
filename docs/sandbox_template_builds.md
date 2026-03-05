# 简介
`sandbox template builds`（别名 `bds`）查看模板的构建状态和日志。

# 格式
```
qshell sandbox template builds <templateID> <buildID>
qshell sbx tpl bds <templateID> <buildID>
```

# 帮助文档
```
$ qshell sandbox template builds -h
$ qshell sandbox template builds --doc
```

# 鉴权
需要配置 `QINIU_API_KEY` 或 `E2B_API_KEY` 环境变量。

# 参数
- `templateID`：模板 ID（必填）
- `buildID`：构建 ID（必填）

# 示例
```
$ qshell sandbox template builds tmpl-xxxxxxxxxxxx build-xxxxxxxxxxxx
$ qshell sbx tpl bds tmpl-xxxxxxxxxxxx build-xxxxxxxxxxxx
```
