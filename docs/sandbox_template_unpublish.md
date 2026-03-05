# 简介
`sandbox template unpublish`（别名 `upb`）将模板设为私有（private），其他用户将无法使用。

# 格式
```
qshell sandbox template unpublish [templateIDs...] [-y] [-s]
qshell sbx tpl upb [templateIDs...] [-y] [-s]
```

# 帮助文档
```
$ qshell sandbox template unpublish -h
$ qshell sandbox template unpublish --doc
```

# 鉴权
需要配置 `QINIU_API_KEY` 或 `E2B_API_KEY` 环境变量。

# 参数
- `templateIDs`：一个或多个模板 ID（与 `--select` 二选一）
- `-y, --yes`：跳过确认提示
- `-s, --select`：交互式选择模板

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

3. 交互式选择
```
$ qshell sandbox template unpublish -s
```
