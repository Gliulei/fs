/*
Copyright © 2025 SHEN YI <1245332635@qq.com>
*/
package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use [name]",
	Short: "设置当前使用的主机",
	Long: `设置当前使用的主机，后续命令将作用于该目标。

示例：
  fs use web1
  fs use prod-group`,
	Args: cobra.ExactArgs(1), // 确保必须传一个参数
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		if name == "" {
			fmt.Fprintln(os.Stderr, "❌ 错误：名称不能为空")
			os.Exit(1)
		}

		configFile := getUsedConfigFile()

		// 确保目录存在
		if err := os.MkdirAll(fsRootDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "❌ 无法创建配置目录 %s: %v\n", fsRootDir, err)
			os.Exit(1)
		}

		// 写入当前选中的名称（覆盖模式）
		err := os.WriteFile(configFile, []byte(name), 0600) // ✅ 仅用户可读写
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ 写入配置失败 %s: %v\n", configFile, err)
			os.Exit(1)
		}

		fmt.Printf("✅ 当前已切换到: %s\n", name)
	},
}

func init() {
	rootCmd.AddCommand(useCmd)
}
