/*
Copyright © 2025 SHEN YI <1245332635@qq.com>
*/

package command

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/ngaut/log"
	"github.com/spf13/cobra"
)

const (
	historyFileName = "history.txt"
)

var (
	tailLines int
	clearHist bool
	rawOutput bool
)

// historyCmd represents the history command
var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "查看命令执行历史",
	Long: `显示所有已执行的 fs 命令历史记录。

支持：
  - 查看最近 N 条记录
  - 清空历史
  - 原始命令输出（便于脚本处理）`,
	Example: `
  fs history                # 查看全部历史
  fs history -n 10          # 查看最近 10 条
  fs history --clear        # 清空历史
  fs history --raw          # 仅输出命令，无编号`,
	Run: func(cmd *cobra.Command, args []string) {
		historyFile := GetHistoryFilePath()

		// 处理 --clear
		if clearHist {
			err := os.Remove(historyFile)
			if err != nil && !os.IsNotExist(err) {
				cmd.PrintErrf("❌ 清空历史失败: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("✅ 历史记录已清空")
			return
		}

		// 打开历史文件
		f, err := os.Open(historyFile)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("📭 历史记录为空")
			} else {
				cmd.PrintErrf("❌ 无法读取历史文件: %v\n", err)
			}
			return
		}
		defer f.Close()

		// 读取所有行
		var lines []string
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			text := scanner.Text()
			if text != "" { // 忽略空行
				lines = append(lines, text)
			}
		}

		if err := scanner.Err(); err != nil && err != io.EOF {
			cmd.PrintErrf("❌ 读取历史时出错: %v\n", err)
			os.Exit(1)
		}

		if len(lines) == 0 {
			fmt.Println("📭 历史记录为空")
			return
		}

		// 只显示最后 N 行
		if tailLines > 0 && tailLines < len(lines) {
			lines = lines[len(lines)-tailLines:]
		}

		// 输出
		for i, line := range lines {
			if rawOutput {
				fmt.Println(line)
			} else {
				fmt.Printf("[%d] %s\n", i+1, line)
			}
		}

		if !rawOutput {
			fmt.Printf("\n📌 共 %d 条记录\n", len(lines))
		}
	},
}

func init() {
	rootCmd.AddCommand(historyCmd)

	// 添加标志
	historyCmd.Flags().IntVarP(&tailLines, "lines", "n", 0, "只显示最近 N 条记录")
	historyCmd.Flags().BoolVar(&clearHist, "clear", false, "清空历史记录")
	historyCmd.Flags().BoolVar(&rawOutput, "raw", false, "仅输出原始命令（无编号）")
}

// GetHistoryFile 返回历史文件路径，并确保目录存在
func GetHistoryFilePath() string {
	histFile := path.Join(fsRootDir, historyFileName)

	return histFile
}

// record 记录命令到历史文件
func record(cmdLog []string) {
	cmdString := strings.Join(cmdLog, " ") + "\n"
	historyFile := GetHistoryFilePath()

	f, err := os.OpenFile(historyFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Errorf("无法打开历史文件: %s: %v", historyFile, err)
		return
	}
	defer f.Close() // ✅ 重要：确保关闭文件

	if _, err := f.WriteString(cmdString); err != nil {
		log.Errorf("写入历史文件失败: %v", err)
	}
}
