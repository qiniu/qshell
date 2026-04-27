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
持久化模板参数并自动管理 `template_id`，方便团队协作与 CI 集成。

# 完整字段说明

```toml
# 模板身份（template_id 由 qshell 首次构建后自动回写，勿手动修改）
template_id = "tmpl-xxxxxxxxxxxx"
name        = "claude"

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
| template_id | string | 模板 ID，首次构建后自动回写 |
| name | string | 模板名称（仅创建时使用） |
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

# 查找规则

1. `--config <path>` 显式指定路径（仅 `build` 命令支持）
2. 当前工作目录下的 `qshell.sandbox.toml`
3. 未找到时按纯 CLI 模式运行（向后兼容）

# 自动回写行为

- 首次构建（配置文件存在、`template_id` 为空）成功后，qshell 自动把新的 `template_id` 写入文件
- 未指定 `--wait` 时，在构建成功启动后回写；指定 `--wait` 时，在构建完成且状态为 `ready` 后回写
- 回写会替换已有的 `template_id` 赋值行（无论原值是否为空），或在文件头插入一行
- 注释、字段顺序、空白均保留
- 回写完成后 stdout 打印：`Written template_id to <path> (please commit this file)`
- `template_id` 存在后再次执行 build 会进入 rebuild 流程，配置文件中保留的 `from_image` / `from_template` 仅作为首次创建记录保留，不参与 rebuild

# 团队协作约定

将 `qshell.sandbox.toml` 加入版本控制：
- 团队成员克隆仓库后直接 `qshell sandbox template build` 即可重建同一模板
- CI 脚本只需一行命令，无需维护散落的参数
- 首次 commit 由构建者完成；后续通常无需修改（除非切换模板）

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
# 构建完成后，qshell.sandbox.toml 中 template_id 已被回写
git add qshell.sandbox.toml && git commit -m "chore: record template_id"
```

下次重建：
```bash
qshell sandbox template build --wait  # 幂等，不再 409
```
