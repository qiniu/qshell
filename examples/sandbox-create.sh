#!/bin/bash
# sandbox 创建示例
#
# 使用前请确保：
#   1. 已配置 QINIU_API_KEY 或 E2B_API_KEY
#   2. 已编译 qshell 或使用已安装的 qshell 命令
#   3. 已准备可用的沙箱模板 ID

set -e

QSHELL="${QSHELL:-qshell}"
TEMPLATE_ID="${TEMPLATE_ID:-my-template}"

echo "=== 示例 1: 创建并连接到沙箱 ==="
echo "${QSHELL} sbx cr ${TEMPLATE_ID}"
# ${QSHELL} sbx cr "${TEMPLATE_ID}"

echo ""
echo "=== 示例 2: 创建时附加已存在的注入规则 ==="
echo "${QSHELL} sbx cr ${TEMPLATE_ID} --injection-rule rule-openai --injection-rule rule-http"
# ${QSHELL} sbx cr "${TEMPLATE_ID}" --injection-rule rule-openai --injection-rule rule-http

echo ""
echo "=== 示例 3: 创建时附加内联 OpenAI 注入配置 ==="
echo "${QSHELL} sbx cr ${TEMPLATE_ID} --inline-injection 'type=openai,api-key=sk-xxx'"
# ${QSHELL} sbx cr "${TEMPLATE_ID}" --inline-injection 'type=openai,api-key=sk-xxx'

echo ""
echo "=== 示例 4: 创建时附加多个内联注入配置 ==="
echo "${QSHELL} sbx cr ${TEMPLATE_ID} --inline-injection 'type=gemini,api-key=sk-gem' --inline-injection 'type=http,base-url=https://api.example.com,headers=Authorization=Bearer token,X-Env=prod'"
# ${QSHELL} sbx cr "${TEMPLATE_ID}" \
#   --inline-injection 'type=gemini,api-key=sk-gem' \
#   --inline-injection 'type=http,base-url=https://api.example.com,headers=Authorization=Bearer token,X-Env=prod'

echo ""
echo "取消注释上述命令即可执行，或设置 QSHELL/TEMPLATE_ID 环境变量："
echo "  QSHELL=./qshell TEMPLATE_ID=my-template bash ${0}"
