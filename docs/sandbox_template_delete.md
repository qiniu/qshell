# 简介
`sandbox template delete`（别名 `dl`）删除一个或多个沙箱模板。

# 格式
```
qshell sandbox template delete [templateIDs...] [-y]
qshell sbx tpl dl [templateIDs...] [-y]
```

# 帮助文档
```
$ qshell sandbox template delete -h
$ qshell sandbox template delete --doc
```

# 鉴权
需要配置 `QINIU_API_KEY` 或 `E2B_API_KEY` 环境变量。

# 参数
- `templateIDs`：一个或多个模板 ID。未传入时，自动读取当前目录 `qshell.sandbox.toml`；优先使用 `template_id`，否则按 `name` 查找远端模板
- `-y, --yes`：跳过确认提示

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

4. 删除当前目录配置文件对应的模板（`template_id` 或 `name`）
```
$ qshell sandbox template delete -y
$ qshell sbx tpl dl -y
```

# 非交互式调用（CI / AI Agent / 管道）

当 stdin 不是终端时，缺省的确认提示会立即报错并退出。自动化场景必须传 `-y` / `--yes`。
