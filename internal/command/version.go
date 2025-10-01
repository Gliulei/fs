/*
Copyright © 2025 SHEN YI <1245332635@qq.com>
*/
package command

import (
	"fmt"
	"fs/internal/version"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "打印版本信息",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("fs version %s\n", version.Version)
		fmt.Printf("  Commit: %s\n", version.Commit)
		fmt.Printf("  Built:  %s\n", version.Date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
