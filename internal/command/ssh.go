/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package command

import (
	"fmt"
	"fs/internal/utils"
	"io"
	"os"
	"sync"
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

		// 5. 将标准输入/输出/错误绑定到远程会话
		stdin, err := session.StdinPipe()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ 无法获取 stdin 管道: %v\n", err)
			os.Exit(1)
		}

		stdout, err := session.StdoutPipe()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ 无法获取 stdout 管道: %v\n", err)
			os.Exit(1)
		}

		stderr, err := session.StderrPipe()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ 无法获取 stderr 管道: %v\n", err)
			os.Exit(1)
		}

		// 6. 请求伪终端（PTY）
		modes := ssh.TerminalModes{
			ssh.ECHO:          1,     // 启用回显
			ssh.TTY_OP_ISPEED: 14400, // 输入速率
			ssh.TTY_OP_OSPEED: 14400, // 输出速率
		}

		// 获取本地终端尺寸
		fd := int(os.Stdin.Fd())
		width, height, _ := term.GetSize(fd)

		err = session.RequestPty("xterm", height, width, modes)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ 请求 PTY 失败: %v\n", err)
			os.Exit(1)
		}

		// 7. 开始远程会话
		err = session.Start("exec /bin/bash -l")
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ 启动远程 shell 失败: %v\n", err)
			os.Exit(1)
		}

		// 8. 并发复制数据流
		var wg sync.WaitGroup
		wg.Add(3)

		// 将远程输出打印到本地 stdout
		go func() {
			defer wg.Done()
			io.Copy(os.Stdout, stdout)
		}()

		// 将远程错误输出到本地 stderr
		go func() {
			defer wg.Done()
			io.Copy(os.Stderr, stderr)
		}()

		// 将本地输入发送到远程 stdin
		go func() {
			defer wg.Done()
			io.Copy(stdin, os.Stdin)
		}()

		// 9. 等待会话结束
		wg.Wait()
		fmt.Println("\n✅ 已断开连接")
	},
}

func init() {
	rootCmd.AddCommand(sshCmd)
}
