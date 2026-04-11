# qshell — 开发维护指南

qshell 是七牛云官方命令行工具，模块路径 `github.com/qiniu/qshell/v2`，Go 1.24+，基于 Cobra 框架构建。

本文档面向维护者和贡献者，帮助理解项目结构、架构模式和开发流程。
命令使用说明请参考 `docs/` 目录和 `README.md`。

## 项目结构

### 顶层目录

| 目录/文件 | 说明 |
|-----------|------|
| `main/` | 程序入口，调用 `cmd.Execute()` |
| `cmd/` | Cobra 命令定义（每个文件对应一组命令） |
| `cmd_test/` | 命令层集成/单元测试 |
| `iqshell/` | 核心业务逻辑 |
| `docs/` | 命令文档（`.md` 文件）和文档注册（`docs.go`） |
| `examples/` | 使用示例脚本 |
| `.github/` | Issue/PR 模板、CI 工作流 |

### cmd/ 命令文件

| 文件 | 职责 |
|------|------|
| `root.go` | 根命令定义、全局 flag |
| `user.go` | 账号管理（account、user） |
| `bucket.go` | Bucket 管理 |
| `rs.go` | 单文件操作（stat、delete、move、copy、rename 等） |
| `rsbatch.go` | 批量文件操作 |
| `upload.go` | 文件上传（表单上传、分片上传） |
| `download.go` | 文件下载 |
| `cdn.go` | CDN 操作（刷新、预取） |
| `fop.go` | 数据处理（pfop、prefop） |
| `asyncfetch.go` | 异步抓取 |
| `sandbox.go` | 沙箱环境管理 |
| `sandbox_template.go` | 沙箱模板管理 |
| `ali.go` | 阿里云 OSS 数据迁移 |
| `aws.go` | AWS S3 数据迁移 |
| `tools.go` | 工具命令（base64、urlencode、token 等） |
| `match.go` | 文件匹配测试 |
| `servers.go` | 文件服务 |
| `share.go` | 文件分享 |
| `m3u8.go` | M3U8 操作 |
| `version.go` | 版本信息 |
| `autocompletion.go` | Shell 自动补全 |

### iqshell/ 核心业务

| 包 | 职责 |
|---|------|
| `iqshell/common/` | 通用工具（版本、配置、日志、账号、数据结构等） |
| `iqshell/storage/` | 七牛对象存储操作实现 |
| `iqshell/cdn/` | CDN 操作实现 |
| `iqshell/sandbox/` | 沙箱环境操作实现（基于七牛 Go SDK sandbox 包） |
| `iqshell/ali/` | 阿里云 OSS 迁移实现 |
| `iqshell/aws/` | AWS S3 迁移实现 |

### 其他重要文件

| 文件 | 说明 |
|------|------|
| `iqshell/load.go` | 命令加载器接口定义 |
| `iqshell/common/version/version.go` | 版本号（通过 ldflags 注入） |
| `readme.go` | 通过 `go:embed` 嵌入 README.md |
| `CHANGELOG.md` | 变更日志 |
| `Makefile` | 构建、测试、发布目标 |

## 核心架构模式

### Cobra 命令模式

