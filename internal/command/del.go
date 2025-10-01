/*
Copyright © 2025 SHEN YI <1245332635@qq.com>
*/
package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// delCmd represents the del command
var delCmd = &cobra.Command{
	Use:   "del <name>",
	Short: "删除一个 SSH 服务器配置",
	Long: `删除指定名称的主机配置文件。

示例：
  fs del web1
  fs del db-prod --force  # 跳过确认`,
	Args: cobra.ExactArgs(1), // 必须且只能传 1 个参数
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		force, _ := cmd.Flags().GetBool("force")

		if name == "" {
			fmt.Fprintln(os.Stderr, "❌ 错误：主机名称不能为空")
			os.Exit(1)
		}

		configFile := filepath.Join(fsHostsConfigDir, name+".yaml")

		// 检查文件是否存在
		_, err := os.Stat(configFile)
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "❌ 未找到配置文件: %s\n", configFile)
			os.Exit(1)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ 检查文件时出错: %v\n", err)
			os.Exit(1)
		}

		// 交互式确认（除非 --force）
		if !force {
			fmt.Printf("⚠️  确定要删除主机 '%s' 的配置吗？这将永久删除文件。\n", name)
			fmt.Print("   输入 'yes' 确认: ")
			var confirm string
			fmt.Scanln(&confirm)
			if strings.ToLower(confirm) != "yes" {
				fmt.Println("❌ 删除已取消。")
				return
			}
		}

		// 执行删除
		if err := os.Remove(configFile); err != nil {
			fmt.Fprintf(os.Stderr, "❌ 删除文件失败: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ 已删除主机配置: %s\n", name)

		// 如果删除的是当前选中的主机，清除 current 文件
		currentFile := getUsedConfigFile()
		currentData, err := os.ReadFile(currentFile)
		if err == nil && strings.TrimSpace(string(currentData)) == name {
			if err := os.Remove(currentFile); err != nil {
				fmt.Fprintf(os.Stderr, "⚠️ 无法清除当前选中项: %v\n", err)
				// 不退出，只是提示
			} else {
				fmt.Printf("ℹ️  当前选中项 '%s' 已清除。\n", name)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(delCmd)

	// 添加 --force 标志
	delCmd.Flags().BoolP("force", "f", false, "跳过确认提示")
}
