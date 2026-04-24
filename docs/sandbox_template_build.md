# 简介
`sandbox template build`（别名 `bd`）创建新模板并触发构建，或对已有模板重新构建。

**创建新模板** 支持三种构建模式：`--from-image` 基于 Docker 镜像、`--from-template` 基于已有模板、`--dockerfile` 基于 Dockerfile（v2 构建系统）。

**重新构建已有模板**（`--template-id`）必须提供 `--dockerfile`——服务端 rebuild 接口要求在请求体中携带 Dockerfile 内容。

支持 `--no-cache` 强制完整构建和 `--wait` 流式查看构建日志。

# 格式
```
qshell sandbox template build [--name <name>] [--template-id <id>] [--from-image <image>] [--from-template <template>] [--dockerfile <path>] [--path <dir>] [--start-cmd <cmd>] [--ready-cmd <cmd>] [--cpu <N>] [--memory <N>] [--wait] [--no-cache]
qshell sbx tpl bd [--name <name>] [--template-id <id>] [--from-image <image>] [--from-template <template>] [--dockerfile <path>] [--path <dir>] [--start-cmd <cmd>] [--ready-cmd <cmd>] [--cpu <N>] [--memory <N>] [--wait] [--no-cache]
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
- `--template-id`：已有模板 ID（重新构建时使用，与 --name 二选一，必须同时提供 `--dockerfile`）
- `--from-image`：基础 Docker 镜像
- `--from-template`：基础模板
- `--dockerfile`：Dockerfile 路径（启用 v2 构建系统，自动解析 FROM/RUN/COPY 等指令）
- `--path`：构建上下文目录（默认为 Dockerfile 所在目录，与 --dockerfile 配合使用）
- `--start-cmd`：构建完成后执行的启动命令（Dockerfile 模式下可从 CMD/ENTRYPOINT 自动提取）
- `--ready-cmd`：就绪检查命令（Dockerfile 模式下默认为 "sleep 20"）
- `--cpu`：沙箱 CPU 核数
- `--memory`：沙箱内存大小（MiB）
- `--wait`：等待构建完成，实时流式显示构建日志（带彩色级别标签）
- `--no-cache`：强制完整构建，忽略缓存

# 示例
1. 从 Docker 镜像创建并构建模板
```
$ qshell sandbox template build --name my-template --from-image ubuntu:22.04 --wait
$ qshell sbx tpl bd --name my-template --from-image ubuntu:22.04 --wait
```

2. 从 Dockerfile 构建模板
```
$ qshell sandbox template build --name my-template --dockerfile ./Dockerfile --wait
$ qshell sbx tpl bd --name my-template --dockerfile ./Dockerfile --wait
```

3. 从 Dockerfile 构建，指定构建上下文目录
```
$ qshell sandbox template build --name my-template --dockerfile ./docker/Dockerfile --path ./src --wait
```

4. 重新构建已有模板（rebuild 必须提供 Dockerfile）
```
$ qshell sandbox template build --template-id tmpl-xxxxxxxxxxxx --dockerfile ./Dockerfile --wait
```

5. 使用 Dockerfile 重新构建已有模板（忽略缓存）
```
$ qshell sandbox template build --template-id tmpl-xxxxxxxxxxxx --dockerfile ./Dockerfile --no-cache --wait
```

6. 强制完整构建（忽略缓存）
```
$ qshell sandbox template build --template-id tmpl-xxxxxxxxxxxx --dockerfile ./Dockerfile --no-cache --wait
```

7. 指定启动命令和资源配置
```
$ qshell sandbox template build --name my-app --from-image node:18 --start-cmd "npm start" --cpu 2 --memory 1024 --wait
```
