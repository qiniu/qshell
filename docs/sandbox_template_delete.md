# 简介
`sandbox template delete`（别名 `dl`）删除一个或多个沙箱模板。支持变参和交互式多选。

# 格式
```
qshell sandbox template delete [templateIDs...] [-y] [-s]
qshell sbx tpl dl [templateIDs...] [-y] [-s]
```

# 帮助文档
```
$ qshell sandbox template delete -h
$ qshell sandbox template delete --doc
```

# 鉴权
需要配置 `QINIU_API_KEY` 或 `E2B_API_KEY` 环境变量。

# 参数
- `templateIDs`：一个或多个模板 ID（与 `--select` 二选一）
- `-y, --yes`：跳过确认提示
- `-s, --select`：交互式选择模板进行删除

# 示例
1. 删除单个模板（需确认）
```
$ qshell sandbox template delete tmpl-xxxxxxxxxxxx
```

2. 直接删除（跳过确认）
```
$ qshell sandbox template delete tmpl-xxxxxxxxxxxx -y
$ qshell sbx tpl dl tmpl-xxxxxxxxxxxx -y
```

3. 删除多个模板
```
$ qshell sandbox template delete tmpl-aaa tmpl-bbb -y
```

4. 交互式选择删除
```
$ qshell sandbox template delete -s
$ qshell sbx tpl dl -s
```
