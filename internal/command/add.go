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
	"gopkg.in/yaml.v3" // 注意是 v3
)

// 定义变量（init 之前）
var (
	user           string
	password       string
	host           string
	port           int
	uploadDir      string
	downloadDir    string
	name           string
	tags           []string
	privateKeyPath string
	passphrase     string
)

// applySSHFlags 将 SSH 相关的 flags 绑定到指定的 Command 上
// 可被 add、update 等命令复用
func applySSHFlags(cmd *cobra.Command) {
	flags := cmd.Flags()

	flags.StringVarP(&user, "user", "u", "", "用户名")
	flags.StringVarP(&password, "password", "p", "", "密码（可选）")
	flags.StringVarP(&host, "host", "H", "", "服务器地址")
	flags.IntVarP(&port, "port", "P", 22, "SSH 端口号（默认 22）")
	flags.StringVarP(&privateKeyPath, "private-key-path", "", "", "私钥路径（可选）")
	flags.StringVarP(&passphrase, "passphrase", "", "", "私钥密码（可选）")
	flags.StringVarP(&name, "name", "n", "", "服务器名称（用于快速引用）")
	flags.StringSliceVarP(&tags, "tags", "t", nil, "服务器标签，用空格分隔，如: web db prod")
	flags.StringVar(&uploadDir, "upload-dir", "", "指定默认上传文件目录")
	flags.StringVar(&downloadDir, "download-dir", "", "指定默认下载文件目录")
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "添加一个 SSH 服务器配置",
	Long: `将一个新的 SSH 服务器配置添加到本地配置文件中。

	示例：
	  fs add -n web1 -u alice -H 192.168.1.100 -P 22 -p 123456
	  fs add --name db-prod --user root --host 10.0.0.5 --password secret --upload-dir /tmp/uploads --download-dir /tmp/downloads --private-key-path /home/alice/.ssh/id_rsa`,
	Run: func(cmd *cobra.Command, args []string) {
		if name == "" {
			name = host
		}

		config := types.SshConfig{
			User:               user,
			Password:           password,
			Host:               host,
			Port:               port,
			DefaultUploadDir:   uploadDir,
			DefaultDownloadDir: downloadDir,
			Name:               name,
			Tags:               tags,
			PrivateKeyPath:     privateKeyPath,
			Passphrase:         passphrase,
		}

		// 🔐 用 name 命名文件
		filename := fmt.Sprintf("%s.yaml", name)
		filepath := path.Join(fsHostsConfigDir, filename)

		// ✅ 确保目录存在
		if err := os.MkdirAll(fsHostsConfigDir, 0755); err != nil {
			log.Fatalf("无法创建配置目录 %s: %v", fsHostsConfigDir, err)
		}

		// ✅ 将结构体序列化为 YAML
		data, err := yaml.Marshal(&config)
		if err != nil {
			log.Fatalf("序列化配置失败: %v", err)
		}

		// ✅ 写入文件（覆盖模式）
		if err := os.WriteFile(filepath, data, 0600); err != nil {
			log.Fatalf("写入文件失败 %s: %v", filepath, err)
		}

		// ✅ 成功提示
		displayName := name
		if displayName == "" {
			displayName = fmt.Sprintf("%s:%d", host, port)
		}
		fmt.Printf("✅ 已添加服务器: %s → %s", displayName, filepath)
	},
}

func init() {
	applySSHFlags(addCmd)

	// 标记必填字段
	_ = addCmd.MarkFlagRequired("user")
	_ = addCmd.MarkFlagRequired("name")
	_ = addCmd.MarkFlagRequired("host")

	rootCmd.AddCommand(addCmd)
}