所有命令基于 [spf13/cobra](https://github.com/spf13/cobra) 框架：

```
main/main.go → cmd.Execute() → cmd/root.go（根命令）→ 各子命令
```

命令注册流程：
1. 每个 `cmd/*.go` 文件定义命令加载函数
2. `cmd/root.go` 中注册所有命令加载器
3. 命令执行时调用 `iqshell/` 中的业务逻辑

### 分层架构

```
cmd/（命令定义层）→ iqshell/（业务逻辑层）→ go-sdk/（SDK 层）
```

- `cmd/`：参数解析、flag 定义、命令注册
- `iqshell/`：业务逻辑实现、配置管理、输出格式化
- `go-sdk/`：底层 API 调用、认证、HTTP 请求

### 配置管理

- 使用 [spf13/viper](https://github.com/spf13/viper) 管理配置
- 账号信息存储在 `~/.qshell/` 目录（LevelDB）
- 支持多账号管理

### 版本注入

版本号通过 `ldflags` 在构建时注入：

```bash
go build -ldflags '-X github.com/qiniu/qshell/v2/iqshell/common/version.version=vX.Y.Z' ./main/
```

## 编码规范

### 格式化

- 使用 `gofmt -s` 格式化代码（**CI 强制检查**）
- 提交前确保 `gofmt -s -l .` 无输出

### 静态检查（本地提交前）

- 使用 `make lint` 运行 `go vet` + `staticcheck`
- 注意：`vet` 和 `staticcheck` 目前**不在 CI 中运行**，仅作为本地提交前检查
- CI 仅强制检查 `gofmt` 和单元测试

### 注释规范

- 注释默认使用**中文**
- 导出的 API、函数、类型、常量必须添加注释
- 注释以被注释的标识符名称开头（符合 Go 官方 godoc 规范）

### 命名规范

- 命令文件按功能命名：`cmd/<功能>.go`
- 业务逻辑按模块组织：`iqshell/<模块>/`
- 测试文件放在 `cmd_test/` 或对应包内

### 错误处理

- 错误使用 Go 标准的 `fmt.Errorf("context: %w", err)` 包装
- 使用 `errors.Is()` / `errors.As()` 判断错误类型
- 命令层统一处理错误输出

## 测试规范

### Build Tags

测试文件可选添加 build tag：

```go
//go:build unit

package xxx_test
```

- `unit` — 单元测试，不依赖外部服务
- `integration` — 集成测试，需要七牛账号和凭证

### Makefile 测试命令

| 命令 | 说明 |
|------|------|
| `make test` | 运行所有测试（无 build tag） |
| `make test-unit` | 运行 unit 标签测试 |
| `make test-integration` | 运行集成测试（需要凭证，超时 600s） |
| `make test-sandbox-unit` | Sandbox 单元测试 |
| `make test-sandbox-integration` | Sandbox 集成测试（需要 `QINIU_API_KEY`） |
| `make test-sandbox` | 所有 Sandbox 测试 |

### 测试实践

- 使用 `testify/assert` 进行断言
- 推荐使用表驱动测试
- 集成测试需要配置环境变量（七牛 AK/SK 等）

## CI 流程

CI 在 push 和 pull_request 时触发（`.github/workflows/ut-check.yml`）：

1. **gofmt 检查**：`gofmt -s -l .` 检查未格式化文件
2. **单元测试**：`go test -coverprofile=coverage.txt ./...`
3. **代码覆盖率**：上传至 Codecov

Go 版本矩阵：1.24.x

### 发布流程

发布通过 GitHub Release 触发（`.github/workflows/release.yaml`）：

1. 构建 16 个平台二进制文件（darwin/linux/windows × 多架构）
2. 上传至 GitHub Releases
3. 上传至七牛云 devtools 存储

## 开发工作流

### 环境准备

```bash
# 验证当前状态
make test
make lint
```

### 日常开发

1. 修改代码前先运行 `make test` 确认当前状态
2. 新增命令时在 `cmd/` 创建命令定义，在 `iqshell/` 实现业务逻辑
3. 提交前运行 `make test` 和 `make lint` 确保通过
4. 代码格式化：`gofmt -s -w .`
5. 同步更新 `docs/` 中的命令文档

### 新增命令的标准流程

1. 在 `cmd/<功能>.go` 中定义 Cobra 命令和 flag
2. 在 `iqshell/<模块>/` 中实现业务逻辑
3. 在 `cmd/root.go` 中注册命令加载器
4. 在 `docs/` 中添加命令文档（`.md` 文件），并在 `docs/docs.go` 中添加 embed 声明和注册
5. 在 `cmd_test/` 中添加测试
6. 更新 `CHANGELOG.md`

## 重要约定

1. **入口在 `./main/`** — 构建时使用 `go build ./main/`，不是 `go build .`
2. **版本号通过 ldflags 注入** — 不要硬编码版本号，通过 `-ldflags` 设置 `version.version`
3. **文档同步** — 新增或修改命令时必须同步更新 `docs/` 目录
4. **依赖七牛 Go SDK** — 核心存储和沙箱功能依赖 `github.com/qiniu/go-sdk/v7`
5. **跨平台兼容** — 支持 16 个平台，使用 `filepath.Join` 等跨平台写法
6. **保持向后兼容** — 命令行接口变更需要考虑用户兼容性
7. **CHANGELOG 更新** — 功能变更必须记录在 `CHANGELOG.md` 中
