/*
Copyright © 2025 SHEN YI <1245332635@qq.com>
*/
package command

import (
	"fmt"
	"fs/internal/types"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "查看所有已配置的 SSH 服务器",
	Long: `列出所有已保存的 SSH 主机配置，并高亮当前选中项。

支持按 ID 过滤：
  fs show prod        # 显示 ID 包含 'prod' 的
  fs show web db      # 显示 ID 包含 'web' 或 'db' 的

示例：
  fs show
  fs show staging`,
	Args: cobra.ArbitraryArgs, // 支持 0 个或多个参数
	Run: func(cmd *cobra.Command, args []string) {
		// 获取所有关键词（不区分大小写）
		var filters []string
		for _, arg := range args {
			filters = append(filters, strings.ToLower(arg))
		}

		// 1. 获取所有 .yaml 配置文件
		files, err := filepath.Glob(filepath.Join(fsHostsConfigDir, "*.yaml"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ 配置目录访问失败: %v\n", err)
			os.Exit(1)
		}

		if len(files) == 0 {
			fmt.Println("📭 暂无任何主机配置，请使用 'fs add' 添加。")
			return
		}

		// 2. 存储所有配置，key 可以是文件名或 host_port
		configs := make(map[string]types.SshConfig)
		var keys []string

		for _, file := range files {
			data, err := os.ReadFile(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "⚠️ 跳过文件 %s: 读取失败: %v\n", file, err)
				continue
			}

			var config types.SshConfig
			if err := yaml.Unmarshal(data, &config); err != nil {
				fmt.Fprintf(os.Stderr, "⚠️ 跳过文件 %s: YAML 解析失败: %v\n", file, err)
				continue
			}

			// 使用文件名（不含扩展名）作为 key
			key := strings.TrimSuffix(filepath.Base(file), ".yaml")

			// 在遍历文件时：
			match := len(filters) == 0 // 无过滤条件 → 匹配
			for _, f := range filters {
				if strings.Contains(strings.ToLower(key), f) {
					match = true
					break
				}
			}
			if !match {
				continue
			}

			configs[key] = config
			keys = append(keys, key)
		}

		// 3. 排序输出（按 key 字典序）
		sort.Strings(keys)

		// 4. 读取当前选中的 ID（可选）
		var currentKey string
		configFile := getUsedConfigFile()
		currentKeyBytes, err := os.ReadFile(configFile)
		if err == nil && len(currentKeyBytes) > 0 {
			currentKey = string(currentKeyBytes)
		}

		// 5. 创建表格
		table := tablewriter.NewWriter(os.Stdout)
		table.SetAlignment(tablewriter.ALIGN_CENTER)
		header := []string{"Name", "Host", "Port", "User", "Private-Key-Path", "Default-Uploaddir", "Default-Downloaddir"}
		table.SetHeader(header)
		table.SetColWidth(30) // 防止内容被截断
		highlightStyle := make([]tablewriter.Colors, len(header))
		for i := range highlightStyle {
			highlightStyle[i] = tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiWhiteColor, tablewriter.BgHiBlueColor}
		}

		// 6. 填充数据
		for _, key := range keys {
			cfg := configs[key]
			row := []string{
				cfg.Name,
				cfg.Host,
				strconv.Itoa(cfg.Port),
				cfg.User,
				cfg.PrivateKeyPath,
				cfg.DefaultUploadDir,
				cfg.DefaultDownloadDir,
			}

			if key == currentKey {
				// 高亮当前选中行（红底白字加粗）
				table.Rich(row, highlightStyle)
			} else {
				table.Append(row)
			}
		}

		// 7. 输出表格
		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
