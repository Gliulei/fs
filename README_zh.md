<div align="center">

[English](./README_en.md) | **中文**

[![许可证](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go版本](https://img.shields.io/badge/go-%3E%3D1.24-blue.svg)](https://golang.org/)
[![构建状态](https://img.shields.io/badge/build-passing-brightgreen.svg)](https://github.com/yourusername/fs/actions)

</div>


# 📁 fs：下一代 SSH 管理与文件同步利器

**让远程操作更简单、更智能、更高效**

在当今快节奏的开发与运维场景中，SSH 连接和文件传输是每日高频操作。然而，传统的 `scp` 命令和手动管理 `~/.ssh/config` 的方式，早已显得笨重而低效。配置繁琐、参数复杂、缺乏统一管理——这些痛点，正在被一款名为 **fs** 的新一代命令行工具彻底改变。

**fs** 是一款专为开发者与运维工程师打造的综合性远程操作利器，集 SSH 管理、文件同步、环境切换与智能补全于一体，旨在重塑远程工作流，让繁琐的服务器管理成为过去式。

## 🌟 为什么需要 fs？

你是否曾遇到以下场景？

- 想传个文件，却要翻查 IP、用户名、端口，再拼出一长串 `scp` 命令？
- 多套环境（开发、测试、生产）来回切换，配置混乱，容易连错服务器？
- 在不同机器间迁移工作环境时，不得不手动重建所有 SSH 配置？
- Shell 输入时没有补全，拼错主机名或路径，调试半小时才发现问题？

**fs 正是为解决这些问题而生。**

它不仅是一个工具，更是一种现代化、系统化的远程操作解决方案。通过简洁的命令、智能的交互与强大的功能集成，fs 显著提升了跨主机操作的效率与准确性。

---

## ✨ 核心亮点

### 🔧 1. 极简 SSH 服务器管理

无需手动编辑配置文件。只需一条命令：

```
fs add my-server --host=192.168.1.100 --user=dev --port=22
```

即可完成服务器注册。支持命名、分组、别名设置，轻松管理数十台甚至上百台主机。

### 📤 2. 一键文件传输

告别复杂的 `scp` 语法。使用 `fs cp` 命令，像本地拷贝一样简单：

```
fs cp ./app.zip my-server:/opt/deploy/
```

支持双向传输，本地 ↔ 远程一键直达。

### ⚡ 3. 高效传输优化

基于 SSH 协议深度优化数据流，内置压缩机制，大幅提升文件传输速度，尤其适合大文件或批量同步场景。

### 🧠 4. 智能 Shell 补全

支持 Bash 与 Zsh 自动补全。输入 `fs cp` 后按 Tab，即可自动提示服务器名、主机组、路径等选项，极大减少输入错误，提升操作流畅度。

### 🌐 5. 多环境快速切换

通过 `fs use <group>` 快速切换当前上下文环境：

```
fs use production
```

立即进入生产环境操作模式，避免误操作风险，特别适合多集群、多角色运维场景。

### 💻 6. 全平台兼容

完美支持 Linux、macOS、WSL 以及 Windows PowerShell，无论你使用何种系统，都能无缝接入 fs 的高效生态。

### 🔗 7. 直连 SSH 会话

无需额外命令，直接通过 `fs ssh <server>` 建立交互式连接，即刻进入远程终端，操作自然流畅。

### 📚 8. 操作可追溯

自动记录常用命令与操作历史，支持快速回溯与复用，便于审计、调试与团队知识沉淀。

### 📦 9. 配置打包与共享

支持将整个服务器配置导出为加密包，或批量导入他人分享的配置。新成员入职、环境迁移、团队协作从此变得轻而易举。

## 🚀 使用场景示例

- **开发部署**：快速将本地构建产物同步到测试服务器。
- **运维巡检**：一键切换至生产组，批量执行健康检查脚本。
- **团队协作**：导出“标准环境配置包”，新人一天内完成环境搭建。
- **跨区域管理**：统一管理分布在不同云厂商的服务器资源。

---

## 📦安装

### 自动化安装（推荐）

#### Linux / macOS（Bash / Zsh）

```bash
# 建议安装前清空可能冲突的环境变量
unset GOROOT GOBIN GOPATH

# 下载并执行安装脚本
curl -sSL https://gitee.com/liulei152/fs/raw/master/install.sh | bash

#可选：避免别名冲突（如已存在 'fs'别）
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

###手动安装

```bash
tar -xzf fs-linux-amd64.tar.gz
sudo mv fs /usr/local/bin/
```

##🚀 快速开始

通过几个简单命令即可开始服务器管理和文件操作。

### 1. 初始化配置

```bash
fs init
```

### 2. SSH服务器管理

```bash
# 添加新的SSH服务器
fs add --name web1 --user alice --host 192.168.1.100 --port 22

# 查看所有配置的服务器
fs show

#切换默认服务器
fs use web1

# 通过SSH直接连接
fs ssh
```

### 3. 文件传输操作

```bash
# 上传文件到远程服务器
fs upload filename

# 从远程服务器下载文件
fs download filename
```

### 4. 配置管理

```bash
#打包所有配置用于备份/团队分享
fs archive

#打包到指定文件
fs archive --output my-configs.tar.gz

#打包内容
fs import config-archive.tar.gz --preview

#导入配置（交互式处理冲突）
fs import config-archive.tar.gz

#强制覆盖导入
fs import config-archive.tar.gz --force

#跳过已存在配置
fs import config-archive.tar.gz --skip-existing
```

---

## 🧭路线图

- [x] 上传文件加个在当前目录下find
- [x]记录命令历史
- [x] 配置增加
- [x] SHELL Zsh|Bash 自动补全
- [x] 配置删除
- [ ] 多文件下载|上传

## 📞联方式

如有问题或反馈，请联系：

- **姓名**: shenyi
- **邮箱**: 1245332635@qq.com
- **GitHub**: [您的 GitHub 个人主页](https://github.com/Gliulei/)

## 贡献

欢迎提交 Issue 和 Pull Request！

1. Fork 本仓库
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 📄许证

本项目采用 MIT许证 - 详情请见 [LICENSE](LICENSE) 文件。
