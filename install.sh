#!/bin/bash
# ==============================================================================
# fs 一键安装脚本
# 支持: Linux, macOS, Windows (Git Bash / MSYS2 / WSL)
# 所有平台均使用 .tar.gz 发布包
# 安装路径: $HOME/.fs
# 用法: curl -fsSL https://example.com/install.sh | bash
# ==============================================================================

set -euo pipefail
IFS=$'\n\t'

# -----------------------------
# 🔧 配置区域
# -----------------------------

SOFTWARE_NAME="fs"
GITHUB_REPO="liulei152/fs"
DEFAULT_VERSION="1.0.2"

FS_ROOT="${HOME}/.fs"
BIN_DIR="${FS_ROOT}/bin"
AUTO_CONFIG_SHELL=true

# 临时目录（退出时清理）
TEMP_DIR=""

# -----------------------------
# 🎨 输出函数
# -----------------------------

info() {
    printf '\033[1;34m[i]\033[0m %s\n' "$*"
}

success() {
    printf '\033[1;32m[✓]\033[0m %s\n' "$*"
}

error() {
    printf '\033[1;31m[✗]\033[0m %s\n' "$*" >&2
    cleanup
    exit 1
}

warn() {
    printf '\033[1;33m[!]\033[0m %s\n' "$*"
}

# -----------------------------
# 🧹 清理函数
# -----------------------------

cleanup() {
    if [ -n "${TEMP_DIR:-}" ] && [ -d "$TEMP_DIR" ]; then
        info "🧹 清理临时文件..."
        rm -rf "$TEMP_DIR"
    fi
}

# -----------------------------
# 🖥️ 系统探测函数
# -----------------------------

get_os() {
    local os
    os=$(uname -s | tr '[:upper:]' '[:lower:]')
    case "$os" in
        linux*)   echo "linux" ;;
        darwin*)  echo "darwin" ;;
        mingw*|cygwin*) echo "windows" ;;
        *)        error "不支持的操作系统: $os" ;;
    esac
}

get_arch() {
    local arch
    arch=$(uname -m | tr '[:upper:]' '[:lower:]')
    case "$arch" in
        x86_64|amd64)   echo "amd64" ;;
        aarch64|arm64)  echo "arm64" ;;
        i?86|x86)       echo "386" ;;
        armv*)          echo "arm" ;;
        *)              error "不支持的架构: $arch" ;;
    esac
}

# -----------------------------
# 🐚 配置 Shell 环境（bash/zsh/fish）
# -----------------------------

configure_shell() {
    local shell_rc=""
    local shell_name=$(basename "${SHELL:-unknown}")

    case "$shell_name" in
        "bash")
            shell_rc="${HOME}/.bashrc"
            ;;
        "zsh")
            shell_rc="${HOME}/.zshrc"
            ;;
        "fish")
            local fish_conf="${HOME}/.config/fish/config.fish"
            mkdir -p "$(dirname "$fish_conf")"
            if ! grep -qF '# fs shell setup' "$fish_conf" 2>/dev/null; then
                {
                    echo ""
                    echo "# fs shell setup"
                    echo "set -gx PATH \"\$HOME/.fs/bin\" \$PATH"
                    echo "# fs end"
                } >> "$fish_conf"
                success "已为 fish 配置 PATH"
            fi
            return
            ;;
        *)
            warn "未知 shell: $shell_name，跳过自动配置"
            return
            ;;
    esac

    local marker="# fs shell setup (auto generated)"
    local export_line="export PATH=\"\$HOME/.fs/bin:\$PATH\""

    if [ ! -f "$shell_rc" ]; then
        touch "$shell_rc"
    fi

    if grep -qF "$marker" "$shell_rc" 2>/dev/null; then
        info "Shell 已配置，跳过。"
        return
    fi

    info "配置 $shell_rc 以添加 PATH..."

    {
        echo ""
        echo "$marker"
        echo "$export_line"
        echo "# fs end"
    } >> "$shell_rc"

    success "已配置 $shell_rc"
}

