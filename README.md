# 📁 fs - File Sync Tool

<div align="center">

**English** | [中文](./README_zh.md)

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.24-blue.svg)](https://golang.org/)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)](https://github.com/yourusername/fs/actions)

</div>

## Introduction

> ⚡ Lightweight command-line tool for SSH server management and file synchronization that makes `scp` and manual SSH management obsolete

`fs` is a comprehensive command-line tool designed for developers and DevOps engineers, providing an integrated solution for SSH server management and file transfer operations. It aims to replace traditional `scp` commands and manual SSH management with a more concise, efficient, and intelligent cross-host operations experience.

Forget complex parameters and manual server management - add, configure, and connect to SSH servers with simple commands, support multi-environment switching, file upload/download, and shell auto-completion, significantly improving remote server management efficiency.

---

##✨ Features

- ✅ **SSH Server Management**: Add, configure, and manage multiple SSH server connections with simple commands
- ✅ **One-click File Transfer**: No need to remember complex parameters, complete local-to-remote file copying with one click
- ✅ **Efficient Transfer**: SSH protocol optimized data stream, supports compression transfer, significantly improves speed
- ✅ **Smart Completion**: Supports Bash/Zsh auto-completion, automatically recognizes host groups and paths, reduces input errors
-✅ **Cross-platform Support**: Compatible with Linux, macOS, WSL, and Windows PowerShell
- ✅ **Multi-environment Management**: Quickly switch between development, testing, and production server groups via `fs use <group>`
- ✅ **Direct SSH Connection**: Connect directly to configured servers with `fs ssh`
- ✅ **Operation Traceability**: Automatically records command history for auditing and reuse
- ✅ **Configuration Packaging & Import**: Supports packaging and batch import of configuration files for environment migration and team collaboration

---

## 📦 Installation

### Automated Installation (Recommended)

#### Linux / macOS (Bash / Zsh)

```bash
# It's recommended to clear potentially conflicting environment variables before installation
unset GOROOT GOBIN GOPATH

# Download and execute installation script
curl -sSL https://gitee.com/liulei152/fs/raw/master/install.sh | bash

# Optional: Avoid alias conflicts (if 'fs' alias already exists)
cat << 'EOF' >> ~/.bashrc
if [[ -n $(alias fs 2>/dev/null) ]]; then
    unalias fs
fi
EOF

# Load fs environment variables
source "$HOME/.fs/env"
```

#### Windows (PowerShell)

```powershell
iwr https://gitee.com/liulei152/fs/raw/master/install.ps1 -UseBasicParsing | iex
```

### Manual Installation

```bash
tar -xzf fs-linux-amd64.tar.gz
sudo mv fs /usr/local/bin/
```

## 🚀 Quick Start

Get started with server management and file operations in just a few commands.

### 1. Initialize Configuration

```bash
fs init
```

### 2. SSH Server Management

```bash
# Add a new SSH server
fs add --name web1 --user alice --host 192.168.1.100 --port 22

# View all configured servers
fs show

# Switch default server
fs use web1

# Connect directly via SSH
fs ssh
```

### 3. File Transfer Operations

```bash
# Upload files to remote server
fs upload filename

# Download files from remote server
fs download filename
```

### 4. Configuration Management

```bash
# Package all configurations for backup/team sharing
fs archive

# Package to specified file
fs archive --output my-configs.tar.gz

# Preview package contents
fs import config-archive.tar.gz --preview

# Import configuration (interactive conflict handling)
fs import config-archive.tar.gz

# Force overwrite import
fs import config-archive.tar.gz --force

# Skip existing configurations
fs import config-archive.tar.gz --skip-existing
```

---

##🧭 Roadmap

- [x] Add find functionality in current directory for file upload
- [x] Record command history
- [x] Configuration addition
- [x] SHELL Zsh/Bash auto-completion
- [x] Configuration deletion
- [ ] Multi-file download/upload

## 📞 Contact

For questions or feedback, please contact:

- **Name**: shenyi
- **Email**: 1245332635@qq.com
- **GitHub**: [Your GitHub Profile](https://github.com/yourusername)

##🤝uting

Contributions are welcome! Please feel free to submit pull requests or open issues.

##📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.