# 简介
`sandbox template publish`（别名 `pb`）将模板设为公开（public），允许其他用户使用。

# 格式
```
qshell sandbox template publish [templateIDs...] [-y] [-s]
qshell sbx tpl pb [templateIDs...] [-y] [-s]
```

# 帮助文档
```
$ qshell sandbox template publish -h
$ qshell sandbox template publish --doc
```

# 鉴权
需要配置 `QINIU_API_KEY` 或 `E2B_API_KEY` 环境变量。

# 参数
- `templateIDs`：一个或多个模板 ID。未传入且未使用 `--select` 时，自动读取当前目录 `qshell.sandbox.toml` 中的 `template_id`
- `-y, --yes`：跳过确认提示
- `-s, --select`：交互式选择模板

# 示例
1. 发布单个模板
```
$ qshell sandbox template publish tmpl-xxxxxxxxxxxx -y
$ qshell sbx tpl pb tmpl-xxxxxxxxxxxx -y
```

2. 发布多个模板
```
$ qshell sandbox template publish tmpl-aaa tmpl-bbb -y
```

3. 交互式选择发布
```
$ qshell sandbox template publish -s
```

4. 发布当前目录配置文件对应的模板
```
$ qshell sandbox template publish -y
$ qshell sbx tpl pb -y
```
