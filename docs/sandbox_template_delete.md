# 简介
`sandbox template delete`（别名 `dl`）删除指定的沙箱模板。

# 格式
```
qshell sandbox template delete <templateID> [-y]
qshell sbx tpl dl <templateID> [-y]
```

# 帮助文档
```
$ qshell sandbox template delete -h
$ qshell sandbox template delete --doc
```

# 鉴权
需要配置 `QINIU_API_KEY` 或 `E2B_API_KEY` 环境变量。

# 参数
- `templateID`：模板 ID（必填）
- `-y, --yes`：跳过确认提示

# 示例
1. 删除模板（需确认）
```
$ qshell sandbox template delete tmpl-xxxxxxxxxxxx
```

2. 直接删除（跳过确认）
```
$ qshell sandbox template delete tmpl-xxxxxxxxxxxx -y
$ qshell sbx tpl dl tmpl-xxxxxxxxxxxx -y
```
