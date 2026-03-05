# 简介
`sandbox template get`（别名 `gt`）查看模板的详细信息，包括构建记录。

# 格式
```
qshell sandbox template get <templateID>
qshell sbx tpl gt <templateID>
```

# 帮助文档
```
$ qshell sandbox template get -h
$ qshell sandbox template get --doc
```

# 鉴权
需要配置 `QINIU_API_KEY` 或 `E2B_API_KEY` 环境变量。

# 参数
- `templateID`：模板 ID（必填）

# 示例
```
$ qshell sandbox template get tmpl-xxxxxxxxxxxx
$ qshell sbx tpl gt tmpl-xxxxxxxxxxxx
```
