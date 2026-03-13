# qshell 下载与安装指南

## 安装步骤

### 1. 获取最新版本号

从七牛官方文档页面获取最新版本号：
```bash
VERSION=$(curl -sL -e https://developer.qiniu.com "https://developer.qiniu.com/kodo/1302/qshell" | grep -oE 'qshell-v[0-9]+\.[0-9]+\.[0-9]+' | head -1 | sed 's/^qshell-v//')
echo "最新版本: v${VERSION}"
```

### 2. 检测平台和架构，下载二进制

下载地址：`https://kodo-toolbox-new.qiniu.com/qshell-v${VERSION}-${SUFFIX}.tar.gz`（必须带 Referer `-e https://developer.qiniu.com`）

```bash
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

if [ "$OS" = "darwin" ]; then
  if [ "$ARCH" = "arm64" ]; then
    SUFFIX="darwin-arm64"
  else
    SUFFIX="darwin-amd64"
  fi
elif [ "$OS" = "linux" ]; then
  if [ "$ARCH" = "x86_64" ]; then
    SUFFIX="linux-amd64"
  elif [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
    SUFFIX="linux-arm64"
  elif [ "$ARCH" = "i386" ] || [ "$ARCH" = "i686" ]; then
    SUFFIX="linux-386"
  fi
fi

URL="https://kodo-toolbox-new.qiniu.com/qshell-v${VERSION}-${SUFFIX}.tar.gz"
curl -sL -e https://developer.qiniu.com -o /tmp/qshell.tar.gz "$URL"
```

### 3. 解压并安装到系统 PATH

```bash
tar -xzf /tmp/qshell.tar.gz -C /tmp/
chmod +x /tmp/qshell

# 安装到用户可写的 PATH 目录
if [ -d "$HOME/.local/bin" ]; then
  mv /tmp/qshell "$HOME/.local/bin/qshell"
elif [ -d "/usr/local/bin" ] && [ -w "/usr/local/bin" ]; then
  mv /tmp/qshell /usr/local/bin/qshell
else
  sudo mv /tmp/qshell /usr/local/bin/qshell
fi

rm -f /tmp/qshell.tar.gz
```

### 4. 验证安装

```bash
qshell version
```

## 配置账号

安装完成后配置七牛账号：
```bash
qshell account <AccessKey> <SecretKey> <Name>
```

- `AccessKey` 和 `SecretKey`：从七牛控制台获取（https://portal.qiniu.com/user/key）
- `Name`：自定义名称，用于本地区分多个账号

### 账号管理

```bash
# 查看当前账号
qshell account

# 列出所有已配置账号
qshell user ls

# 切换账号
qshell user cu <Name>
```

## 更新 qshell

重新执行上述安装步骤即可覆盖更新，会自动从官方文档获取最新版本。
