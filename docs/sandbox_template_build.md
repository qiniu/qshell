# 简介
`sandbox template build`（别名 `bd`）创建新模板并触发构建，或对已有模板重新构建。

**创建新模板** 支持三种基础构建模式：`--from-image` 基于 Docker 镜像、`--from-template` 基于已有模板、`--dockerfile` 基于 Dockerfile（v2 构建系统）。

`--dockerfile` 也可以和 `--from-template` 或 `--from-image` 组合使用，用选定的基础环境作为 base，再应用 Dockerfile 中的 `RUN` / `COPY` / `ENV` 等增量步骤。和 `--from-template` 组合时，Dockerfile 里的 `FROM` 会被解析但不会作为真实基础镜像。

**重新构建已有模板**（`--template-id`）必须提供 `--dockerfile`——服务端 rebuild 接口要求在请求体中携带 Dockerfile 内容。

支持 `qshell.sandbox.toml` 配置文件、`--config` 显式指定配置文件、`--no-cache` 强制完整构建和 `--wait` 流式查看构建日志。

# 格式
```
qshell sandbox template build [--name <name>] [--template-id <id>] [--from-image <image>] [--from-template <template>] [--dockerfile <path>] [--path <dir>] [--start-cmd <cmd>] [--ready-cmd <cmd>] [--cpu <N>] [--memory <N>] [--wait] [--no-cache] [--config <path>]
qshell sbx tpl bd [--name <name>] [--template-id <id>] [--from-image <image>] [--from-template <template>] [--dockerfile <path>] [--path <dir>] [--start-cmd <cmd>] [--ready-cmd <cmd>] [--cpu <N>] [--memory <N>] [--wait] [--no-cache] [--config <path>]
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
- `--from-image`：基础 Docker 镜像，仅创建新模板时可用
- `--from-template`：基础模板，仅创建新模板时可用
- `--dockerfile`：Dockerfile 路径（启用 v2 构建系统，自动解析 FROM/RUN/COPY 等指令；可与 --from-template 或 --from-image 组合叠加增量步骤）
- `--path`：构建上下文目录（默认为 Dockerfile 所在目录，与 --dockerfile 配合使用）
- `--start-cmd`：构建完成后执行的启动命令（Dockerfile 模式下可从 CMD/ENTRYPOINT 自动提取）
- `--ready-cmd`：就绪检查命令（Dockerfile 模式下默认为 "sleep 20"）
- `--cpu`：沙箱 CPU 核数
- `--memory`：沙箱内存大小（MiB）
- `--wait`：等待构建完成，实时流式显示构建日志（带彩色级别标签）
- `--no-cache`：强制完整构建，忽略缓存
- `--config`：显式指定 `qshell.sandbox.toml` 配置文件路径；未指定时，如果当前目录存在 `qshell.sandbox.toml`，命令会自动读取

# 配置文件
`sandbox template build` 会按 `CLI flag > 配置文件 > 内置默认值` 合并参数。配置文件可提供 `template_id`、`name`、`dockerfile`、`path`、`from_image`、`from_template`、`start_cmd`、`ready_cmd`、`cpu_count`、`memory_mb` 和 `no_cache`。

首次构建成功且配置文件存在、`template_id` 为空时，命令会自动把新模板 ID 回写到配置文件。后续再次运行同一配置文件会进入 rebuild 流程；此时配置文件中保留的 `from_image` / `from_template` 会被忽略，rebuild 只使用 `template_id`、`dockerfile`、资源规格和启动/就绪命令等参数。更多字段说明见 `docs/sandbox_template_config.md`。

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

4. 基于已有模板构建，并叠加 Dockerfile 步骤
```
$ qshell sandbox template build --name claude --from-template agents-base --dockerfile ./Dockerfile --wait
```

这种模式下实际 base 来自 `agents-base`，Dockerfile 里的 `FROM` 仅用于兼容 Dockerfile 解析，不会触发基础镜像拉取。

5. 重新构建已有模板（rebuild 必须提供 Dockerfile）
```
$ qshell sandbox template build --template-id tmpl-xxxxxxxxxxxx --dockerfile ./Dockerfile --wait
```

6. 使用 Dockerfile 重新构建已有模板（忽略缓存）
```
$ qshell sandbox template build --template-id tmpl-xxxxxxxxxxxx --dockerfile ./Dockerfile --no-cache --wait
```

7. 强制完整构建（忽略缓存）
```
$ qshell sandbox template build --template-id tmpl-xxxxxxxxxxxx --dockerfile ./Dockerfile --no-cache --wait
```

8. 使用当前目录的 `qshell.sandbox.toml`
```
$ qshell sandbox template build --wait
$ qshell sbx tpl bd --wait
```

9. 显式指定配置文件
```
$ qshell sandbox template build --config ./configs/prod.toml --wait
$ qshell sbx tpl bd --config ./configs/prod.toml --wait
```

10. 指定启动命令和资源配置
```
$ qshell sandbox template build --name my-app --from-image node:18 --start-cmd "npm start" --cpu 2 --memory 1024 --wait
```