# -----------------------------
# 💡 智能提示函数（核心优化）
# -----------------------------

show_post_install_tip() {
    local shell_name
    shell_name=$(basename "${SHELL:-unknown}" 2>/dev/null || echo "unknown")
    local config_file=""

    case "$shell_name" in
        bash)
            config_file="$HOME/.bashrc"
            ;;
        zsh)
            config_file="$HOME/.zshrc"
            ;;
        fish)
            config_file="$HOME/.config/fish/config.fish"
            ;;
        *)
            # 通用 fallback（适用于 dash, sh 等）
            config_file="$HOME/.profile"
            ;;
    esac

    if [ -n "$config_file" ]; then
        echo "💡 如果命令未生效，请运行: source $config_file"
    else
        echo "💡 如果命令未生效，请重启终端或手动将 \$HOME/.fs/bin 加入 PATH"
    fi
}

# -----------------------------
# 🚀 主安装逻辑
# -----------------------------

main() {
    local version="${1:-$DEFAULT_VERSION}"
    local os=$(get_os)
    local arch=$(get_arch)

    # ✅ 所有平台都使用 .tar.gz
    local archive_name="${os}-${arch}.tar.gz"
    local url="https://gitee.com/${GITHUB_REPO}/releases/download/v${version}/${archive_name}"

    info "准备安装 ${SOFTWARE_NAME} v${version}"
    info "平台: ${os}/${arch}"
    info "下载: ${url}"

    # 创建临时目录
    TEMP_DIR=$(mktemp -d 2>/dev/null || mktemp -d -t fs-install)
    info "使用临时目录: $TEMP_DIR"
    cd "$TEMP_DIR"

    # 下载 .tar.gz
    info "📥 正在下载..."
    if command -v wget >/dev/null 2>&1; then
        wget -q --show-progress -O "fs.tar.gz" "$url" || \
            wget -q -O "fs.tar.gz" "$url"
    elif command -v curl >/dev/null 2>&1; then
        curl -f# -L -o "fs.tar.gz" "$url"
    else
        error "需要 curl 或 wget"
    fi

    # 创建安装目录
    info "🔧 创建安装目录: $BIN_DIR"
    mkdir -p "$BIN_DIR"

    # ✅ 解压 .tar.gz（所有平台通用）
    info "📦 解压 fs.tar.gz..."
    tar -xzf "fs.tar.gz" -C "$FS_ROOT" || error "解压失败"

    # ✅ 判断是否为 Windows 平台，决定二进制名
    local expected_binary="fs"
    if [ "$os" = "windows" ]; then
        expected_binary="fs.exe"
    fi

    # 检查解压后的二进制是否存在
    if [ ! -f "$FS_ROOT/$expected_binary" ]; then
        error "未找到二进制文件: $FS_ROOT/$expected_binary，请检查发布包内容"
    fi

    # 移动到 bin 目录（保持统一名称）
    mv "$FS_ROOT/$expected_binary" "$BIN_DIR/fs" || error "移动二进制失败"
    chmod +x "$BIN_DIR/fs"

    # 配置 shell
    if [ "$AUTO_CONFIG_SHELL" = true ]; then
        configure_shell
    fi

    # 清理
    cleanup

    # 完成
    success "${SOFTWARE_NAME} v${version} 安装成功！"
    echo
    echo "🎉 使用方式: fs --help"
    show_post_install_tip  # ← 核心优化：动态提示
    echo "📖 项目地址: https://gitee.com/liulei152/fs"
    echo

    # 验证
    if command -v fs >/dev/null 2>&1; then
        local ver
        ver=$(fs --version 2>&1 | head -n1)
        success "版本: $ver"
    else
        warn "建议重启终端或运行: source ~/.profile"
    fi
}

# -----------------------------
# 🏁 入口点
# -----------------------------

# 检测 'curl | bash'
if [ ! -t 0 ]; then
    info "检测到 'curl | bash' 安装模式"
fi

# 注册清理函数
trap cleanup EXIT

# 执行
main "$@"
