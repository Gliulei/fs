/*
Copyright © 2025 SHEN YI <1245332635@qq.com>
*/

package command

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ngaut/log"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	forceImport   bool
	skipExisting  bool
	previewImport bool
)

var importCmd = &cobra.Command{
	Use:   "import [archive_file]",
	Short: "从打包文件导入主机配置",
	Long: `从 tar.gz 打包文件导入主机配置。

支持：
  - 配置冲突处理（覆盖、跳过、重命名）
  - 导入预览功能
  - 批量导入操作
  - 配置文件验证

示例：
  fs import config-archive.tar.gz
  fs import config-archive.tar.gz --force
  fs import config-archive.tar.gz --skip-existing
  fs import config-archive.tar.gz --preview`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := runImport(args[0]); err != nil {
			log.Errorf("❌ 导入失败: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	importCmd.Flags().BoolVarP(&forceImport, "force", "f", false, "强制覆盖已存在的配置")
	importCmd.Flags().BoolVarP(&skipExisting, "skip-existing", "s", false, "跳过已存在的配置")
	importCmd.Flags().BoolVarP(&previewImport, "preview", "p", false, "预览将要导入的配置")
}

func runImport(archivePath string) error {
	// 验证打包文件
	if err := validateArchive(archivePath); err != nil {
		return fmt.Errorf("打包文件验证失败: %v", err)
	}

	// 解析打包文件
	archiveInfo, err := parseArchive(archivePath)
	if err != nil {
		return fmt.Errorf("解析打包文件失败: %v", err)
	}

	// 预览模式
	if previewImport {
		return showImportPreview(archiveInfo)
	}

	// 检查冲突
	conflicts := detectConflicts(archiveInfo.ConfigFiles)
	if len(conflicts) > 0 && !forceImport && !skipExisting {
		if err := handleConflicts(conflicts); err != nil {
			return err
		}
	}

	// 执行导入
	result, err := executeImport(archivePath, conflicts)
	if err != nil {
		return fmt.Errorf("导入执行失败: %v", err)
	}

	// 显示导入结果
	showImportResult(result)
	return nil
}

type ArchiveInfo struct {
	Version     string
	Metadata    map[string]interface{}
	ConfigFiles []ConfigFile
}

type ConfigFile struct {
	Name    string
	Content []byte
}

type ImportResult struct {
	SuccessCount int
	SkipCount    int
	ErrorCount   int
	Errors       []string
	Imported     []string
	Skipped      []string
}

func validateArchive(archivePath string) error {
	// 检查文件是否存在
	if _, err := os.Stat(archivePath); os.IsNotExist(err) {
		return fmt.Errorf("打包文件不存在: %s", archivePath)
	}

	// 检查文件扩展名
	if !strings.HasSuffix(archivePath, ".tar.gz") {
		return fmt.Errorf("不支持的文件格式，仅支持 .tar.gz 文件")
	}

	return nil
}

func parseArchive(archivePath string) (*ArchiveInfo, error) {
	file, err := os.Open(archivePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 创建 gzip reader
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return nil, fmt.Errorf("创建 gzip reader 失败: %v", err)
	}
	defer gzipReader.Close()

	// 创建 tar reader
	tarReader := tar.NewReader(gzipReader)

	info := &ArchiveInfo{
		ConfigFiles: make([]ConfigFile, 0),
	}

	// 读取文件
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("读取打包文件失败: %v", err)
		}

		switch header.Name {
		case "VERSION":
			versionData := make([]byte, header.Size)
			if _, err := io.ReadFull(tarReader, versionData); err != nil {
				return nil, fmt.Errorf("读取版本信息失败: %v", err)
			}
			info.Version = string(versionData)
		case "metadata.json":
			metadataData := make([]byte, header.Size)
			if _, err := io.ReadFull(tarReader, metadataData); err != nil {
				return nil, fmt.Errorf("读取元数据失败: %v", err)
			}
			if err := json.Unmarshal(metadataData, &info.Metadata); err != nil {
				return nil, fmt.Errorf("解析元数据失败: %v", err)
			}
		default:
			// 处理配置文件
			if strings.HasPrefix(header.Name, "config/") && strings.HasSuffix(header.Name, ".yaml") {
				content := make([]byte, header.Size)
				if _, err := io.ReadFull(tarReader, content); err != nil {
					return nil, fmt.Errorf("读取配置文件失败 %s: %v", header.Name, err)
				}

				// 验证 YAML 格式
				var config map[string]interface{}
				if err := yaml.Unmarshal(content, &config); err != nil {
					return nil, fmt.Errorf("配置文件格式无效 %s: %v", header.Name, err)
				}

				filename := filepath.Base(header.Name)
				configName := strings.TrimSuffix(filename, ".yaml")

				info.ConfigFiles = append(info.ConfigFiles, ConfigFile{
					Name:    configName,
					Content: content,
				})
			}
		}
	}

	return info, nil
}

