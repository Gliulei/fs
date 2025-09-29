# 📁 fs — 端到端高效文件传输工具

> ⚡ 轻量级命令行文件同步工具，让 `scp` 成为过去式

`fs` 是一款专为开发者与运维人员设计的轻量级文件传输工具，旨在替代传统 `scp` 命令，提供更简洁、高效、智能的跨主机文件操作体验。

无需记忆复杂参数，一键上传下载，支持多环境切换与 Shell 自动补全，大幅提升远程文件管理效率。

---

## ✨ 特性

- ✅ **简洁语法**：无需记忆复杂参数，一键完成本地与远程主机间的文件拷贝。
- ✅ **高效传输**：基于 SSH 协议优化数据流，支持压缩传输，显著提升速度。
- ✅ **智能补全**：支持 Bash / Zsh 自动补全，自动识别主机组和路径，减少输入错误。
- ✅ **跨平台支持**：兼容 Linux、macOS、WSL 及 Windows PowerShell。
- ✅ **多环境管理**：通过 `fs use <group>` 快速切换开发、测试、生产等服务器组。
- ✅ **操作可追溯**：自动记录命令历史，便于审计与复用。

---

## 📦 安装

### 自动化安装（推荐）

#### Linux / macOS（Bash / Zsh）

```bash
# 建议安装前清空可能冲突的环境变量
unset GOROOT GOBIN GOPATH

# 下载并执行安装脚本
curl -sSL https://gitee.com/liulei152/fs/raw/master/install.sh | bash

# 可选：避免别名冲突（如已存在 'fs' 别名）
cat << 'EOF' >> ~/.bashrc
if [[ -n $(alias fs 2>/dev/null) ]]; then
    unalias fs
fi
EOF

# 加载 fs 环境变量
source "$HOME/.fs/env"
```

#### Windows（PowerShell）

```powershell
iwr https://gitee.com/liulei152/fs/raw/master/install.ps1 -UseBasicParsing | iex
```

### 手动安装

```bash
tar -xzf fs-linux-amd64.tar.gz
sudo mv fs /usr/local/bin/
```

## 🚀 快速开始

只需几个命令，即可上手使用。

### 1. 初始化配置

```bash
fs init
```
### 2. 文件传输

上传文件到远程主机：
```bash
fs upload filename
```

下载文件到本地：

```bash
fs download filename
```
### 3. 切换主机组（分组）

支持多环境快速切换：

```bash
fs use group
```

📌 示例：

```
fs use production   # 切换到生产环境
fs use staging      # 切换到测试环境
```

---



## 🧭 路线图（ROADMAP）

+ [x] 上传文件加个在当前目录下find
+ [x] 记录命令历史
+ [x] 配置增加
+ [x] SHELL Zsh|Bash 自动补全
+ [ ] 实现类似scp /home/space/music/ root@www.runoob.com:/home/root/others/功能，自动记忆
+ [ ] 配置删除
+ [ ] 多文件下载|上传
