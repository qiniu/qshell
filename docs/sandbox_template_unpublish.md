# 简介
`sandbox template unpublish`（别名 `upb`）将模板设为私有（private），其他用户将无法使用。

# 格式
```
qshell sandbox template unpublish [templateIDs...] [-y]
qshell sbx tpl upb [templateIDs...] [-y]
```

# 帮助文档
```
$ qshell sandbox template unpublish -h
$ qshell sandbox template unpublish --doc
```

# 鉴权
需要配置 `QINIU_API_KEY` 或 `E2B_API_KEY` 环境变量。

# 参数
- `templateIDs`：一个或多个模板 ID。未传入时，自动读取当前目录 `qshell.sandbox.toml`；优先使用 `template_id`，否则按 `name` 查找远端模板
- `-y, --yes`：跳过确认提示

# 示例
1. 取消发布单个模板
```
$ qshell sandbox template unpublish tmpl-xxxxxxxxxxxx -y
$ qshell sbx tpl upb tmpl-xxxxxxxxxxxx -y
```

2. 取消发布多个模板
```
$ qshell sandbox template unpublish tmpl-aaa tmpl-bbb -y
```

3. 取消发布当前目录配置文件对应的模板（`template_id` 或 `name`）
```
$ qshell sandbox template unpublish -y
$ qshell sbx tpl upb -y
```

# 非交互式调用（CI / AI Agent / 管道）

当 stdin 不是终端时，缺省的确认提示会立即报错并退出。自动化场景必须传 `-y` / `--yes`。
