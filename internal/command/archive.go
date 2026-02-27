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
	"time"

	"github.com/ngaut/log"
	"github.com/spf13/cobra"
)

var (
	outputFile       string
	compressionLevel int
)

var archiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "打包所有主机配置到压缩文件",
	Long: `将所有主机配置文件打包成一个 tar.gz 文件，便于备份和分发。

支持：
  - 自定义输出文件路径
  - 可调节压缩级别
  - 包含元数据和版本信息

示例：
  fs archive
  fs archive --output my-configs.tar.gz
  fs archive --compression-level 9`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runArchive(); err != nil {
			log.Errorf("❌ 打包失败: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	archiveCmd.Flags().StringVarP(&outputFile, "output", "o", "", "输出文件路径 (默认: fs-config-archive-时间戳.tar.gz)")
	archiveCmd.Flags().IntVarP(&compressionLevel, "compression-level", "l", 6, "压缩级别 (1-9, 默认: 6)")
}

func runArchive() error {
	// 确定输出文件名
	if outputFile == "" {
		timestamp := time.Now().Format("20060102-150405")
		outputFile = fmt.Sprintf("fs-config-archive-%s.tar.gz", timestamp)
	}

	// 检查输出文件是否已存在
	if _, err := os.Stat(outputFile); err == nil {
		return fmt.Errorf("输出文件已存在: %s", outputFile)
	}

	// 收集所有配置文件
	configFiles, err := collectConfigFiles()
	if err != nil {
		return fmt.Errorf("收集配置文件失败: %v", err)
	}

	if len(configFiles) == 0 {
		return fmt.Errorf("未找到任何配置文件")
	}

	log.Infof("✅ 找到 %d 个配置文件", len(configFiles))

	// 创建打包文件
	if err := createArchive(outputFile, configFiles); err != nil {
		return fmt.Errorf("创建打包文件失败: %v", err)
	}

	log.Infof("✅ 配置打包完成: %s", outputFile)
	return nil
}

func collectConfigFiles() ([]string, error) {
	var files []string

	// 检查 hosts 目录是否存在
	if _, err := os.Stat(fsHostsConfigDir); os.IsNotExist(err) {
		return files, nil
	}

	// 遍历 hosts 目录
	err := filepath.Walk(fsHostsConfigDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 只处理 .yaml 文件
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".yaml") {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

func createArchive(archivePath string, configFiles []string) error {
	// 创建输出文件
	file, err := os.Create(archivePath)
	if err != nil {
		return fmt.Errorf("创建输出文件失败: %v", err)
	}
	defer file.Close()

	// 创建 gzip writer
	gzipWriter, err := gzip.NewWriterLevel(file, compressionLevel)
	if err != nil {
		return fmt.Errorf("创建 gzip writer 失败: %v", err)
	}
	defer gzipWriter.Close()

	// 创建 tar writer
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	// 添加配置文件
	for _, configFile := range configFiles {
		if err := addFileToArchive(tarWriter, configFile, "config/"); err != nil {
			return fmt.Errorf("添加文件到打包失败 %s: %v", configFile, err)
		}
	}

	// 创建并添加 metadata.json
	if err := addMetadata(tarWriter, configFiles); err != nil {
		return fmt.Errorf("添加元数据失败: %v", err)
	}

	// 创建并添加 VERSION 文件
	if err := addVersionFile(tarWriter); err != nil {
		return fmt.Errorf("添加版本文件失败: %v", err)
	}

	return nil
}

func addFileToArchive(tarWriter *tar.Writer, filePath, prefix string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	// 获取相对路径作为文件名
	relPath, err := filepath.Rel(fsHostsConfigDir, filePath)
	if err != nil {
		return err
	}

	// 标准化路径分隔符为正斜杠（tar 标准）
	standardPath := filepath.ToSlash(filepath.Join(prefix, relPath))

	header := &tar.Header{
		Name:    standardPath,
		Mode:    0644,
		Size:    info.Size(),
		ModTime: info.ModTime(),
	}

	if err := tarWriter.WriteHeader(header); err != nil {
		return err
	}

	_, err = io.Copy(tarWriter, file)
	return err
}

func addMetadata(tarWriter *tar.Writer, configFiles []string) error {
	// 准备元数据
	metadata := map[string]interface{}{
		"created_at":        time.Now().Format(time.RFC3339),
		"fs_version":        "unknown", // TODO: 从 version 包获取
		"config_count":      len(configFiles),
		"config_names":      getConfigNames(configFiles),
		"compression_level": compressionLevel,
	}

	// 序列化为 JSON
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return err
	}

	// 创建 tar header
	header := &tar.Header{
		Name: "metadata.json",
		Mode: 0644,
		Size: int64(len(data)),
	}

	if err := tarWriter.WriteHeader(header); err != nil {
		return err
	}

	_, err = tarWriter.Write(data)
	return err
}

func addVersionFile(tarWriter *tar.Writer) error {
	// TODO: 从 version 包获取实际版本
	versionData := []byte("0.1.0") // 临时版本号

	header := &tar.Header{
		Name: "VERSION",
		Mode: 0644,
		Size: int64(len(versionData)),
	}

	if err := tarWriter.WriteHeader(header); err != nil {
		return err
	}

	_, err := tarWriter.Write(versionData)
	return err
}

func getConfigNames(configFiles []string) []string {
	var names []string
	for _, file := range configFiles {
		name := strings.TrimSuffix(filepath.Base(file), ".yaml")
		names = append(names, name)
	}
	return names
}

func init() {
	rootCmd.AddCommand(archiveCmd)
}
