# 简介
`sandbox template init`（别名 `it`）初始化一个新的模板项目，生成指定语言的脚手架文件。

# 格式
```
qshell sandbox template init [--name <name>] [--language <lang>] [--path <path>]
qshell sbx tpl it [--name <name>] [--language <lang>] [--path <path>]
```

# 帮助文档
```
$ qshell sandbox template init -h
$ qshell sandbox template init --doc
```

# 鉴权
此命令不需要 API Key，仅在本地生成文件。

# 参数
- `--name`：模板项目名称（必须匹配 `[a-z0-9][a-z0-9_-]*`）
- `--language`：编程语言（支持：go, typescript, python）
- `--path`：输出目录（默认为 `./<name>`）

未提供参数时将使用交互式提示。

# 支持的语言和生成文件
| 语言 | 生成文件 |
|------|---------|
| Go | `main.go`, `go.mod`, `Makefile` |
| TypeScript | `template.ts`, `package.json` |
| Python | `template.py`, `requirements.txt` |

# 示例
1. 交互式初始化
```
$ qshell sandbox template init
$ qshell sbx tpl it
```

2. 初始化 Go 项目
```
$ qshell sandbox template init --name my-template --language go
```

3. 初始化 TypeScript 项目到指定目录
```
$ qshell sandbox template init --name my-api --language typescript --path ./my-api
```

4. 初始化 Python 项目
```
$ qshell sandbox template init --name my-script --language python
```