func showImportPreview(info *ArchiveInfo) error {
	fmt.Println("🔍 导入预览:")
	fmt.Printf("📦 打包文件版本: %s\n", info.Version)
	fmt.Printf("📋 配置文件数量: %d\n", len(info.ConfigFiles))
	fmt.Println("\n📁 将要导入的配置:")

	for _, configFile := range info.ConfigFiles {
		fmt.Printf("  - %s.yaml\n", configFile.Name)
	}

	if info.Metadata != nil {
		if createdAt, ok := info.Metadata["created_at"].(string); ok {
			fmt.Printf("\n📅 打包时间: %s\n", createdAt)
		}
	}

	return nil
}

func detectConflicts(configFiles []ConfigFile) []string {
	var conflicts []string

	for _, configFile := range configFiles {
		configPath := filepath.Join(fsHostsConfigDir, configFile.Name+".yaml")
		if _, err := os.Stat(configPath); err == nil {
			conflicts = append(conflicts, configFile.Name)
		}
	}

	return conflicts
}

func handleConflicts(conflicts []string) error {
	if len(conflicts) == 0 {
		return nil
	}

	fmt.Printf("⚠️  检测到 %d 个配置冲突:\n", len(conflicts))
	for _, conflict := range conflicts {
		fmt.Printf("  - %s\n", conflict)
	}

	fmt.Println("\n请选择处理方式:")
	fmt.Println("1. 覆盖所有冲突配置")
	fmt.Println("2. 跳过所有冲突配置")
	fmt.Println("3. 重命名导入的配置")
	fmt.Println("4. 取消导入")

	var choice int
	fmt.Print("请输入选择 (1-4): ")
	if _, err := fmt.Scanf("%d", &choice); err != nil {
		return fmt.Errorf("输入无效")
	}

	switch choice {
	case 1:
		forceImport = true
	case 2:
		skipExisting = true
	case 3:
		// TODO: 实现重命名逻辑
		return fmt.Errorf("重命名功能暂未实现")
	case 4:
		return fmt.Errorf("用户取消导入")
	default:
		return fmt.Errorf("无效选择")
	}

	return nil
}

func executeImport(archivePath string, conflicts []string) (*ImportResult, error) {
	result := &ImportResult{
		Errors: make([]string, 0),
	}

	info, err := parseArchive(archivePath)
	if err != nil {
		return nil, err
	}

	// 确保 hosts 目录存在
	if err := os.MkdirAll(fsHostsConfigDir, 0755); err != nil {
		return nil, fmt.Errorf("创建 hosts 目录失败: %v", err)
	}

	for _, configFile := range info.ConfigFiles {
		configPath := filepath.Join(fsHostsConfigDir, configFile.Name+".yaml")

		// 检查是否已存在
		if _, err := os.Stat(configPath); err == nil {
			if skipExisting {
				result.SkipCount++
				result.Skipped = append(result.Skipped, configFile.Name)
				continue
			}
			if !forceImport {
				// 检查是否在冲突列表中
				isConflict := false
				for _, conflict := range conflicts {
					if conflict == configFile.Name {
						isConflict = true
						break
					}
				}
				if isConflict {
					result.SkipCount++
					result.Skipped = append(result.Skipped, configFile.Name)
					continue
				}
			}
		}

		// 写入配置文件
		if err := os.WriteFile(configPath, configFile.Content, 0644); err != nil {
			result.ErrorCount++
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", configFile.Name, err))
			continue
		}

		result.SuccessCount++
		result.Imported = append(result.Imported, configFile.Name)
	}

	return result, nil
}

func showImportResult(result *ImportResult) {
	fmt.Println("\n✅ 导入完成:")
	fmt.Printf("  成功导入: %d 个\n", result.SuccessCount)
	fmt.Printf("  跳过: %d 个\n", result.SkipCount)
	fmt.Printf("  错误: %d 个\n", result.ErrorCount)

	if len(result.Imported) > 0 {
		fmt.Println("\n📥 成功导入的配置:")
		for _, name := range result.Imported {
			fmt.Printf("  - %s\n", name)
		}
	}

	if len(result.Skipped) > 0 {
		fmt.Println("\n⏭️  跳过的配置:")
		for _, name := range result.Skipped {
			fmt.Printf("  - %s\n", name)
		}
	}

	if len(result.Errors) > 0 {
		fmt.Println("\n❌ 导入错误:")
		for _, err := range result.Errors {
			fmt.Printf("  - %s\n", err)
		}
	}
}

func init() {
	rootCmd.AddCommand(importCmd)
}
