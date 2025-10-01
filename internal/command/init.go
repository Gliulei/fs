/*
Copyright © 2025 SHEN YI <1245332635@qq.com>
*/
package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化 fs 目录结构",
	Long:  `创建 fs 所需的目录：主目录、bin 目录、hosts 配置目录。`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("初始化fs开始...")

		// 确保目录存在
		dirs := []string{fsRootDir, fsBinDir, fsHostsConfigDir}
		for _, dir := range dirs {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "创建目录失败 %s: %v\n", dir, err)
				os.Exit(1)
			}
		}

		fmt.Printf("✅ %s\n", "初始化完成")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
