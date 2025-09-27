#!/bin/bash
# ==============================================================================
# 一键安装脚本 - 用于通过 curl | bash 方式安装软件
# 示例用法: curl -fsSL https://example.com/install | bash
# 
# 作者: 通义千问
# 注意: 请将下方的占位符替换为你的实际软件信息。
# ==============================================================================

# -----------------------------
# 🔧 配置区域 - 必须修改!
# -----------------------------

# 软件名称 (用于显示)
SOFTWARE_NAME="MyAwesomeApp"

# 软件的下载URL (请替换为你的实际下载链接)
SOFTWARE_DOWNLOAD_URL="https://example.com/releases/myapp-latest-linux-x64.tar.gz"

# 下载后的文件名
ARCHIVE_NAME="myapp.tar.gz"

# 解压后的主目录名
EXTRACTED_DIR="myapp"

# 安装的目标路径 (通常为 /usr/local 或 /opt)
INSTALL_PREFIX="/usr/local"

# 软件二进制文件的相对路径 (相对于 EXTRACTED_DIR)
BINARY_PATH="bin/myapp"

# -----------------------------
# 🚀 脚本逻辑
# -----------------------------

# 设置严格模式: 遇错退出, 未定义变量报错, 管道错误被捕获
set -euo pipefail

# 检查是否在管道中运行 (curl | bash)
if [ ! -t 0 ]; then
    echo "💡 检测到通过 'curl | bash' 方式安装，开始执行..."
fi

# 函数: 打印彩色信息
info() {
    echo -e "\033[1;34mINFO:\033[0m $*"
}

success() {
    echo -e "\033[1;32mSUCCESS:\033[0m $*"
}

error() {
    echo -e "\033[1;31mERROR:\033[0m $*" >&2
    exit 1
}

# 检查必要工具
for cmd in curl tar; do
    if ! command -v $cmd &> /dev/null; then
        error "缺少必要命令 '$cmd'，请先安装。"
    fi
done

# 创建临时目录
TMP_DIR=$(mktemp -d)
info "创建临时目录: $TMP_DIR"
cd "$TMP_DIR"

# 下载软件
info "正在下载 $SOFTWARE_NAME..."
curl -fL -o "$ARCHIVE_NAME" "$SOFTWARE_DOWNLOAD_URL" || error "下载失败"

success "下载成功!"

# 解压
info "正在解压..."
tar -xzf "$ARCHIVE_NAME" || error "解压失败"

# 安装到目标位置
TARGET_DIR="$INSTALL_PREFIX/$EXTRACTED_DIR"
info "正在安装到 $TARGET_DIR..."
sudo rm -rf "$TARGET_DIR" 2>/dev/null || true
sudo mkdir -p "$INSTALL_PREFIX"
sudo mv "$EXTRACTED_DIR" "$TARGET_DIR"

# 创建全局符号链接
BIN_SOURCE="$TARGET_DIR/$BINARY_PATH"
BIN_TARGET="/usr/local/bin/$(basename "$BINARY_PATH")"
info "创建全局命令链接: $BIN_TARGET"
sudo ln -sf "$BIN_SOURCE" "$BIN_TARGET"

# 清理
info "清理临时文件..."
cd /
sudo rm -rf "$TMP_DIR"

# 完成
success "$SOFTWARE_NAME 安装成功!"
echo ""
echo "🎉 你现在可以运行: $(basename "$BINARY_PATH")"
echo "📖 查看文档: https://example.com/docs"
echo ""

# 可选: 验证安装
if command -v "$(basename "$BINARY_PATH")" &> /dev/null; then
    version=$("$(basename "$BINARY_PATH")" --version 2>&1 | head -n1)
    echo "✅ 版本: $version"
else
    echo "⚠️  如果命令未生效，请尝试重新打开终端或运行: hash -r"
fi