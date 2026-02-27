/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package command

import (
	"fmt"
	"fs/internal/utils"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

// sshCmd represents the ssh command
var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "SSH 连接到当前选中的服务器",
	Long: `SSH 连接到指定或当前选中的主机。
	
示例：
  fs ssh             # 使用当前选中的主机
  fs ssh web1        # 临时连接 web1，不影响当前选中`,
	Run: func(cmd *cobra.Command, args []string) {
		var configName string

		// 判断是否传入了主机名
		if len(args) > 0 {
			configName = args[0]
			fmt.Printf("🎯 正在连接指定主机: %s\n", configName)
			tempCfg, err := utils.LoadConfigByName(fsHostsConfigDir, configName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "❌ 加载配置失败: %v\n", err)
				os.Exit(1)
			}
			cfg = tempCfg

		} else {
			// 否则使用当前选中的主机
			configName = cfg.Name
			fmt.Printf("🎯 正在连接当前主机: %s\n", configName)
		}

		if cfg == nil {
			fmt.Fprintln(os.Stderr, "❌ 未加载任何主机配置，请先使用 'fs use <name>' 选择一个主机")
			os.Exit(1)
		}

		// 获取认证方式（密钥优先）
		authMethods, err := utils.GetAuthMethod(cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}

		// 1. 构建 SSH 配置
		sshConfig := &ssh.ClientConfig{
			User:            cfg.User,
			Auth:            authMethods,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(), // ⚠️ 测试用，生产环境应验证 HostKey
			Timeout:         10 * time.Second,
		}

		// 2. 构建地址
		address := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
		if cfg.Port == 0 {
			address = fmt.Sprintf("%s:%d", cfg.Host, 22) // 默认 22
		}

		// 3. 拨号连接
		client, err := ssh.Dial("tcp", address, sshConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ 连接失败: %v\n", err)
			os.Exit(1)
		}
		defer client.Close()
		fmt.Printf("✅ 已连接到 %s (%s@%s)\n", cfg.Name, cfg.User, cfg.Host)

		// 4. 创建会话
		session, err := client.NewSession()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ 创建会话失败: %v\n", err)
			os.Exit(1)
		}
		defer session.Close()

		// === 关键：直接绑定标准流，不要用 Pipe ===
		session.Stdin = os.Stdin
		session.Stdout = os.Stdout
		session.Stderr = os.Stderr

		// 设置本地终端为 raw 模式
		fd := int(os.Stdin.Fd())
		oldState, err := term.MakeRaw(fd)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ 设置 raw 模式失败: %v\n", err)
			os.Exit(1)
		}
		defer term.Restore(fd, oldState)

		// 获取终端尺寸
		width, height, _ := term.GetSize(fd)
		termType := os.Getenv("TERM")
		if termType == "" {
			termType = "xterm-256color"
		}
		modes := ssh.TerminalModes{
			ssh.ECHO:          1,
			ssh.TTY_OP_ISPEED: 14400,
			ssh.TTY_OP_OSPEED: 14400,
		}
		if err := session.RequestPty(termType, height, width, modes); err != nil {
			fmt.Fprintf(os.Stderr, "❌ 请求 PTY 失败: %v\n", err)
			os.Exit(1)
		}

		// 监听窗口大小变化
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, syscall.SIGWINCH)
		go func() {
			for range sigc {
				if w, h, err := term.GetSize(fd); err == nil {
					_ = session.WindowChange(h, w)
				}
			}
		}()

		// 启动交互式 shell（推荐用 Shell()）
		if err := session.Shell(); err != nil {
			fmt.Fprintf(os.Stderr, "❌ 启动 shell 失败: %v\n", err)
			os.Exit(1)
		}

		// 等待用户退出
		if err := session.Wait(); err != nil {
			// 可选：记录退出状态，但通常无需处理
		}

		fmt.Println("\n✅ 已断开连接")
	},
}

func init() {
	rootCmd.AddCommand(sshCmd)
}
