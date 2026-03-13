# qshell 下载与安装指南

## 安装步骤

### 1. 获取最新版本号

通过 GitHub API 获取最新版本号：
```bash
VERSION=$(curl -s "https://api.github.com/repos/qiniu/qshell/releases/latest" | grep '"tag_name":' | sed -E 's/.*"v([^"]+)".*/\1/')
if [ -z "$VERSION" ]; then
  echo "Error: Failed to detect latest version. Please check https://github.com/qiniu/qshell/releases manually." >&2
  exit 1
fi
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

if [ -z "$SUFFIX" ]; then
  echo "Error: Unsupported platform: $OS/$ARCH" >&2
  echo "Please download manually from https://github.com/qiniu/qshell/releases" >&2
  exit 1
fi

URL="https://kodo-toolbox-new.qiniu.com/qshell-v${VERSION}-${SUFFIX}.tar.gz"
curl -sL -e https://developer.qiniu.com -o /tmp/qshell.tar.gz "$URL"
```

### 3. 解压并安装到用户目录

```bash
tar -xzf /tmp/qshell.tar.gz -C /tmp/
chmod +x /tmp/qshell

# 安装到用户目录
mkdir -p "$HOME/.local/bin"
mv /tmp/qshell "$HOME/.local/bin/qshell"

rm -f /tmp/qshell.tar.gz
```

如果 `$HOME/.local/bin` 不在 PATH 中，需要将其添加到 shell 配置文件：
```bash
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc   # bash
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc    # zsh
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

重新执行上述安装步骤即可覆盖更新，会自动获取最新版本。
