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
min_free_disk_mb = 5120

# 构建选项
no_cache = false
```

| 字段 | 类型 | 说明 |
|------|------|------|
| template_id | string | 构建所用的底模板 ID，与 from_template 等价；二者只能选一个 |
| name | string | 模板名称。build 按 name 定位模板；get/publish/delete/unpublish 未传 ID 时也按 name 查找 |
| dockerfile | string | Dockerfile 路径 |
| path | string | 构建上下文目录 |
| from_image | string | 基础 Docker 镜像 |
| from_template | string | 基础模板 |
| start_cmd | string | 容器启动命令 |
| ready_cmd | string | 就绪检查命令 |
| cpu_count | int | CPU 核数 |
| memory_mb | int | 内存（MiB） |
| min_free_disk_mb | int | 构建后的最小剩余空间（MiB），可设为 0；省略时使用团队默认值 |
| disk_size_mb | int | **已废弃**：模板 API 移除了磁盘大小参数，该字段不再生效，`build` 会在 stderr 给出迁移提示 |
| no_cache | bool | 强制忽略缓存 |

# 优先级

`CLI flag > 配置文件 > 内置默认值`

同时提供时，CLI 会覆盖配置文件，并在 stderr 打印一次覆盖提示。

# 构建输入组合

`from_image` 和 `from_template` 二选一，不能同时设置。`template_id` 与 `from_template` 等价，同样不能和它们同时使用。

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

1. `build` 按 `name` 定位模板，并申请构建；同名模板已存在时沿用其模板 ID
2. `template_id` / `from_template` / `from_image` 表示构建所用的底模板，不参与目标定位
3. `get` / `publish` / `delete` / `unpublish` 未传模板 ID 时，先从当前目录配置读取 `template_id`，否则按 `name` 查找
4. `build` 若发现 `name` 命中其他团队的公共模板，会直接报错，避免构建结果无法按 name 访问

# 团队协作约定

将 `qshell.sandbox.toml` 加入版本控制：
- 团队成员克隆仓库后直接 `qshell sandbox template build` 即可在当前环境按 `name` 定位或创建模板
- CI 脚本只需一行命令，无需维护散落的参数
- 通常只提交 `name` / `dockerfile` / 资源参数即可

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
min_free_disk_mb = 5120
```

运行：
```bash
qshell sandbox template build --wait
# 每次都会按 name 申请构建：同名模板已存在时沿用其模板 ID，并应用资源设置
```

再次构建：
```bash
qshell sandbox template build --wait  # 按 name 定位模板并应用当前配置
```

`min_free_disk_mb`：创建模板时构建后的最小剩余空间（MiB），可设为 0。省略时使用团队默认值；CLI 的 `--min-free-disk`（含显式 0）优先。

该设置在申请构建时提交，因此同名模板每次构建都会应用当前配置。省略时使用团队默认值；显式 0 表示不为剩余空间扩容，已有文件系统不会缩小。
