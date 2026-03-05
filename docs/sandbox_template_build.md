# 简介
`sandbox template build`（别名 `bd`）创建新模板并触发构建，或对已有模板重新构建。

# 格式
```
qshell sandbox template build [--name <name>] [--template-id <id>] [--from-image <image>] [--from-template <template>] [--start-cmd <cmd>] [--ready-cmd <cmd>] [--cpu <N>] [--memory <N>] [--wait]
qshell sbx tpl bd [--name <name>] [--template-id <id>] [--from-image <image>] [--from-template <template>] [--start-cmd <cmd>] [--ready-cmd <cmd>] [--cpu <N>] [--memory <N>] [--wait]
```

# 帮助文档
```
$ qshell sandbox template build -h
$ qshell sandbox template build --doc
```

# 鉴权
需要配置 `QINIU_API_KEY` 或 `E2B_API_KEY` 环境变量。

# 参数
- `--name`：模板名称（创建新模板时使用，与 --template-id 二选一）
- `--template-id`：已有模板 ID（重新构建时使用，与 --name 二选一）
- `--from-image`：基础 Docker 镜像
- `--from-template`：基础模板
- `--start-cmd`：构建完成后执行的启动命令
- `--ready-cmd`：就绪检查命令
- `--cpu`：沙箱 CPU 核数
- `--memory`：沙箱内存大小（MiB）
- `--wait`：等待构建完成

# 示例
1. 从 Docker 镜像创建并构建模板
```
$ qshell sandbox template build --name my-template --from-image ubuntu:22.04 --wait
$ qshell sbx tpl bd --name my-template --from-image ubuntu:22.04 --wait
```

2. 重新构建已有模板
```
$ qshell sandbox template build --template-id tmpl-xxxxxxxxxxxx --from-image ubuntu:22.04
```

3. 指定启动命令和资源配置
```
$ qshell sandbox template build --name my-app --from-image node:18 --start-cmd "npm start" --cpu 2 --memory 1024 --wait
```
