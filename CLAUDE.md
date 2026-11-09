# qshell — Claude Code Instructions

## 项目概要

七牛云命令行工具 qshell，模块路径 `github.com/qiniu/qshell/v2`，Go 1.25+，基于 Cobra 框架。

## 核心开发规范

@AGENTS.md

## Claude Code 特定指引

### 常用命令

| 命令 | 说明 | 时机 |
|------|------|------|
| `make test` | 运行所有测试 | 修改前确认状态 + 提交前必须通过 |
| `make lint` | 运行静态检查（vet + staticcheck） | 提交前必须通过 |
| `make test-unit` | 运行 unit 标签测试 | 修改 cmd_test/ 或 iqshell/ 后 |
| `make test-sandbox-unit` | 运行 Sandbox 单元测试 | 修改 sandbox 相关代码后 |
| `gofmt -s -w .` | 格式化代码 | 提交前 |
| `go build ./main/` | 构建二进制 | 验证编译是否通过 |

### 开发注意事项

- 修改代码前先运行 `make test` 确认当前状态
- 构建入口是 `./main/`，不是项目根目录
- 版本号通过 ldflags 注入到 `iqshell/common/version.version`，不要硬编码
- 注释默认使用**中文**，导出标识符注释以标识符名称开头
- 新增或修改命令时同步更新 `docs/` 目录的文档
- 更新 `CHANGELOG.md` 记录功能变更
- 跨平台文件路径使用 `filepath.Join`，不要硬编码路径分隔符
