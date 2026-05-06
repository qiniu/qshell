# 简介
`sandbox template config`（别名 `cfg`）展示 `qshell.sandbox.toml` 模板配置文件说明。该命令只输出文档，不会读取、创建或修改本地配置文件。

# 格式
```
qshell sandbox template config
qshell sbx tpl cfg
```

# 帮助文档
```
$ qshell sandbox template config -h
$ qshell sandbox template config --doc
```

# qshell.sandbox.toml

`qshell sandbox template` 系列命令支持从项目根目录读取 `qshell.sandbox.toml`，
持久化模板参数。推荐使用环境无关的 `name` 定位模板；已有的 `template_id` 配置仍保持兼容。

# 完整字段说明

```toml
# 模板身份：推荐只写 name，让不同环境按 name 自动定位各自模板。
name = "claude"

# 兼容旧配置：已有 template_id 仍优先生效。
# template_id = "tmpl-xxxxxxxxxxxx"

# 构建输入
dockerfile    = "./Dockerfile"
path          = "."
from_image    = ""
from_template = ""

# 运行时
start_cmd = "/root/.jupyter/start-up.sh"
ready_cmd = ""

# 资源
cpu_count = 2
memory_mb = 2048

# 构建选项
no_cache = false
```

| 字段 | 类型 | 说明 |
|------|------|------|
| template_id | string | 模板 ID。存在时优先生效，进入 rebuild 流程 |
| name | string | 模板名称。未提供 template_id 时，build/get/publish/delete/unpublish 会按 name 查找远端模板 |
| dockerfile | string | Dockerfile 路径 |
| path | string | 构建上下文目录 |
| from_image | string | 基础 Docker 镜像 |
| from_template | string | 基础模板 |
| start_cmd | string | 容器启动命令 |
| ready_cmd | string | 就绪检查命令 |
| cpu_count | int | CPU 核数 |
| memory_mb | int | 内存（MiB） |
| no_cache | bool | 强制忽略缓存 |

# 优先级

`CLI flag > 配置文件 > 内置默认值`

同时提供时，CLI 会覆盖配置文件，并在 stderr 打印一次覆盖提示。

# 构建输入组合

`from_image` 和 `from_template` 二选一，不能同时设置。

`dockerfile` 可以单独使用，也可以和 `from_image` 或 `from_template` 组合使用：
- 只有 `dockerfile`：使用 Dockerfile 中的 `FROM` 作为基础镜像。
- `from_image + dockerfile`：使用 `from_image` 作为基础镜像，Dockerfile 中的 `FROM` 会被覆盖。
- `from_template + dockerfile`：使用 `from_template` 作为基础模板，Dockerfile 中的 `FROM` 会被解析但不会作为真实基础镜像。

`from_template + dockerfile` 适合在统一基础模板上叠加少量依赖，例如多个 agent 模板共用一个 `agents-base`。

# 配置文件查找规则

1. `--config <path>` 显式指定路径（仅 `build` 命令支持）
2. 当前工作目录下的 `qshell.sandbox.toml`
3. 未找到时按纯 CLI 模式运行（向后兼容）

# 模板定位规则

1. 显式传入的 `template_id` 或配置文件中的 `template_id` 优先生效
2. 未提供 `template_id` 时，使用 `name` 调用远端按 alias 点查
3. `build` 命中 name 时进入 rebuild；未命中时创建新模板
4. `get` / `publish` / `delete` / `unpublish` 未传模板 ID 时，也会先从当前目录配置读取 `template_id`，否则按 `name` 查找

# 自动回写行为

- 首次创建新模板（配置文件存在、`template_id` 为空）成功后，qshell 自动把新的 `template_id` 写入文件，兼容旧脚本
- 按 `name` 命中已有模板并进入 rebuild 时，不回写 `template_id`，配置文件可继续跨环境复用
- 未指定 `--wait` 时，在构建成功启动后回写；指定 `--wait` 时，在构建完成且状态为 `ready` 后回写
- 回写会替换已有的 `template_id` 赋值行（无论原值是否为空），或在文件头插入一行
- 注释、字段顺序、空白均保留
- 回写完成后 stdout 打印：`Written template_id to <path> (please commit this file)`
- 通过 `template_id` 或按 `name` 命中后再次执行 build 会进入 rebuild 流程，配置文件中保留的 `from_image` / `from_template` 会继续参与构建；如果从 CLI 显式传入这两个参数，命令会报错

# 团队协作约定

将 `qshell.sandbox.toml` 加入版本控制：
- 团队成员克隆仓库后直接 `qshell sandbox template build` 即可在当前环境按 `name` 定位或创建模板
- CI 脚本只需一行命令，无需维护散落的参数
- 通常只提交 `name` / `dockerfile` / 资源参数即可；如果选择保留 `template_id`，它会固定指向某一个环境里的模板

# 示例：最小项目

```
my-template/
├── Dockerfile
└── qshell.sandbox.toml
```

`qshell.sandbox.toml`:
```toml
name = "my-template"
dockerfile = "./Dockerfile"
cpu_count = 2
memory_mb = 2048
```

运行：
```bash
qshell sandbox template build --wait
# 首次运行时未找到同名模板则创建；同名模板已存在则 rebuild
```

下次重建：
```bash
qshell sandbox template build --wait  # 按 name 命中后 rebuild，不再 409
```
