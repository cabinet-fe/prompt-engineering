#!/bin/sh
# docs-server 中国大陆加速安装引导脚本
# 大陆网络直连 raw.githubusercontent.com 通常不可达，本脚本按序尝试多个 GitHub 加速代理，
# 择优下载正式安装脚本（install-docs-server.sh），并让 Release 二进制走同一个可用代理。
# 用法（URL 自带加速前缀，整条命令在大陆网络可直接执行）:
#   curl -fsSL https://ghfast.top/https://raw.githubusercontent.com/cabinet-fe/prompt-engineering/main/scripts/install-docs-server-cn.sh | bash
# 固定或自定义加速源（空格分隔按序尝试，直连始终兜底）:
#   GH_PROXY=https://gh-proxy.com/ curl -fsSL ... | bash

set -e

REPO="${REPO:-cabinet-fe/prompt-engineering}"
MAIN_URL="https://raw.githubusercontent.com/${REPO}/main/scripts/install-docs-server.sh"

norm_prefix() {
  case "$1" in
    */) printf '%s' "$1" ;;
    *) printf '%s/' "$1" ;;
  esac
}

fetch() {
  if command -v curl >/dev/null 2>&1; then
    curl -fsSL -m 30 "$1"
  elif command -v wget >/dev/null 2>&1; then
    wget -q -T 30 -O - "$1"
  else
    echo "错误: 未检测到 curl 或 wget，请先安装其中之一后重试。" >&2
    exit 1
  fi
}

# 组装候选列表：指定了 GH_PROXY 用指定的，否则用内置公共加速源；直连始终兜底
if [ -n "${GH_PROXY:-}" ]; then
  MIRRORS=""
  for m in $GH_PROXY; do
    MIRRORS="$MIRRORS $(norm_prefix "$m")"
  done
else
  MIRRORS="https://ghfast.top/ https://gh-proxy.com/ https://ghproxy.net/ https://gh.llkk.cc/ https://ghproxy.cn/"
fi

TMP_FILE="$(mktemp "${TMPDIR:-/tmp}/docs-server-install-cn.XXXXXX")"
cleanup() {
  rm -f "$TMP_FILE"
}
trap cleanup EXIT INT TERM

echo ">> docs-server 大陆加速安装：正在选择可用的下载源..."

for prefix in $MIRRORS ""; do
  if [ -n "$prefix" ]; then
    SOURCE_DESC="加速源 $prefix"
    echo ">> 尝试加速源: $prefix"
  else
    SOURCE_DESC="直连"
    echo ">> 尝试直连 GitHub..."
  fi

  if fetch "${prefix}${MAIN_URL}" > "$TMP_FILE" 2>/dev/null && [ -s "$TMP_FILE" ]; then
    # 全角标点紧贴 $var 时 macOS bash 会吞变量内容，含变量的行一律用 printf
    printf '>> 下载成功（%s），开始安装...\n' "$SOURCE_DESC"
    # 把可用的代理前缀传给正式安装脚本，Release 二进制走同一条路
    GH_PROXY="$prefix"
    export GH_PROXY
    if sh "$TMP_FILE"; then
      exit 0
    fi
    printf '错误: 安装脚本执行失败（%s）。\n' "$SOURCE_DESC" >&2
    echo "可换一个加速源重试，例如:" >&2
    echo "  GH_PROXY=https://gh-proxy.com/ sh install-docs-server-cn.sh" >&2
    exit 1
  fi
  echo ">> 不可用，切换下一个源..."
done

echo "错误: 所有下载源均不可用，请检查网络或指定加速代理后重试，例如:" >&2
echo "  GH_PROXY=https://gh-proxy.com/ sh install-docs-server-cn.sh" >&2
exit 1
