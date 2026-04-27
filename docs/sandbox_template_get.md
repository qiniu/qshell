# 简介
`sandbox template get`（别名 `gt`）查看模板的详细信息，包括构建记录。

# 格式
```
qshell sandbox template get [templateID]
qshell sbx tpl gt [templateID]
```

# 帮助文档
```
$ qshell sandbox template get -h
$ qshell sandbox template get --doc
```

# 鉴权
需要配置 `QINIU_API_KEY` 或 `E2B_API_KEY` 环境变量。

# 参数
- `templateID`：模板 ID。未传入时，自动读取当前目录 `qshell.sandbox.toml` 中的 `template_id`

# 示例
1. 通过参数查看模板详情
```
$ qshell sandbox template get tmpl-xxxxxxxxxxxx
$ qshell sbx tpl gt tmpl-xxxxxxxxxxxx
```

2. 读取当前目录 `qshell.sandbox.toml` 中的 `template_id`
```
$ qshell sandbox template get
$ qshell sbx tpl gt
```
