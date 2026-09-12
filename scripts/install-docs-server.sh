#!/bin/sh
# docs-server 一键安装脚本
# 用法:
#   curl -fsSL https://raw.githubusercontent.com/cabinet-fe/prompt-engineering/main/scripts/install-docs-server.sh | bash
# 或者指定安装目录 / 版本 / 代理:
#   INSTALL_DIR=/usr/local/bin VERSION=latest curl -fsSL ... | bash

set -e

REPO="${REPO:-cabinet-fe/prompt-engineering}"
VERSION="${VERSION:-latest}"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
GH_PROXY="${GH_PROXY:-}"

# 规范化代理前缀（若提供）
if [ -n "$GH_PROXY" ]; then
  case "$GH_PROXY" in
    */) ;;
    *) GH_PROXY="${GH_PROXY}/" ;;
  esac
fi

# 1. 检测操作系统
OS="$(uname -s)"
case "$OS" in
  Linux)
    OS_TYPE="linux"
    ;;
  Darwin)
    OS_TYPE="darwin"
    ;;
  *)
    echo "错误: 不支持的操作系统 '$OS'。docs-server 仅支持 Linux 和 macOS (Darwin)。" >&2
    exit 1
    ;;
esac

# 2. 检测 CPU 架构
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64|amd64)
    ARCH_TYPE="x64"
    ;;
  aarch64|arm64)
    ARCH_TYPE="arm64"
    ;;
  *)
    echo "错误: 不支持的 CPU 架构 '$ARCH'。docs-server 仅支持 x64 (x86_64) 和 arm64 (aarch64)。" >&2
    exit 1
    ;;
esac

ASSET_NAME="docs-server-${OS_TYPE}-${ARCH_TYPE}"

# 3. 构造下载地址
if [ "$VERSION" = "latest" ]; then
  DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${ASSET_NAME}"
else
  DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${ASSET_NAME}"
fi

if [ -n "$GH_PROXY" ]; then
  DOWNLOAD_URL="${GH_PROXY}${DOWNLOAD_URL}"
fi

# 4. 检测下载工具
download_file() {
  target_url="$1"
  dest_file="$2"
  if command -v curl >/dev/null 2>&1; then
    curl -fL --progress-bar -o "$dest_file" "$target_url"
  elif command -v wget >/dev/null 2>&1; then
    if wget --help 2>&1 | grep -q -- '--show-progress'; then
      wget -q --show-progress -O "$dest_file" "$target_url"
    else
      wget -q -O "$dest_file" "$target_url"
    fi
  else
    echo "错误: 未检测到 curl 或 wget，请先安装其中之一后重试。" >&2
    exit 1
  fi
}

has_sudo() {
  command -v sudo >/dev/null 2>&1 && [ "$(id -u)" -ne 0 ]
}

# 5. 下载到临时目录
TMP_DIR="$(mktemp -d 2>/dev/null || mktemp -d -t 'docs-server-install')"
cleanup() {
  rm -rf "$TMP_DIR"
}
trap cleanup EXIT INT TERM

TMP_BIN="${TMP_DIR}/${ASSET_NAME}"

echo ">> 正在从 GitHub Releases 获取 docs-server (${OS_TYPE}/${ARCH_TYPE})..."
echo ">> 下载地址: $DOWNLOAD_URL"

if ! download_file "$DOWNLOAD_URL" "$TMP_BIN"; then
  echo "错误: 下载失败。若处于网络受限环境，可尝试设置 GH_PROXY 环境变量后重试，例如:" >&2
  echo "  GH_PROXY=https://ghfast.top/ curl -fsSL ... | bash" >&2
  exit 1
fi

chmod +x "$TMP_BIN"

# 6. 安装到目标目录
TARGET_BIN="${INSTALL_DIR}/docs-server"

if [ ! -d "$INSTALL_DIR" ]; then
  if [ -w "$(dirname "$INSTALL_DIR" 2>/dev/null || echo ".")" ] || [ "$(id -u)" -eq 0 ]; then
    mkdir -p "$INSTALL_DIR"
  elif has_sudo; then
    echo ">> 正在通过 sudo 创建目录 $INSTALL_DIR..."
    sudo mkdir -p "$INSTALL_DIR"
  else
    echo "错误: 目标目录 $INSTALL_DIR 不存在且无权限创建。" >&2
    echo "提示: 请使用 root / sudo 运行，或设置 INSTALL_DIR（如 INSTALL_DIR=\$HOME/.local/bin）。" >&2
    exit 1
  fi
