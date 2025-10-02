/*
Copyright © 2025 SHEN YI <1245332635@qq.com>
*/

package command

import (
	"fmt"
	"fs/internal/types"
	"os"
	"path"

	"github.com/ngaut/log"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "更新一个已存在的 SSH 服务器配置",
	Long: `更新本地配置文件中的 SSH 服务器信息。

支持只更新部分字段，未提供的字段将保留原值。

示例：
  fs update web1 -u alice --port 2222
  fs update db-prod --tags database,production --upload-dir /opt/backups
  fs update legacy --password newpass --private-key-path ~/.ssh/new_key
`,
	Args: cobra.ExactArgs(1), // 必须传一个参数：配置名
	Run: func(cmd *cobra.Command, args []string) {
		// 从 args 获取要更新的配置名
		name = args[0]

		filename := fmt.Sprintf("%s.yaml", name)
		filepath := path.Join(fsHostsConfigDir, filename)

		// ✅ 检查文件是否存在
		if _, err := os.Stat(filepath); os.IsNotExist(err) {
			log.Fatalf("❌ 配置不存在: %s\n请使用 'fs list' 查看可用配置", name)
		}

		// ✅ 读取现有配置
		data, err := os.ReadFile(filepath)
		if err != nil {
			log.Fatalf("❌ 读取配置文件失败 %s: %v", filepath, err)
		}

		var config types.SshConfig
		if err := yaml.Unmarshal(data, &config); err != nil {
			log.Fatalf("❌ 解析 YAML 失败: %v", err)
		}

		// ✅ 逐字段更新（只更新用户指定的字段）
		updated := false

		if cmd.Flags().Changed("user") {
			config.User = user
			updated = true
		}
		if cmd.Flags().Changed("password") {
			config.Password = password
			updated = true
		}
		if cmd.Flags().Changed("host") {
			config.Host = host
			updated = true
		}
		if cmd.Flags().Changed("port") {
			config.Port = port
			updated = true
		}
		if cmd.Flags().Changed("upload-dir") {
			config.DefaultUploadDir = uploadDir
			updated = true
		}
		if cmd.Flags().Changed("download-dir") {
			config.DefaultDownloadDir = downloadDir
			updated = true
		}
		if cmd.Flags().Changed("private-key-path") {
			config.PrivateKeyPath = privateKeyPath
			updated = true
		}
		if cmd.Flags().Changed("passphrase") {
			config.Passphrase = passphrase
			updated = true
		}
		if cmd.Flags().Changed("tags") {
			config.Tags = tags
			updated = true
		}

		// ✅ 如果没有任何字段被修改，提示用户
		if !updated {
			fmt.Printf("🟡 未提供任何更新字段，配置 %s 保持不变。\n", name)
			return
		}

		// ✅ 重新序列化为 YAML
		newData, err := yaml.Marshal(&config)
		if err != nil {
			log.Fatalf("❌ 序列化配置失败: %v", err)
		}

		// ✅ 写回文件
		if err := os.WriteFile(filepath, newData, 0600); err != nil {
			log.Fatalf("❌ 写入更新后的配置失败 %s: %v", filepath, err)
		}

		fmt.Printf("✅ 已更新服务器配置: %s → %s\n", name, filepath)
	},
}

func init() {
	applySSHFlags(updateCmd)

	rootCmd.AddCommand(updateCmd)
}
