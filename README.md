<div align="center">

**English** | [中文](./README_zh.md)

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.24-blue.svg)](https://golang.org/)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)](https://github.com/yourusername/fs/actions)

</div>

---

# 📁 fs: The Next-Gen SSH Management & File Sync Tool

**Simpler, Smarter, and More Efficient Remote Operations**

In today’s fast-paced development and DevOps environments, SSH connections and file transfers are daily essentials. Yet, traditional tools like `scp` and manual `~/.ssh/config` management have become increasingly cumbersome and inefficient. Complex configurations, cryptic command syntax, and lack of centralized control — these pain points are now being redefined by **fs**, a next-generation command-line tool designed to revolutionize how we interact with remote servers.

**fs** is a powerful, all-in-one remote operations toolkit built for developers and system administrators. It integrates SSH management, file synchronization, environment switching, and intelligent auto-completion into a single cohesive workflow — reimagining remote server management for the modern era.

No more juggling configurations or memorizing command flags. With fs, managing SSH servers becomes effortless. Add, configure, and connect to remote hosts with simple commands. Switch between environments, transfer files seamlessly, and enjoy shell auto-completion — all designed to dramatically boost your productivity.

## 🌟 Why Do You Need fs?

Have you ever faced these scenarios?

- Need to send a file but have to look up the IP, username, and port, then piece together a messy `scp` command?
- Juggle between multiple environments (dev, staging, production), only to accidentally connect to the wrong server?
- Set up a new machine and spend hours re-creating your SSH configurations?
- Type commands blindly without auto-completion, leading to typos and debugging nightmares?

**fs is built to solve exactly these problems.**

More than just a tool, fs offers a modern, systematic approach to remote operations. Through intuitive commands, smart interactions, and robust feature integration, fs significantly improves efficiency and accuracy across multi-host workflows.

## ✨ Key Features

### 🔧 1. Simplified SSH Server Management

No more manual config edits. Register a server with just one command:

```
fs add my-server --host=192.168.1.100 --user=dev --port=22
```

Supports naming, grouping, and aliases — easily manage dozens or even hundreds of servers.

### 📤 2. One-Click File Transfer

Say goodbye to complicated `scp` syntax. Use the `fs upload` command just like local copying:

```
fs upload ./app.zip my-server:/opt/deploy/
```

Supports bidirectional transfer — local ↔ remote has never been easier.

### ⚡ 3. Optimized Transfer Performance

Leverages SSH protocol optimizations with built-in compression to accelerate file transfers, especially beneficial for large files or bulk syncing.

### 🧠 4. Intelligent Shell Auto-Completion

Supports auto-completion for Bash and Zsh. Press `Tab` after `fs cp` to get real-time suggestions for server names, groups, and paths — minimizing errors and boosting speed.

### 🌐 5. Instant Environment Switching

Switch contexts across environments using:

```
fs use production
```

Immediately enter production mode with confidence, reducing the risk of accidental operations — ideal for multi-cluster, multi-role operations.

### 💻 6. Cross-Platform Compatibility

Fully supports Linux, macOS, WSL, and Windows PowerShell. No matter your OS, fs integrates smoothly into your workflow.

### 🔗 7. Direct SSH Access

Connect directly via SSH with:

```
fs ssh <server>
```

Start an interactive session instantly — no extra steps required.

### 📚 8. Audit-Ready Operation History

Automatically logs frequently used commands and operation history. Easy to review, reuse, and audit — great for debugging and team knowledge sharing.

### 📦 9. Configurable Export & Import

Export your entire server configuration as an encrypted package, or import configurations from others. Simplify onboarding, environment migration, and team collaboration.

## 🚀 Real-World Use Cases

- **Development & Deployment**: Quickly sync build artifacts to test servers.
- **Operations & Monitoring**: Switch to production group and run health checks in bulk.
- **Team Collaboration**: Share standardized environment packages so new members can be up and running in minutes.
- **Multi-Cloud Management**: Centrally manage servers across different cloud providers.

## 📦 Installation

### Automated Installation (Recommended)

#### Linux / macOS (Bash / Zsh)

```
# Clear conflicting environment variables if needed
unset GOROOT GOBIN GOPATH

# Download and run the install script
curl -sSL https://gitee.com/liulei152/fs/raw/master/install.sh | bash

# Optional: Avoid alias conflicts (e.g., if 'fs' is already aliased)
cat << 'EOF' >> ~/.bashrc
if [[ -n $(alias fs 2>/dev/null) ]]; then
    unalias fs
fi
EOF

# Load fs environment variables
source "$HOME/.fs/env"
```

#### Windows (PowerShell)

```
iwr https://gitee.com/liulei152/fs/raw/master/install.ps1 -UseBasicParsing | iex
```

### Manual Installation

```
tar -xzf fs-linux-amd64.tar.gz
sudo mv fs /usr/local/bin/
```

## 🚀 Quick Start

Get started with a few simple commands for server management and file operations.

### 1. Initialize Configuration

```
fs init
```

### 2. Manage SSH Servers

```
# Add a new SSH server
fs add --name web1 --user alice --host 192.168.1.100 --port 22

# List all configured servers
fs show

# Switch default server
fs use web1

# Connect via SSH
fs ssh
```

### 3. File Transfer

```
# Upload file to remote server
fs upload filename

# Download file from remote server
fs download filename
```

### 4. Configuration Management

```
# Archive all configs for backup or team sharing
fs archive

# Archive to a specific file
fs archive --output my-configs.tar.gz

# Preview contents of an archive
fs import config-archive.tar.gz --preview

# Import configs (interactive conflict handling)
fs import config-archive.tar.gz

# Force override during import
fs import config-archive.tar.gz --force

# Skip existing configurations
fs import config-archive.tar.gz --skip-existing
```

## 🧭 Roadmap

- Add `find` in current directory for file upload
- Command history logging
- Configuration addition
- Bash/Zsh shell auto-completion
- Configuration deletion
- Multi-file upload/download support

## 📞 Contact Us

For questions, feedback, or collaboration, feel free to reach out:

- **Name**: Shenyi
- **Email**: 1245332635@qq.com
- **GitHub**: https://github.com/Gliulei/

## 🤝 Contributions Welcome

We welcome bug reports, feature suggestions, and pull requests!

**How to contribute:**

1. Fork this repository
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the **MIT License** — see the LICENSE file for details.

**fs — Making remote operations simple again.**
📅 **Today is: February 27, 2026**

---