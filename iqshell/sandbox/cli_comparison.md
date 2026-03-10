# qshell sandbox vs E2B CLI 功能对比

本文档对比 qshell sandbox 与 [E2B CLI](https://github.com/e2b-dev/E2B) 的功能差异，帮助了解当前实现状态和后续改进方向。

## 功能对比总览

| 功能 | qshell sandbox | E2B CLI | 状态 |
|------|---------------|---------|------|
| **认证 - 登录/登出** | 无，完全依赖环境变量 | `auth login/logout` 浏览器 OAuth 流程 | 缺失 |
| **认证 - 用户信息** | 无 | `auth info` 显示邮箱/团队 | 缺失 |
| **认证 - 团队切换** | 无 | `auth configure` 交互式选择团队 | 缺失 |
| **认证 - 凭证持久化** | 无，每次设环境变量 | `~/.e2b/config.json` 自动存储 | 缺失 |
| **认证 - .env 文件** | 支持当前目录 `.env` 文件加载 | 无 | qshell 优势 |
| **沙箱 - 创建** | `sbx cr`（支持 `--detach/-d`） | `sbx cr`（支持 `--detach/-d`） | 一致 |
| **沙箱 - 连接终端** | `sbx cn` | `sbx cn` | 一致 |
| **沙箱 - 列表** | `sbx ls` | `sbx ls` | 一致 |
| **沙箱 - 杀死** | `sbx kl`（并发 goroutine） | `sbx kl` | 一致 |
| **沙箱 - 暂停** | `sbx ps`（支持 `--all`、并发） | `sbx ps` (beta) | 一致 |
| **沙箱 - 恢复** | `sbx rs`（支持 `--all`、并发） | `sbx rs` | 一致 |
| **沙箱 - 远程执行命令** | `sbx ex` 支持后台/环境变量/信号转发/退出码传递 | `sbx ex` 支持 stdin 管道/后台/环境变量/信号转发 | 一致 |
| **沙箱 - 日志** | `sbx lg` | `sbx lg` | 一致 |
| **沙箱 - 指标** | `sbx mt` | `sbx mt` | 一致 |
| **模板 - 列表** | `sbx tpl ls` | `tpl ls` | 一致 |
| **模板 - 详情** | `sbx tpl gt` | 无 | qshell 优势 |
| **模板 - 删除** | `sbx tpl dl` | `tpl dl` | 一致 |
| **模板 - 发布/取消** | `sbx tpl pb/upb` | `tpl pb/upb` | 一致 |
| **模板 - 构建模式** | 三种：`--from-image` / `--from-template` / `--dockerfile` | 仅 Dockerfile（v2 远程构建） | qshell 优势 |
| **模板 - 构建状态查询** | `sbx tpl bds` 独立命令 | 无，构建时内联流式输出 | qshell 优势 |
| **模板 - 初始化脚手架** | Go / TypeScript / Python | TypeScript / Python sync / Python async | 各有侧重 |
| **模板 - 项目配置文件** | 无，所有参数通过 CLI 标志 | `e2b.toml` 持久化模板配置 | 缺失 |
| **模板 - v1→v2 迁移** | 无（无历史包袱） | `tpl migrate` | 不适用 |
| **Dockerfile 解析** | 自研轻量解析器，无 buildkit 依赖 | 使用 SDK 的 `Template.fromDockerfile()` | qshell 优势 |
| **COPY 步骤上传** | SHA-256 内容寻址缓存 + `io.Pipe` 流式上传 | SDK 内部处理 | qshell 优势 |
| **终端 - PTY raw mode** | 支持 | 支持 | 一致 |
| **终端 - 批量输入** | 10ms batchedWriter | 10ms BatchedQueue | 一致 |
| **终端 - 窗口大小** | SIGWINCH 处理 | resize 处理 | 一致 |
| **终端 - Keep-alive** | 每 5 秒心跳，30 秒超时 | SDK 内部处理 | 一致 |
| **输出 - 彩色** | `fatih/color` | `chalk` 自定义调色板 | 一致 |
| **输出 - 表格** | `text/tabwriter` | `console-table-printer` | 一致 |
| **输出 - 框形消息** | `lipgloss` 圆角边框 | `boxen` 错误/警告/成功框 | 一致 |
| **输出 - 语法高亮** | `lipgloss` 代码块样式 | `cli-highlight` 代码示例 | 一致 |
| **输出 - 可点击链接** | `termenv` ANSI 超链接 | ANSI 超链接（Dashboard URL） | 一致 |
| **输出 - SDK 使用示例** | 构建成功后显示 Go/Python/TS 代码 | 构建成功后显示 Python/TS 代码 | 一致 |
| **输出 - JSON 格式** | `--format json` | `--format json` | 一致 |
| **嵌入式文档** | `--doc` 显示 Markdown 帮助 | 无 | qshell 优势 |
| **更新检查** | 无 | `update-notifier` 每 8 小时检查 | 缺失 |
| **交互式 TUI** | `charmbracelet/huh` | `inquirer` / `@inquirer/prompts` | 一致 |
| **运行时依赖** | 无（单一 Go 二进制） | 需要 Node.js | qshell 优势 |
| **E2B 环境变量兼容** | 支持 `E2B_API_KEY` / `E2B_API_URL` 回退 | 原生 | 一致 |

---

## 子命令与参数详细对比

### 认证命令

E2B 提供完整的 `auth` 命令组，qshell 无对应命令，完全依赖环境变量。

#### `auth login`

| 维度 | qshell | E2B |
|------|--------|-----|
| **命令** | 无 | `e2b auth login` |
| **说明** | - | 浏览器 OAuth 登录，本地起 HTTP 服务接收回调 |
| **参数** | - | 无 |
| **标志** | - | 无 |

#### `auth logout`

| 维度 | qshell | E2B |
|------|--------|-----|
| **命令** | 无 | `e2b auth logout` |
| **说明** | - | 删除 `~/.e2b/config.json` |
| **参数** | - | 无 |
| **标志** | - | 无 |

#### `auth info`

| 维度 | qshell | E2B |
|------|--------|-----|
| **命令** | 无 | `e2b auth info` |
| **说明** | - | 显示当前用户邮箱、团队名、团队 ID |
| **参数** | - | 无 |
| **标志** | - | 无 |

#### `auth configure`

| 维度 | qshell | E2B |
|------|--------|-----|
| **命令** | 无 | `e2b auth configure` |
| **说明** | - | 交互式选择用户所属团队 |
| **参数** | - | 无 |
| **标志** | - | 无 |

---

### 沙箱命令

#### `sandbox list`

| 维度 | qshell (`sbx ls`) | E2B (`sbx ls`) |
|------|-------------------|----------------|
| **别名** | `ls` | `ls` |
| **说明** | List sandboxes | list all sandboxes, by default it list only running ones |
| **参数** | 无 | 无 |

**标志对比：**

| 标志 | qshell | E2B | 差异 |
|------|--------|-----|------|
| `--state` / `-s` | `string`，默认运行时为 `running` | `string`，默认 `running` | 一致 |
| `--metadata` / `-m` | `string`，格式 `key1=value1,key2=value2` | `string`，格式 `key1=value1` | 一致 |
| `--limit` / `-l` | `int32`，默认 `0` | `int`，默认 `1000` | 默认值不同 |
| `--format` / `-f` | `string`，默认 `pretty` | `string`，默认 `pretty` | 一致 |

#### `sandbox create`

| 维度 | qshell (`sbx cr`) | E2B (`sbx cr`) |
|------|-------------------|----------------|
| **别名** | `cr` | `cr` |
| **说明** | Create a sandbox and connect to its terminal | create sandbox and connect terminal to it |
| **参数** | `[template]`（必填） | `[template]`（可选，可从 `e2b.toml` 读取） |

**标志对比：**

| 标志 | qshell | E2B | 差异 |
|------|--------|-----|------|
| `--timeout` / `-t` | `int32`，沙箱超时秒数 | 无 | E2B 缺失 |
| `--metadata` / `-m` | `string`，元数据键值对 | 无 | E2B 缺失 |
| `--detach` / `-d` | `bool`，创建但不连接终端 | `bool`，创建但不连接终端 | 一致 |
| `--env-var` / `-e` | `stringArray`，环境变量 `KEY=VALUE`（可重复） | 无 | E2B 缺失 |
| `--auto-pause` | `bool`，超时后自动暂停（而非杀死） | 无 | E2B 缺失 |
| `--path` / `-p` | 无 | `string`，根目录路径 | qshell 缺失 |
| `--config` | 无 | `string`，`e2b.toml` 路径 | qshell 缺失 |

#### `sandbox connect`

| 维度 | qshell (`sbx cn`) | E2B (`sbx cn`) |
|------|-------------------|----------------|
| **别名** | `cn` | `cn` |
| **说明** | Connect to an existing sandbox terminal | connect terminal to already running sandbox |
| **参数** | `<sandboxID>`（必填） | `<sandboxID>`（必填） |
| **标志** | 无 | 无 |
| **差异** | 完全一致 | |

#### `sandbox kill`

| 维度 | qshell (`sbx kl`) | E2B (`sbx kl`) |
|------|-------------------|----------------|
| **别名** | `kl` | `kl` |
| **说明** | Kill one or more sandboxes | kill sandbox |
| **参数** | `[sandboxIDs...]`（可变参数） | `[sandboxIDs...]`（可变参数） |

**标志对比：**

| 标志 | qshell | E2B | 差异 |
|------|--------|-----|------|
| `--all` / `-a` | `bool` | `bool` | 一致 |
| `--state` / `-s` | `string`，配合 `--all` 使用 | `string`，配合 `--all` 使用 | 一致 |
| `--metadata` / `-m` | `string`，配合 `--all` 使用 | `string`，配合 `--all` 使用 | 一致 |

#### `sandbox pause`

| 维度 | qshell (`sbx ps`) | E2B (`sbx ps`) |
|------|-------------------|----------------|
| **别名** | `ps` | `ps` |
| **说明** | Pause one or more sandboxes | pause sandbox (beta) |
| **参数** | `[sandboxIDs...]`（可变参数） | `<sandboxID>`（必填） |

**标志对比：**

| 标志 | qshell | E2B | 差异 |
|------|--------|-----|------|
| `--all` / `-a` | `bool`，暂停所有沙箱 | 无 | E2B 缺失 |
| `--state` / `-s` | `string`，配合 `--all` 筛选状态 | 无 | E2B 缺失 |
| `--metadata` / `-m` | `string`，配合 `--all` 筛选元数据 | 无 | E2B 缺失 |

#### `sandbox resume`

| 维度 | qshell (`sbx rs`) | E2B (`sbx rs`) |
|------|-------------------|----------------|
| **别名** | `rs` | `rs` |
| **说明** | Resume one or more paused sandboxes | resume paused sandbox |
| **参数** | `[sandboxIDs...]`（可变参数） | `<sandboxID>`（必填） |

**标志对比：**

| 标志 | qshell | E2B | 差异 |
|------|--------|-----|------|
| `--all` / `-a` | `bool`，恢复所有暂停的沙箱 | 无 | E2B 缺失 |
| `--metadata` / `-m` | `string`，配合 `--all` 筛选元数据 | 无 | E2B 缺失 |

#### `sandbox exec`

| 维度 | qshell (`sbx ex`) | E2B (`sbx ex`) |
|------|-------------------|----------------|
| **别名** | `ex` | `ex` |
| **说明** | Execute a command in a sandbox | execute a command in a running sandbox |
| **参数** | `<sandboxID> -- <command...>`（`--` 分隔） | `<sandboxID> <command...>`（空格分隔） |

**标志对比：**

| 标志 | qshell | E2B | 差异 |
|------|--------|-----|------|
| `--background` / `-b` | `bool`，后台运行 | `bool`，后台运行 | 一致 |
| `--cwd` / `-c` | `string`，工作目录 | `string`，工作目录 | 一致 |
| `--user` / `-u` | `string`，执行用户 | `string`，执行用户 | 一致 |
| `--env` / `-e` | `stringArray`，可重复 | `string`，可重复 | 一致 |

#### `sandbox logs`

| 维度 | qshell (`sbx lg`) | E2B (`sbx lg`) |
|------|-------------------|----------------|
| **别名** | `lg` | `lg` |
| **说明** | View sandbox logs | show logs for sandbox |
| **参数** | `<sandboxID>`（必填） | `<sandboxID>`（必填） |

**标志对比：**

| 标志 | qshell | E2B | 差异 |
|------|--------|-----|------|
| `--level` | `string`，默认 `INFO`，可选 `DEBUG/INFO/WARN/ERROR` | `string`，默认 `INFO`，可选 `DEBUG/INFO/WARN/ERROR` | 一致 |
| `--limit` | `int32`，默认 `0` | 无 | E2B 缺失 |
| `--format` | `string`，默认 `pretty` | `string`，默认 `pretty`，可选 `json/pretty` | 一致 |
| `--follow` / `-f` | `bool` | `bool` | 一致 |
| `--loggers` | `string`，逗号分隔 | `string`（可选值），逗号分隔 | 一致 |

#### `sandbox metrics`

| 维度 | qshell (`sbx mt`) | E2B (`sbx mt`) |
|------|-------------------|----------------|
| **别名** | `mt` | `mt` |
| **说明** | View sandbox resource metrics | show metrics for sandbox |
| **参数** | `<sandboxID>`（必填） | `<sandboxID>`（必填） |

**标志对比：**

| 标志 | qshell | E2B | 差异 |
|------|--------|-----|------|
| `--format` | `string`，默认 `pretty` | `string`，默认 `pretty` | 一致 |
| `--follow` / `-f` | `bool` | `bool` | 一致 |

---

### 模板命令

#### `template list`

| 维度 | qshell (`sbx tpl ls`) | E2B (`tpl ls`) |
|------|----------------------|----------------|
| **别名** | `ls` | `ls` |
| **说明** | List sandbox templates | list sandbox templates |
| **参数** | 无 | 无 |

**标志对比：**

| 标志 | qshell | E2B | 差异 |
|------|--------|-----|------|
| `--format` | `string`，默认 `pretty` | `string`，默认 `pretty` | 一致 |
| `--team` / `-t` | 无（服务端不支持） | `string`，指定团队 ID | 不适用 |

#### `template get`

| 维度 | qshell (`sbx tpl gt`) | E2B |
|------|----------------------|-----|
| **命令** | `sandbox template get` | 无 |
| **别名** | `gt` | - |
| **说明** | Get template details | - |
| **参数** | `<templateID>`（必填） | - |
| **标志** | 无 | - |

#### `template delete`

| 维度 | qshell (`sbx tpl dl`) | E2B (`tpl dl`) |
|------|----------------------|----------------|
| **别名** | `dl` | `dl` |
| **说明** | Delete one or more templates | delete sandbox template and `e2b.toml` config |
| **参数** | `[templateIDs...]`（可变参数） | `[template]`（可选，单个） |

**标志对比：**

| 标志 | qshell | E2B | 差异 |
|------|--------|-----|------|
| `--yes` / `-y` | `bool`，跳过确认 | `bool`，跳过确认 | 一致 |
| `--select` / `-s` | `bool`，交互式选择 | `bool`，交互式选择 | 一致 |
| `--path` / `-p` | 无 | `string`，根目录路径 | qshell 缺失 |
| `--config` | 无 | `string`，`e2b.toml` 路径 | qshell 缺失 |
| `--team` / `-t` | 无 | `string`，指定团队 ID | qshell 缺失 |

#### `template build`

| 维度 | qshell (`sbx tpl bd`) | E2B (`tpl ct` / `tpl bd`) |
|------|----------------------|---------------------------|
| **别名** | `bd` | `ct`（v2 当前）/ `bd`（v1 已废弃） |
| **说明** | Build a template，支持三种构建模式 | build Dockerfile as a Sandbox template |
| **参数** | 无（通过标志指定） | `<template-name>`（v2 必填）/ `[template]`（v1 可选） |

**标志对比：**

| 标志 | qshell | E2B (v2 create) | E2B (v1 build) | 差异 |
|------|--------|-----------------|----------------|------|
| `--name` | `string`，模板名称 | 无（位置参数） | `--name` / `-n` | qshell 用标志，E2B v2 用位置参数 |
| `--template-id` | `string`，重建已有模板 | 无 | `[template]` 位置参数 | qshell 用标志 |
| `--from-image` | `string`，基础 Docker 镜像 | 无 | 无 | qshell 独有 |
| `--from-template` | `string`，基于已有模板 | 无 | 无 | qshell 独有 |
| `--dockerfile` | `string`，Dockerfile 路径 | `--dockerfile` / `-d` | `--dockerfile` / `-d` | 一致 |
| `--path` | `string`，构建上下文目录 | `--path` / `-p` | `--path` / `-p` | 一致 |
| `--start-cmd` | `string`，启动命令 | `--cmd` / `-c` | `--cmd` / `-c` | 标志名不同 |
| `--ready-cmd` | `string`，就绪检查命令 | `--ready-cmd` | `--ready-cmd` | 一致 |
| `--cpu` | `int32`，CPU 核数 | `--cpu-count` | `--cpu-count` | 标志名不同 |
| `--memory` | `int32`，内存 MiB | `--memory-mb` | `--memory-mb` | 标志名不同 |
| `--no-cache` | `bool`，忽略缓存 | `--no-cache` | `--no-cache` | 一致 |
| `--wait` | `bool`，等待构建完成 | 无（默认等待） | 无（默认轮询等待） | qshell 需显式指定 |
| `--build-arg` | 无 | 无 | `--build-arg`（可变参数） | qshell 缺失（v1 独有） |
| `--config` | 无 | 无 | `--config` | qshell 缺失 |
| `--team` / `-t` | 无 | 无 | `--team` / `-t` | qshell 缺失 |

#### `template builds`

| 维度 | qshell (`sbx tpl bds`) | E2B |
|------|------------------------|-----|
| **命令** | `sandbox template builds` | 无 |
| **别名** | `bds` | - |
| **说明** | View template build status | - |
| **参数** | `<templateID>`（必填）、`<buildID>`（必填） | - |
| **标志** | 无 | - |

#### `template publish`

| 维度 | qshell (`sbx tpl pb`) | E2B (`tpl pb`) |
|------|----------------------|----------------|
| **别名** | `pb` | `pb` |
| **说明** | Publish templates (make public) | publish sandbox template |
| **参数** | `[templateIDs...]`（可变参数） | `[template]`（可选，单个） |

**标志对比：**

| 标志 | qshell | E2B | 差异 |
|------|--------|-----|------|
| `--yes` / `-y` | `bool`，跳过确认 | `bool`，跳过确认 | 一致 |
| `--select` / `-s` | `bool`，交互式选择 | `bool`，交互式选择 | 一致 |
| `--path` / `-p` | 无 | `string`，根目录路径 | qshell 缺失 |
| `--config` | 无 | `string`，`e2b.toml` 路径 | qshell 缺失 |
| `--team` / `-t` | 无 | `string`，指定团队 ID | qshell 缺失 |

#### `template unpublish`

| 维度 | qshell (`sbx tpl upb`) | E2B (`tpl upb`) |
|------|------------------------|-----------------|
| **别名** | `upb` | `upb` |
| **说明** | Unpublish templates (make private) | unpublish sandbox template |
| **参数** | `[templateIDs...]`（可变参数） | `[template]`（可选，单个） |

**标志对比：**

| 标志 | qshell | E2B | 差异 |
|------|--------|-----|------|
| `--yes` / `-y` | `bool`，跳过确认 | `bool`，跳过确认 | 一致 |
| `--select` / `-s` | `bool`，交互式选择 | `bool`，交互式选择 | 一致 |
| `--path` / `-p` | 无 | `string`，根目录路径 | qshell 缺失 |
| `--config` | 无 | `string`，`e2b.toml` 路径 | qshell 缺失 |
| `--team` / `-t` | 无 | `string`，指定团队 ID | qshell 缺失 |

#### `template init`

| 维度 | qshell (`sbx tpl it`) | E2B (`tpl it`) |
|------|----------------------|----------------|
| **别名** | `it` | `it` |
| **说明** | Initialize a new template project | initialize a new sandbox template using the SDK |
| **参数** | 无 | 无 |

**标志对比：**

| 标志 | qshell | E2B | 差异 |
|------|--------|-----|------|
| `--name` | `string`，项目名称 | `--name` / `-n`，模板名称 | 一致 |
| `--language` | `string`，可选 `go/typescript/python` | `--language` / `-l`，可选 `typescript/python-sync/python-async` | 语言选项不同 |
| `--path` | `string`，输出目录 | `--path` / `-p`，根目录路径 | 一致 |

#### `template migrate`

| 维度 | qshell | E2B (`tpl migrate`) |
|------|--------|---------------------|
| **命令** | 无 | `e2b template migrate` |
| **说明** | - | 将 `e2b.Dockerfile` + `e2b.toml` 迁移到 v2 SDK 格式 |
| **参数** | - | 无 |

**E2B 标志：**

| 标志 | 类型 | 说明 |
|------|------|------|
| `--dockerfile` / `-d` | `string` | Dockerfile 路径，默认 `e2b.Dockerfile` |
| `--config` | `string` | `e2b.toml` 路径 |
| `--language` / `-l` | `string` | 目标语言，可选 `typescript/python-sync/python-async` |
| `--path` / `-p` | `string` | 根目录路径 |

---

## qshell 优势

- **模板详情查看** (`tpl get`)：独立命令查看模板完整信息
- **构建状态查询** (`tpl builds`)：独立查看构建进度和日志
- **三种构建模式**：`--from-image`、`--from-template`、`--dockerfile`
- **内容寻址上传缓存**：COPY 步骤 SHA-256 去重，避免重复上传
- **Go 模板支持**：脚手架初始化支持 Go 语言
- **嵌入式文档**：`--doc` 标志直接显示 Markdown 帮助
- **零运行时依赖**：Go 编译为单一二进制，不需要 Node.js
- **自研 Dockerfile 解析器**：无 buildkit 依赖，轻量可控
- **流式上传**：`io.Pipe` 零内存拷贝
- **批量操作**：`delete/publish/unpublish/pause/resume/kill` 支持多个 ID 的可变参数 + `--all` 批量
- **.env 文件支持**：当前目录 `.env` 文件自动加载，按项目区分 API KEY
- **SDK 使用示例**：构建成功后额外显示 Go 语言示例

## qshell 待补齐功能

- **认证命令**：用户需手动设置环境变量，入门门槛高
- **项目配置文件**：无法像 `e2b.toml` 持久化模板配置（template_id、dockerfile、CPU/内存等）
- **团队管理**：服务端不支持 team 概念，`--team` 标志不适用
- **更新检查**：无自动检测新版本机制
