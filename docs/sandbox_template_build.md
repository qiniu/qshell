# 简介
`sandbox template build`（别名 `bd`）分两步创建或更新模板：先申请构建并提交资源设置，再以选定的底模板启动构建。

每次构建分两步：

1. 按 `name` 申请构建（v3 创建接口）。同名模板已存在时沿用该模板 ID，并把 `cpu_count`、`memory_mb`、`min_free_disk_mb` 等资源设置应用到该模板；不存在时创建新模板。
2. 用选定的底模板启动构建（v2 构建接口）。

底模板有三种来源：`--from-image`（Docker 镜像）、`--from-template`（已有模板）、`--dockerfile`（Dockerfile 中的 `FROM`）。`--template-id` 与 `--from-template` 等价，都表示构建所用的底模板。

`--dockerfile` 可以和 `--from-template` 或 `--from-image` 组合使用，用选定的基础环境作为 base，再应用 `RUN` / `COPY` / `ENV` 等增量步骤。和 `--from-template` 组合时，Dockerfile 里的 `FROM` 会被解析但不会作为真实基础镜像。

支持 `qshell.sandbox.toml` 配置文件、`--config` 显式指定配置文件、`--no-cache` 强制完整构建和 `--wait` 流式查看构建日志。

# 格式
```
qshell sandbox template build [--name <name>] [--template-id <id>] [--from-image <image>] [--from-template <template>] [--dockerfile <path>] [--path <dir>] [--start-cmd <cmd>] [--ready-cmd <cmd>] [--cpu <N>] [--memory <N>] [--min-free-disk <N>] [--wait] [--no-cache] [--config <path>]
qshell sbx tpl bd [--name <name>] [--template-id <id>] [--from-image <image>] [--from-template <template>] [--dockerfile <path>] [--path <dir>] [--start-cmd <cmd>] [--ready-cmd <cmd>] [--cpu <N>] [--memory <N>] [--min-free-disk <N>] [--wait] [--no-cache] [--config <path>]
```

# 帮助文档
```
$ qshell sandbox template build -h
$ qshell sandbox template build --doc
```

# 鉴权
需要配置 `QINIU_API_KEY` 或 `E2B_API_KEY` 环境变量。

# 参数
- `--name`：模板名称，必填。按 name 定位模板；同名模板已存在时沿用其模板 ID
- `--template-id`：构建所用的底模板 ID，与 `--from-template` 等价，二者只能选一个
- `--from-image`：底 Docker 镜像
- `--from-template`：底模板名称或 ID
- `--dockerfile`：Dockerfile 路径（启用 v2 构建系统，自动解析 FROM/RUN/COPY 等指令；可与 --from-template 或 --from-image 组合叠加增量步骤）
- `--path`：构建上下文目录（默认为 Dockerfile 所在目录，与 --dockerfile 配合使用）
- `--start-cmd`：构建完成后执行的启动命令（Dockerfile 模式下可从 CMD/ENTRYPOINT 自动提取）
- `--ready-cmd`：就绪检查命令（Dockerfile 模式下默认为 "sleep 20"）
- `--cpu`：沙箱 CPU 核数
- `--memory`：沙箱内存大小（MiB）
- `--min-free-disk`：构建后的最小剩余空间（MiB），显式 0 表示不额外扩容；省略时使用团队默认值
- `--wait`：等待构建完成，实时流式显示构建日志（带彩色级别标签）
- `--no-cache`：强制完整构建，忽略缓存
- `--config`：显式指定 `qshell.sandbox.toml` 配置文件路径；未指定时，如果当前目录存在 `qshell.sandbox.toml`，命令会自动读取

# 配置文件
`sandbox template build` 会按 `CLI flag > 配置文件 > 内置默认值` 合并参数。配置文件可提供 `template_id`、`name`、`dockerfile`、`path`、`from_image`、`from_template`、`start_cmd`、`ready_cmd`、`cpu_count`、`memory_mb`、`min_free_disk_mb` 和 `no_cache`。

其中 `template_id` 与 `from_template` 都表示底模板，二者只能保留一个。`name` 表示要构建的模板名称。

每次构建都会按 `name` 申请构建：同名模板已存在时沿用其模板 ID 并应用资源设置，不存在时创建。因此同一份配置可以在多个环境复用，无需维护模板 ID。

若 `name` 命中其他团队的公共模板，命令会直接报错，因为后续按 name 解析仍会指向那个公共模板。请使用本团队独有的名称。

更多字段说明见 `docs/sandbox_template_config.md`。

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

5. 构建后至少保留 5 GiB 剩余空间
```
$ qshell sandbox template build --name my-template --from-image ubuntu:22.04 --min-free-disk 5120 --wait
```

6. 以模板 ID 作为底模板，构建到新名字
```
$ qshell sandbox template build --name my-template-v2 --template-id tmpl-xxxxxxxxxxxx --min-free-disk 20480 --wait
```

`--template-id` 与 `--from-template` 等价，二者只能选一个。

7. 更新同名模板的资源配置（忽略缓存）
```
$ qshell sandbox template build --name my-template --from-template base --dockerfile ./Dockerfile --no-cache --wait
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

# 构建后的最小剩余空间
`--min-free-disk <MiB>` 表示构建后的最小剩余空间，不是磁盘总大小。省略时使用团队默认值；显式指定 `--min-free-disk=0` 表示不为剩余空间扩容，已有文件系统不会缩小。
```shell
qshell sandbox template build --name demo --from-image ubuntu:22.04 --min-free-disk 5120 --wait
```
也可在 qshell.sandbox.toml 中配置 `min_free_disk_mb = 5120`。CLI 显式指定的值（包括 0）优先于配置文件。
