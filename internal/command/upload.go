/*
Copyright © 2025 SHEN YI <1245332635@qq.com>
*/

package command

import (
	"context"
	"fmt"
	"fs/internal/types"
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
	names  string // 服务器配置名，如: web1,db-prod
	tagStr string // 标签，如: web,prod
)

var uploadCmd = &cobra.Command{
	Use:   "upload [local_file...] [remote_dir|remote_file]",
	Short: "上传文件到远程服务器",
	Long: `通过 SCP 上传本地文件到远程服务器。

支持：
  - 批量上传多个文件
  - 密钥认证
  - 进度条显示
  - --dry-run 模式预演
  - 指定多个目标服务器（通过 --names 或 --tags）

示例：
  fs upload myfile.txt
  fs upload f1.txt f2.txt /tmp/
  fs upload *.log backup/
  fs upload -n myfile.txt /tmp/
  fs upload f.txt--names server1,server2
  fs upload f.txt --tags web,prod`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.PrintErrln("❌ 未指定要上传的文件")
			cmd.Usage()
			os.Exit(1)
		}

		if chmod == "" {
			chmod = "0644"
		}

		var srcPaths []string
		var remotePath string
		if len(args) == 1 {
			srcPaths = args[:1]
			remotePath = ""
		} else {
			srcPaths = args[:len(args)-1]
			remotePath = args[len(args)-1]
		}

		// 获取目标服务器列表（names / tags / 当前）
		configs, err := getTargetConfigs(cmd)
		if err != nil {
			log.Fatalf("❌ 获取目标服务器失败: %v", err)
		}

		if len(configs) == 0 {
			log.Fatal("❌ 没有可用的目标服务器")
		}

		for _, targetCfg := range configs {
			log.Infof("📤 开始上传到服务器: %s (%s@%s:%d)", targetCfg.Name, targetCfg.User, targetCfg.Host, targetCfg.Port)

			var client scp.Client
			if !dryRun {
				var err error
				client, err = utils.EstablishScpClient(targetCfg)
				if err != nil {
					log.Errorf("❌ 建立 SCP 连接失败 (%s): %v", targetCfg.Name, err)
					continue
				}
				defer client.Close()
			}

			remoteDir, remoteFilename := utils.ParseRemotePath(remotePath)
			if remoteDir == "" {
				remoteDir = targetCfg.DefaultUploadDir
				if remoteDir == "" {
					remoteDir = "~"
				}
			}
			remoteDir = utils.ExpandHome(remoteDir)

			successCount := 0
			for _, srcPath := range srcPaths {
				srcFile, exists := utils.ResolveLocalFile(srcPath)
				if !exists {
					log.Errorf("🟡 跳过：本地文件不存在: %s", srcPath)
					continue
				}

				fileInfo, err := os.Stat(srcFile)
				if err != nil {
					log.Errorf("🟡 无法读取文件信息: %s: %v", srcFile, err)
					continue
				}
				if fileInfo.IsDir() {
					log.Warnf("🟡 跳过目录: %s", srcFile)
					continue
				}

				var finalRemoteFile string
				if len(srcPaths) == 1 && remoteFilename != "" {
					finalRemoteFile = remoteFilename
				} else {
					finalRemoteFile = fileInfo.Name()
				}
				fullRemotePath := path.Join(remoteDir, finalRemoteFile)

				if dryRun {
					fmt.Printf("🟡 [dry-run] 将上传: %s → %s@%s:%s\n", srcFile, targetCfg.User, targetCfg.Host, fullRemotePath)
					successCount++
					continue
				}

				f, err := os.Open(srcFile)
				if err != nil {
					log.Errorf("🟡 打开文件失败: %s: %v", srcFile, err)
					continue
				}

				err = client.CopyFromFilePassThru(context.Background(), *f, fullRemotePath, chmod, passThru)
				f.Close()

				if err != nil {
					log.Errorf("❌ 上传失败 (%s): %s → %s: %v", targetCfg.Name, srcFile, fullRemotePath, err)
					if bar != nil {
						bar.Finish()
					}
					continue
				}

				if bar != nil {
					bar.Finish()
				}

				log.Infof("✅ 上传成功 (%s): %s → %s@%s:%s", targetCfg.Name, srcFile, targetCfg.User, targetCfg.Host, fullRemotePath)
				successCount++
			}

			if dryRun {
				fmt.Printf("🟢 [dry-run] 预计上传 %d 个文件到 %s\n", successCount, targetCfg.Name)
			} else {
				log.Infof("✅ 上传完成 (%s): 成功 %d/%d 个文件", targetCfg.Name, successCount, len(srcPaths))
			}
		}

		// 仅当使用当前服务器时记录历史
		if !dryRun && !cmd.Flags().Changed("names") && !cmd.Flags().Changed("tags") {
			record(append([]string{"fs", "upload"}, args...))
		}
	},
}

// getTargetConfigs 根据 names/tags 获取服务器列表（互斥）
func getTargetConfigs(cmd *cobra.Command) ([]*types.SshConfig, error) {
	usingNames := cmd.Flags().Changed("names")
	usingTags := cmd.Flags().Changed("tags")

	switch {
	case usingNames && usingTags:
		return nil, fmt.Errorf("`--names` 和 `--tags` 不能同时使用，请选择其中一个")

	case usingNames:
		nameList := utils.NonEmptyTrimmedSplit(names, ",")
		var configs []*types.SshConfig
		for _, n := range nameList {
			cfg, err := utils.LoadConfigByName(fsHostsConfigDir, n)
			if err != nil {
				return nil, fmt.Errorf("加载配置失败 (%s): %v", n, err)
			}
			configs = append(configs, cfg)
		}
		return configs, nil

	case usingTags:
		tagList := utils.NonEmptyTrimmedSplit(tagStr, ",")
		allConfigs, err := utils.LoadAllConfigs(fsHostsConfigDir)
		if err != nil {
			return nil, fmt.Errorf("加载所有配置失败: %v", err)
		}

		var matched []*types.SshConfig
		for _, cfg := range allConfigs {
			if utils.HasAnyTag(cfg.Tags, tagList) {
				matched = append(matched, cfg)
			}
		}
		if len(matched) == 0 {
			return nil, fmt.Errorf("没有找到包含标签 %v 的服务器", tagList)
		}
		return matched, nil

	default:
		if cfg == nil {
			return nil, fmt.Errorf("当前无选中服务器，请使用 'fs use <name>' 切换，或使用 --names/--tags 指定目标")
		}
		return []*types.SshConfig{cfg}, nil
	}
}

func init() {
	uploadCmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "预演操作，不实际上传")
	uploadCmd.Flags().StringVar(&chmod, "chmod", "", "设置远程文件权限（如 0644）")
	uploadCmd.Flags().StringVar(&names, "names", "", "指定多个目标服务器（配置名），用逗号分隔")
	uploadCmd.Flags().StringVar(&tagStr, "tags", "", "上传到匹配这些标签的服务器，用逗号分隔")

	rootCmd.AddCommand(uploadCmd)
}