fi

if [ -w "$INSTALL_DIR" ] || [ "$(id -u)" -eq 0 ]; then
  mv "$TMP_BIN" "$TARGET_BIN"
  chmod +x "$TARGET_BIN"
elif has_sudo; then
  echo ">> 目标目录需要 root 权限，正在通过 sudo 安装到 $TARGET_BIN..."
  sudo mv "$TMP_BIN" "$TARGET_BIN"
  sudo chmod +x "$TARGET_BIN"
else
  echo "错误: 对目录 $INSTALL_DIR 无写入权限。" >&2
  echo "提示: 请使用 root / sudo 运行，或设置 INSTALL_DIR，例如:" >&2
  echo "  INSTALL_DIR=\$HOME/.local/bin sh install-docs-server.sh" >&2
  exit 1
fi

# 7. 生成或读取推送令牌 (push_token)
GEN_PUSH_TOKEN="${DOCS_PUSH_TOKEN:-${PUSH_TOKEN:-}}"
if [ -z "$GEN_PUSH_TOKEN" ]; then
  if command -v openssl >/dev/null 2>&1; then
    GEN_PUSH_TOKEN="$(openssl rand -hex 32)"
  elif [ -r /dev/urandom ]; then
    GEN_PUSH_TOKEN="$(head -c 32 /dev/urandom | od -An -tx1 | tr -d ' \t\n')"
  else
    GEN_PUSH_TOKEN="$(date +%s%N 2>/dev/null || date +%s | shasum -a 256 2>/dev/null || sha256sum 2>/dev/null | head -c 64)"
  fi
fi

echo ""
echo "========================================="
echo " 🎉 docs-server 安装成功！"
echo " 二进制路径: $TARGET_BIN"
echo "========================================="
echo ""
echo "🔑 生成的推送令牌 (DOCS_PUSH_TOKEN):"
echo "   $GEN_PUSH_TOKEN"
echo ""
echo "💡 提示: 该令牌用于推送文档时的鉴权（Bearer Token），请妥善保存。"

# 8. PATH 检查
case ":$PATH:" in
  *":$INSTALL_DIR:"*) ;;
  *)
    echo ""
    echo "⚠️  注意: $INSTALL_DIR 不在当前 PATH 中。建议将其加入 PATH："
    echo "   export PATH=\"$INSTALL_DIR:\$PATH\""
    ;;
esac

echo ""
echo "-----------------------------------------"
echo "运行方式 A：仅在本机监听（通过 Nginx 反向代理暴露）"
echo "-----------------------------------------"
echo "  DOCS_ADDR=127.0.0.1:8080 \\"
echo "  DOCS_DB_PATH=/var/lib/docs-server/docs.db \\"
echo "  DOCS_PUSH_TOKEN=$GEN_PUSH_TOKEN \\"
echo "  docs-server"
echo ""
echo "  配套 Nginx 配置参考："
echo "    location / {"
echo "        proxy_pass http://127.0.0.1:8080;"
echo "        proxy_set_header Host \$host;"
echo "        proxy_set_header X-Real-IP \$remote_addr;"
echo "        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;"
echo "        proxy_set_header X-Forwarded-Proto \$scheme;"
echo "        client_max_body_size 50m; # 支持大文档库推送"
echo "    }"
echo ""
echo "-----------------------------------------"
echo "运行方式 B：直接监听所有网卡 (:8080)"
echo "-----------------------------------------"
echo "  DOCS_ADDR=:8080 \\"
echo "  DOCS_DB_PATH=/var/lib/docs-server/docs.db \\"
echo "  DOCS_PUSH_TOKEN=$GEN_PUSH_TOKEN \\"
echo "  docs-server"
echo ""
echo "-----------------------------------------"
echo "或使用配置文件启动："
echo "  docs-server -config /etc/docs-server.yaml"
echo ""
