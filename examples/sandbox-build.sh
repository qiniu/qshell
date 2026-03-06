#!/bin/bash
# sandbox 模板 Dockerfile 构建示例
#
# 使用前请确保：
#   1. 已配置 QINIU_API_KEY 或通过 qshell account 登录
#   2. 已编译 qshell 或使用已安装的 qshell 命令

set -e

QSHELL="${QSHELL:-qshell}"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

# 示例 1: 从简单 Dockerfile 构建（无 COPY 指令）
echo "=== 示例 1: 简单 Dockerfile 构建 ==="
echo "${QSHELL} sbx tpl build --name test-simple --dockerfile ${SCRIPT_DIR}/sandbox/Dockerfile.simple --wait"
# ${QSHELL} sbx tpl build --name test-simple --dockerfile "${SCRIPT_DIR}/sandbox/Dockerfile.simple" --wait

# 示例 2: 从带 COPY 指令的 Dockerfile 构建
echo ""
echo "=== 示例 2: 带 COPY 的 Dockerfile 构建 ==="
echo "${QSHELL} sbx tpl build --name test-copy --dockerfile ${SCRIPT_DIR}/sandbox/Dockerfile --path ${SCRIPT_DIR}/sandbox --wait"
# ${QSHELL} sbx tpl build --name test-copy --dockerfile "${SCRIPT_DIR}/sandbox/Dockerfile" --path "${SCRIPT_DIR}/sandbox" --wait

# 示例 3: 重新构建已有模板（无缓存）
echo ""
echo "=== 示例 3: 重新构建已有模板 ==="
echo "${QSHELL} sbx tpl build --template-id <TEMPLATE_ID> --dockerfile ${SCRIPT_DIR}/sandbox/Dockerfile.simple --no-cache --wait"

# 示例 4: 从 Docker 镜像直接构建
echo ""
echo "=== 示例 4: 从 Docker 镜像构建 ==="
echo "${QSHELL} sbx tpl build --name test-image --from-image m.daocloud.io/docker.io/library/ubuntu:22.04 --wait"

echo ""
echo "取消注释上述命令即可执行，或设置 QSHELL 环境变量指定 qshell 路径："
echo "  QSHELL=./qshell-UNSTABLE-darwin-arm64 bash ${0}"
