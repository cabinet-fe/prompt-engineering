#!/bin/sh
# docs-server 一键安装脚本
# 用法:
#   curl -fsSL https://raw.githubusercontent.com/cabinet-fe/prompt-engineering/main/scripts/install-docs-server.sh | bash
# 或者指定安装目录 / 版本 / 代理:
#   INSTALL_DIR=/usr/local/bin VERSION=latest curl -fsSL ... | bash
# 安装完成后自动生成配置文件 /etc/docs-server.yaml（CONFIG_PATH 可覆盖，已存在则不覆盖），
# 并给出后台常驻（systemd / nohup）运行指引；前台直接运行用于临时验证。

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
  rm -f "${TMP_CONF:-}"
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

# 8. 准备数据目录（Linux 惯例 /var/lib，macOS 落在用户目录）
if [ "$OS_TYPE" = "darwin" ]; then
  DATA_DIR="${DOCS_DATA_DIR:-$HOME/.local/share/docs-server}"
else
  DATA_DIR="${DOCS_DATA_DIR:-/var/lib/docs-server}"
fi
if ! mkdir -p "$DATA_DIR" 2>/dev/null; then
  if has_sudo && sudo mkdir -p "$DATA_DIR" 2>/dev/null; then
    :
  else
    DATA_DIR="$HOME/.docs-server"
    mkdir -p "$DATA_DIR" 2>/dev/null || true
  fi
fi

# 9. 生成配置文件（推荐：之后一律以该文件配置；已存在则不覆盖）
CONFIG_PATH="${CONFIG_PATH:-/etc/docs-server.yaml}"
TMP_CONF="$(mktemp "${TMPDIR:-/tmp}/docs-server-config.XXXXXX")"
cat > "$TMP_CONF" <<EOF
# docs-server 配置（由一键安装脚本生成）
# 优先级：环境变量 > 本文件 > 默认值；仅经 Nginx 反代对外时建议把 addr 改为 "127.0.0.1:8080"
addr: ":8080"
db_path: ${DATA_DIR}/docs.db
push_token: ${GEN_PUSH_TOKEN}
EOF
CONFIG_NOTE="已生成（推送令牌已写入，权限 600）"
if [ -e "$CONFIG_PATH" ]; then
  CONFIG_NOTE="已存在，未覆盖，沿用原配置"
elif [ -w "$(dirname "$CONFIG_PATH")" ] || [ "$(id -u)" -eq 0 ]; then
  { mv "$TMP_CONF" "$CONFIG_PATH" && chmod 600 "$CONFIG_PATH"; } || CONFIG_FAIL=1
elif has_sudo; then
  { sudo mv "$TMP_CONF" "$CONFIG_PATH" && sudo chmod 600 "$CONFIG_PATH"; } || CONFIG_FAIL=1
elif [ -w "$HOME" ]; then
  CONFIG_PATH="$HOME/docs-server.yaml"
  { mv "$TMP_CONF" "$CONFIG_PATH" && chmod 600 "$CONFIG_PATH"; } || CONFIG_FAIL=1
  CONFIG_NOTE="默认位置无权限写入，已改生成到 HOME（可用 CONFIG_PATH 覆盖）"
else
  CONFIG_FAIL=1
fi
rm -f "$TMP_CONF"
if [ -n "${CONFIG_FAIL:-}" ]; then
  CONFIG_PATH=""
  CONFIG_NOTE="写入失败，可改用环境变量方式启动（见 README）"
fi

# 10. 输出安装结果与运行指引
echo ""
echo "========================================="
echo " 🎉 docs-server 安装成功！"
echo " 二进制路径: $TARGET_BIN"
if [ -n "$CONFIG_PATH" ]; then
  echo " 配置文件:   $CONFIG_PATH"
  echo "             ($CONFIG_NOTE)"
else
  echo " 配置文件:   未生成 ($CONFIG_NOTE)"
fi
echo "========================================="
echo ""
echo "🔑 推送令牌 (push_token):"
echo "   $GEN_PUSH_TOKEN"
echo ""
echo "💡 库仓库推送文档时，在仓库 .env 中配置 DOCS_TOKEN 为上述令牌（严禁提交入 git）。"

# 11. PATH 检查
case ":$PATH:" in
  *":$INSTALL_DIR:"*) ;;
  *)
    echo ""
    echo "⚠️  注意: $INSTALL_DIR 不在当前 PATH 中。建议将其加入 PATH："
    echo "   export PATH=\"$INSTALL_DIR:\$PATH\""
    ;;
esac

# 后续运行指引统一走配置文件
if [ -n "$CONFIG_PATH" ]; then
  RUN_CMD="$TARGET_BIN -config $CONFIG_PATH"
else
  RUN_CMD="DOCS_DB_PATH=$DATA_DIR/docs.db DOCS_PUSH_TOKEN=$GEN_PUSH_TOKEN $TARGET_BIN"
fi

echo ""
echo "-----------------------------------------"
echo "▶ 默认运行方式：后台常驻（推荐）"
echo "-----------------------------------------"
if [ "$OS_TYPE" = "linux" ] && [ -n "$CONFIG_PATH" ]; then
  echo "使用 systemd（开机自启、异常自动重启）："
  echo ""
  echo "  sudo tee /etc/systemd/system/docs-server.service >/dev/null <<'EOF'"
  echo "  [Unit]"
  echo "  After=network.target"
  echo ""
  echo "  [Service]"
  echo "  ExecStart=$RUN_CMD"
  echo "  Restart=on-failure"
  echo "  StateDirectory=docs-server"
  echo ""
  echo "  [Install]"
  echo "  WantedBy=multi-user.target"
  echo "  EOF"
  echo ""
  echo "  sudo systemctl daemon-reload"
  echo "  sudo systemctl enable --now docs-server"
  echo "  journalctl -u docs-server -f    # 查看日志"
  echo ""
  echo "无 systemd 的环境可用 nohup："
  echo "  nohup $RUN_CMD > /var/log/docs-server.log 2>&1 &"
elif [ -n "$CONFIG_PATH" ]; then
  echo "  nohup $RUN_CMD > \"\$HOME/docs-server.log\" 2>&1 &"
  echo "  tail -f \"\$HOME/docs-server.log\"    # 查看日志；停止: pkill docs-server"
  echo ""
  echo "如需开机自启，macOS 可配置 launchd。"
else
  echo "未生成配置文件，请按 README 以环境变量或自备配置文件启动。"
fi
echo ""
echo "-----------------------------------------"
echo "▶ 或直接前台运行（临时验证用，Ctrl+C 停止）"
echo "-----------------------------------------"
echo "  $RUN_CMD"
echo ""
echo "仅通过 Nginx 反向代理对外时：把配置文件中的 addr 改为 \"127.0.0.1:8080\"，"
echo "并在 Nginx location 中放宽请求体限制（client_max_body_size 50m），完整参考见 README。"
