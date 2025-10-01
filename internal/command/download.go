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
	"path/filepath"

	"github.com/bramvdbogaerde/go-scp"
	"github.com/ngaut/log"
	"github.com/spf13/cobra"
)

var (
	dryRunDownload bool
	overwrite      bool
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download [remote_file...] [local_dir]",
	Short: "从远程服务器下载文件",
	Long: `通过 SCP 从远程服务器下载文件。

支持：
  - 批量下载多个文件
  - 下载到指定本地目录
  - --dry-run 模式预演
  - 自动展开 ~ 路径

示例：
  fs download myfile.txt                    # 下载到默认下载目录
  fs download f1.txt f2.txt ./local/        # 批量下载到本地目录
  fs download /tmp/*.log                    # 通配符（由 shell 展开）
  fs download -n myfile.txt                 # 预演操作，不实际下载
  fs download file.txt ~/Downloads/         # 指定本地目录`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.PrintErrln("❌ 未指定要下载的远程文件")
			cmd.Usage()
			os.Exit(1)
		}

		// 1. 分离远程文件和本地目录
		var remotePaths []string
		var localDir string

		if len(args) == 1 {
			// fs download file.txt
			remotePaths = args[:1]
			localDir = cfg.DefaultDownloadDir
			if localDir == "" {
				localDir = "." // 当前目录
			}
		} else {
			// fs download f1 f2 ./local/
			remotePaths = args[:len(args)-1]
			localDir = args[len(args)-1]
		}

		// 2. 解析本地目录（支持 ~）
		localDir = utils.ExpandHome(localDir)
		if !filepath.IsAbs(localDir) {
			// 转为绝对路径
			absDir, err := filepath.Abs(localDir)
			if err != nil {
				log.Errorf("无法解析本地目录: %s: %v", localDir, err)
				return
			}
			localDir = absDir
		}

		// 3. 确保本地目录存在
		if !dryRunDownload {
			if err := utils.EnsureDir(localDir); err != nil {
				log.Errorf("无法创建本地目录: %s: %v", localDir, err)
				return
			}
		}

		// 4. 建立 SCP 客户端（dry-run 不需要连接）
		var client scp.Client
		if !dryRunDownload {
			var err error
			client, err = utils.EstablishScpClient(cfg)
			if err != nil {
				log.Errorf("建立 SCP 连接失败: %v", err)
				return
			}
			defer client.Close()
		}

		// 5. 批量下载每个远程文件
		successCount := 0
		for _, remotePath := range remotePaths {
			// 5.1 解析远程路径
			if !path.IsAbs(remotePath) {
				// 非绝对路径，使用默认上传目录
				if cfg.DefaultUploadDir != "" {
					remotePath = path.Join(cfg.DefaultUploadDir, remotePath)
				} else {
					remotePath = path.Join("~", remotePath)
				}
			}
			remotePath = utils.ExpandHome(remotePath) // 支持 ~/file.txt

			filename := path.Base(remotePath)
			localFilePath := filepath.Join(localDir, filename)

			// 5.2 Dry-run 模式：只打印
			if dryRunDownload {
				fmt.Printf("🟡 [dry-run] 将下载: %s@%s:%s → %s\n", cfg.User, cfg.Host, remotePath, localFilePath)
				successCount++
				continue
			}

			// 5.3 检查本地文件是否已存在
			if utils.FileExists(localFilePath) && !overwrite {
				log.Warnf("跳过：文件已存在（使用 --overwrite 覆盖）: %s", localFilePath)
				continue
			}

			// 5.4 创建本地文件
			f, err := os.OpenFile(localFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				log.Errorf("无法创建本地文件: %s: %v", localFilePath, err)
				continue
			}
			// 关闭文件
			defer f.Close()

			// 5.5 下载文件
			err = client.CopyFromRemotePassThru(context.Background(), f, remotePath, passThru)

			if err != nil {
				log.Errorf("下载失败: %s → %s: %v", remotePath, localFilePath, err)
				if bar != nil {
					bar.Finish()
				}
				continue
			}

			if bar != nil {
				bar.Finish()
			}

			log.Infof("✅ 下载成功: %s@%s:%s → %s", cfg.User, cfg.Host, remotePath, localFilePath)
			successCount++
		}

		// 6. 总结
		if dryRunDownload {
			fmt.Printf("🟢 [dry-run] 预计下载 %d 个文件到: %s\n", successCount, localDir)
		} else {
			log.Infof("批量下载完成，成功 %d/%d 个文件", successCount, len(remotePaths))
		}

		// 7. 记录命令历史（仅非 dry-run）
		if !dryRunDownload {
			cmdLog := append([]string{"fs", "download"}, args...)
			record(cmdLog)
		}
	},
}

func init() {
	// 添加标志
	downloadCmd.Flags().BoolVarP(&dryRunDownload, "dry-run", "n", false, "预演操作，不实际下载")
	downloadCmd.Flags().BoolVar(&overwrite, "overwrite", false, "覆盖已存在的本地文件")

	rootCmd.AddCommand(downloadCmd)
}
