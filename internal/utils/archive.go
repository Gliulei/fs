/*
Copyright © 2025 SHEN YI <1245332635@qq.com>
*/

package utils

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ArchiveOptions 打包选项
type ArchiveOptions struct {
	CompressionLevel int
	OutputPath       string
	IncludeMetadata  bool
}

// ArchiveFileInfo 归档文件信息
type ArchiveFileInfo struct {
	Name    string
	Path    string
	Size    int64
	ModTime time.Time
}

// CreateArchive 创建配置文件打包
func CreateArchive(configFiles []string, outputPath string, options ArchiveOptions) error {
	// 创建输出文件
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("创建输出文件失败: %v", err)
	}
	defer file.Close()

	// 创建 gzip writer
	gzipWriter, err := gzip.NewWriterLevel(file, options.CompressionLevel)
	if err != nil {
		return fmt.Errorf("创建 gzip writer 失败: %v", err)
	}
	defer gzipWriter.Close()

	// 创建 tar writer
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	// 添加配置文件
	for _, configFile := range configFiles {
		if err := addFileToTar(tarWriter, configFile, "config/"); err != nil {
			return fmt.Errorf("添加文件到打包失败 %s: %v", configFile, err)
		}
	}

	// 如果需要包含元数据
	if options.IncludeMetadata {
		if err := addMetadataToTar(tarWriter, configFiles); err != nil {
			return fmt.Errorf("添加元数据失败: %v", err)
		}
	}

	return nil
}

// ExtractArchive 解压配置文件打包
func ExtractArchive(archivePath, extractDir string) ([]ArchiveFileInfo, error) {
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

	var files []ArchiveFileInfo

	// 读取文件
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("读取打包文件失败: %v", err)
		}

		// 跳过目录和特殊文件
		if header.Typeflag != tar.TypeReg {
			continue
		}

		// 构造输出路径
		outputPath := filepath.Join(extractDir, header.Name)

		// 确保目录存在
		if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
			return nil, fmt.Errorf("创建目录失败 %s: %v", filepath.Dir(outputPath), err)
		}

		// 创建输出文件
		outputFile, err := os.Create(outputPath)
		if err != nil {
			return nil, fmt.Errorf("创建文件失败 %s: %v", outputPath, err)
		}

		// 复制文件内容
		if _, err := io.Copy(outputFile, tarReader); err != nil {
			outputFile.Close()
			return nil, fmt.Errorf("写入文件失败 %s: %v", outputPath, err)
		}

		outputFile.Close()

		// 记录文件信息
		files = append(files, ArchiveFileInfo{
			Name:    filepath.Base(header.Name),
			Path:    outputPath,
			Size:    header.Size,
			ModTime: header.ModTime,
		})
	}

	return files, nil
}

// ValidateArchive 验证打包文件的完整性
func ValidateArchive(archivePath string) error {
	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 创建 gzip reader
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("无效的 gzip 文件: %v", err)
	}
	defer gzipReader.Close()

	// 创建 tar reader
	tarReader := tar.NewReader(gzipReader)

	// 检查基本结构
	requiredFiles := map[string]bool{
		"VERSION":       false,
		"metadata.json": false,
	}

	configFileCount := 0

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("读取打包文件失败: %v", err)
		}

		// 检查必需文件
		if _, exists := requiredFiles[header.Name]; exists {
			requiredFiles[header.Name] = true
		}

		// 统计配置文件
		if strings.HasPrefix(header.Name, "config/") && strings.HasSuffix(header.Name, ".yaml") {
			configFileCount++
		}
	}

	// 验证必需文件是否存在
	for filename, exists := range requiredFiles {
		if !exists {
			return fmt.Errorf("缺少必需文件: %s", filename)
		}
	}

	// 验证至少有一个配置文件
	if configFileCount == 0 {
		return fmt.Errorf("打包文件中没有配置文件")
	}

	return nil
}

// ListArchiveContents 列出打包文件内容
func ListArchiveContents(archivePath string) ([]ArchiveFileInfo, error) {
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

	var files []ArchiveFileInfo

	// 读取文件列表
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("读取打包文件失败: %v", err)
		}

		// 只处理普通文件
		if header.Typeflag == tar.TypeReg {
			files = append(files, ArchiveFileInfo{
				Name:    header.Name,
				Path:    header.Name,
				Size:    header.Size,
				ModTime: header.ModTime,
			})
		}
	}

	return files, nil
}

// Helper functions

func addFileToTar(tarWriter *tar.Writer, filePath, prefix string) error {
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
	relPath := filepath.Base(filePath)

	header := &tar.Header{
		Name:    filepath.Join(prefix, relPath),
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

func addMetadataToTar(tarWriter *tar.Writer, configFiles []string) error {
	// 这里可以添加元数据生成逻辑
	// 暂时留空，后续可以实现
	return nil
}
