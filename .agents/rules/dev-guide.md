# qshell — 开发维护指引

七牛云命令行工具（`github.com/qiniu/qshell/v2`），Go 1.24+，基于 Cobra 框架。

完整开发维护指南见项目根目录 `AGENTS.md`。

## 项目结构

- `main/` — 程序入口（`go build ./main/`）
- `cmd/` — Cobra 命令定义（user、bucket、rs、upload、download、cdn、sandbox 等）
- `cmd_test/` — 命令层测试
- `iqshell/` — 核心业务逻辑
  - `common/` — 通用工具（版本、配置、日志、账号）
  - `storage/` — 对象存储操作
  - `cdn/` — CDN 操作
  - `sandbox/` — 沙箱环境操作
  - `ali/` — 阿里云 OSS 迁移
  - `aws/` — AWS S3 迁移
- `docs/` — 命令文档（`.md` 文件）和文档注册（`docs.go`）
- `examples/` — 使用示例

## 架构模式

- **Cobra 命令模式**：`main/ → cmd.Execute() → cmd/*.go → iqshell/`
- **分层架构**：cmd（参数解析）→ iqshell（业务逻辑）→ go-sdk（API 调用）
- **配置管理**：Viper + LevelDB（`~/.qshell/`）
- **版本注入**：ldflags `-X version.version=vX.Y.Z`

## 编码规范

- `gofmt -s` 格式化（**CI 强制检查**）
- `make lint`（`go vet` + `staticcheck`）— 本地提交前检查，不在 CI 中运行
- 注释使用**中文**，导出标识符注释以名称开头（godoc 规范）
- 错误使用 `fmt.Errorf("context: %w", err)` 包装

## 测试要求

- `make test` — 运行所有测试（提交前必须通过）
- `make test-unit` — 运行 unit 标签测试
- `make lint` — 静态检查（提交前必须通过）
- 使用 `testify/assert` 断言，推荐表驱动测试

## CI 流程

gofmt 检查 → 单元测试 + 覆盖率上传（Go 1.24.x，ubuntu-latest）

## 关键约定

- 构建入口是 `./main/`，不是项目根目录
- 版本号通过 ldflags 注入，不要硬编码
- 新增/修改命令时同步更新 `docs/` 文档和 `CHANGELOG.md`
- 跨平台路径使用 `filepath.Join`
