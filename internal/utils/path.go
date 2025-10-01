/*
Copyright © 2025 SHEN YI <1245332635@qq.com>
*/

// pkg/utils/path.go
package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// ExpandHome 将路径中的 ~ 替换为用户主目录
func ExpandHome(path string) string {
	if !strings.HasPrefix(path, "~") {
		return path
	}

	var home string
	var err error

	// 处理 ~user 形式（可选扩展）
	if len(path) > 1 && path[1] != '/' && path[1] != '\\' {
		// 这里暂不支持 ~alice，直接返回
		return path
	}

	home, err = os.UserHomeDir()
	if err != nil {
		// 如果无法获取 home，保留原路径
		return path
	}

	if path == "~" {
		return home
	}

	return filepath.Join(home, path[2:])
}

// ResolveLocalFile 解析本地文件路径，支持：
// - 绝对路径
// - 相对路径
// - ~ 展开
// 返回解析后的绝对路径（如果存在），否则返回原路径 + false
func ResolveLocalFile(src string) (string, bool) {
	// 1. 展开 ~
	src = ExpandHome(src)

	// 2. 如果已经是绝对路径，直接检查是否存在
	if filepath.IsAbs(src) {
		if FileExists(src) {
			return src, true
		}
		return src, false
	}

	// 3. 尝试在当前目录下查找
	currentDir := getCurrentDir()
	if currentDir != "" {
		abs := filepath.Join(currentDir, src)
		if FileExists(abs) {
			return abs, true
		}
	}

	// 5. 都找不到，返回原始路径（可能用于报错）
	return src, false
}

// ParseRemotePath 解析远程目标路径，返回 (dir, filename)
// 支持: /tmp/file.txt → (/tmp, file.txt)
// 支持: /tmp/        → (/tmp, "")
// 支持: ""           → ("", "")
func ParseRemotePath(remotePath string) (dir, filename string) {
	if remotePath == "" {
		return "", ""
	}

	// 展开 ~
	remotePath = ExpandHome(remotePath)

	dir, filename = filepath.ToSlash(remotePath), ""
	if dir[len(dir)-1] == '/' {
		// 路径以 / 结尾，说明是目录
		dir = dir[:len(dir)-1] // 去掉末尾 /
		filename = filepath.Base(dir)
		dir = filepath.Dir(dir)
	} else {
		dir, filename = filepath.Split(dir)
	}

	return dir, filename
}

// fileExists 检查文件是否存在
func FileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// IsDir 检查路径是否为目录
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// EnsureDir 确保目录存在，不存在则创建
func EnsureDir(dir string) error {
	if dir == "" {
		return nil
	}
	return os.MkdirAll(dir, 0755)
}

func getCurrentDir() string {
	wd, _ := os.Getwd()
	return wd
}
