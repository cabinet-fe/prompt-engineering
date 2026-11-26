#!/bin/sh
# docs-server 一键安装/更新脚本
# 用法:
#   curl -fsSL https://raw.githubusercontent.com/cabinet-fe/prompt-engineering/main/scripts/install-docs-server.sh | bash
# 或者指定安装目录 / 版本 / 代理:
#   INSTALL_DIR=. VERSION=latest curl -fsSL ... | bash
# 仅下载并安装或覆盖 docs-server 二进制文件，不自动生成配置文件。
# 首次部署请执行 `docs-server init` 生成当前目录下的 docs-server.yaml；
# 已安装二进制升级也可直接执行 `docs-server update`（大陆网络加 `-cn`）。

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
IS_UPGRADE=0
if [ -e "$TARGET_BIN" ]; then
  IS_UPGRADE=1
fi

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

# 7. 输出安装或更新结果（用 -version 读取版本，旧版也会在 flag 解析阶段安全退出）
INSTALLED_VER="$("$TARGET_BIN" -version 2>/dev/null || true)"

echo ""
echo "========================================="
if [ "$IS_UPGRADE" -eq 1 ]; then
  echo " 🎉 docs-server 更新成功！"
else
  echo " 🎉 docs-server 安装成功！"
fi
echo " 二进制路径: $TARGET_BIN"
if [ -n "$INSTALLED_VER" ]; then
  echo " 当前版本:   $INSTALLED_VER"
fi
echo "========================================="
echo ""

if [ "$IS_UPGRADE" -eq 1 ]; then
  echo "提示: 本脚本仅替换二进制，未改动任何配置文件与数据库。"
  echo "      若服务正在后台运行，请重启对应进程（如 sudo systemctl restart docs-server）以使新版本生效。"
  echo "      后续升级也可直接执行: $TARGET_BIN update（大陆网络加 -cn）"
else
  echo "下一步:"
  echo "  1. 在你的服务工作目录下生成配置文件（默认生成 ./docs-server.yaml，已存在时会询问是否覆盖）:"
  echo "     docs-server init"
  echo ""
  echo "  2. 使用配置文件启动服务:"
  echo "     docs-server -config ./docs-server.yaml"
  echo ""
  echo "  3. 后续升级直接原地更新即可（不改动配置与数据）:"
  echo "     docs-server update        # 大陆网络: docs-server update -cn"
fi

# 8. PATH 检查
case ":$PATH:" in
  *":$INSTALL_DIR:"*) ;;
  *)
    echo ""
    echo "⚠️  注意: $INSTALL_DIR 不在当前 PATH 中。建议将其加入 PATH："
    echo "   export PATH=\"$INSTALL_DIR:\$PATH\""
    ;;
esac
