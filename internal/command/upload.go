/*
Copyright © 2025 SHEN YI <1245332635@qq.com>
*/

package command

import (
	"context"
	"fmt"
	"fs/internal/utils"
	"os"
	"path"

	"github.com/bramvdbogaerde/go-scp"
	"github.com/ngaut/log"
	"github.com/spf13/cobra"
)

var (
	dryRun bool
	chmod  string
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload [local_file...] [remote_dir|remote_file]",
	Short: "上传文件到远程服务器",
	Long: `通过 SCP 上传本地文件到远程服务器。

支持：
  - 批量上传多个文件
  - 密钥认证
  - 进度条显示
  - --dry-run 模式预演

示例：
  fs upload myfile.txt                    # 上传到默认上传目录
  fs upload f1.txt f2.txt /tmp/           # 批量上传到 /tmp/
  fs upload *.log backup/                 # 通配符上传
  fs upload myfile.txt /tmp/data.txt      # 指定远程文件名
  fs upload -n myfile.txt /tmp/           # 预演操作，不实际上传`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.PrintErrln("❌ 未指定要上传的文件")
			cmd.Usage()
			os.Exit(1)
		}

		// 解析标志
		if chmod == "" {
			chmod = "0644"
		}

		// 1. 分离本地文件和远程路径
		var srcPaths []string
		var remotePath string

		if len(args) == 1 {
			// fs upload file.txt
			srcPaths = args[:1]
			remotePath = "" // 使用默认上传目录
		} else {
			// fs upload f1 f2 /tmp/ 或 fs upload f1 f2 /tmp/f.txt
			srcPaths = args[:len(args)-1]
			remotePath = args[len(args)-1]
		}

		// 3. 建立 SCP 客户端（dry-run 不需要连接）
		var client scp.Client
		if !dryRun {
			client, err := utils.EstablishScpClient(cfg)
			if err != nil {
				log.Errorf("建立 SCP 连接失败: %v", err)
				return
			}
			defer client.Close()
		}

		// 4. 解析远程路径
		remoteDir, remoteFilename := utils.ParseRemotePath(remotePath)
		if remoteDir == "" {
			remoteDir = cfg.DefaultUploadDir
			if remoteDir == "" {
				remoteDir = "~"
			}
		}
		remoteDir = utils.ExpandHome(remoteDir)

		// 5. 批量处理每个文件
		successCount := 0
		for _, srcPath := range srcPaths {
			// 5.1 解析本地路径
			srcFile, exists := utils.ResolveLocalFile(srcPath)
			if !exists {
				log.Errorf("跳过：本地文件不存在: %s", srcPath)
				continue
			}

			// 5.2 获取文件信息
			fileInfo, err := os.Stat(srcFile)
			if err != nil {
				log.Errorf("无法读取文件信息: %s: %v", srcFile, err)
				continue
			}
			if fileInfo.IsDir() {
				log.Warnf("跳过目录: %s", srcFile)
				continue
			}

			// 5.3 确定远程目标文件名
			var finalRemoteFile string
			if len(srcPaths) == 1 && remoteFilename != "" {
				// 单文件上传且指定了完整路径
				finalRemoteFile = remoteFilename
			} else {
				// 多文件或未指定文件名，使用原文件名
				finalRemoteFile = fileInfo.Name()
			}
			fullRemotePath := path.Join(remoteDir, finalRemoteFile)

			// 5.4 Dry-run 模式：只打印
			if dryRun {
				fmt.Printf("🟡 [dry-run] 将上传: %s → %s@%s:%s\n", srcFile, cfg.User, cfg.Host, fullRemotePath)
				successCount++
				continue
			}

			// 5.5 实际上传
			f, err := os.Open(srcFile)
			if err != nil {
				log.Errorf("打开文件失败: %s: %v", srcFile, err)
				continue
			}

			// 进度条处理
			err = client.CopyFromFilePassThru(context.Background(), *f, fullRemotePath, chmod, passThru)
			f.Close() // 立即关闭文件

			if err != nil {
				log.Errorf("上传失败: %s → %s: %v", srcFile, fullRemotePath, err)
				if bar != nil {
					bar.Finish()
				}
				continue
			}

			if bar != nil {
				bar.Finish()
			}

			log.Infof("✅ 上传成功: %s → %s@%s:%s", srcFile, cfg.User, cfg.Host, fullRemotePath)
			successCount++
		}

		// 6. 总结
		if dryRun {
			fmt.Printf("🟢 [dry-run] 预计上传 %d 个文件\n", successCount)
		} else {
			log.Infof("批量上传完成，成功 %d/%d 个文件", successCount, len(srcPaths))
		}

		// 7. 记录命令历史（仅非 dry-run）
		if !dryRun {
			record(append([]string{"fs", "upload"}, args...))
		}
	},
}

func init() {
	// 添加标志
	uploadCmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "预演操作，不实际上传")
	uploadCmd.Flags().StringVar(&chmod, "chmod", "", "设置远程文件权限（如 0644）")

	rootCmd.AddCommand(uploadCmd)
}
