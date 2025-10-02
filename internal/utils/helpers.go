/*
Copyright © 2025 SHEN YI <1245332635@qq.com>
*/

package utils

import (
	"fmt"
	"fs/internal/types"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ngaut/log"
	"gopkg.in/yaml.v3"
)

/**
 * @Description: 去除空格和空行
 * @param {string} s
 * @param {string} sep
 * @return {[]string}
 */
func NonEmptyTrimmedSplit(s, sep string) []string {
	parts := strings.Split(s, sep)
	var result []string
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

/**
 * HasAnyTag 检查目标标签列表中是否存在任意一个配置标签.
 *
 * @param configTags 需要匹配的标签列表
 * @param targetTags 目标标签列表.
 * @return 果targetTags中包含任意一个configTags中的标签则返回true，否则返回false.
 */
func HasAnyTag(configTags, targetTags []string) bool {
	tagSet := make(map[string]bool)
	for _, t := range configTags {
		tagSet[t] = true
	}
	for _, t := range targetTags {
		if tagSet[t] {
			return true
		}
	}
	return false
}

/*
*
LoadAllConfigs 从配置目录加载所有 .yaml 文件，返回成功解析的 SSH 配置列表
注意：会跳过无法读取或解析失败的文件，仅返回成功加载的配置
*/
func LoadAllConfigs(fsHostsConfigDir string) ([]*types.SshConfig, error) {
	pattern := filepath.Join(fsHostsConfigDir, "*.yaml")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("匹配配置文件失败 %s: %v", pattern, err)
	}

	if len(files) == 0 {
		log.Infof("📭 配置目录中未找到任何 .yaml 文件: %s", fsHostsConfigDir)
		return []*types.SshConfig{}, nil
	}

	// 排序：保证加载顺序一致（按文件名）
	sort.Strings(files)

	var configs []*types.SshConfig
	var loadErrors []string

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			errMsg := fmt.Sprintf("⚠️ 读取配置文件失败 %s: %v", file, err)
			log.Error(errMsg)
			loadErrors = append(loadErrors, errMsg)
			continue
		}

		var config types.SshConfig // 注意：直接定义结构体，避免指针初始化问题
		if err := yaml.Unmarshal(data, &config); err != nil {
			errMsg := fmt.Sprintf("⚠️ YAML 解析失败 %s: %v", file, err)
			log.Error(errMsg)
			loadErrors = append(loadErrors, errMsg)
			continue
		}

		configs = append(configs, &config)
	}

	// 可选：按 Name 排序
	sort.Slice(configs, func(i, j int) bool {
		return configs[i].Name < configs[j].Name
	})

	// 如果有错误，返回警告信息但不中断（除非完全没加载成功）
	if len(loadErrors) > 0 {
		if len(configs) == 0 {
			return nil, fmt.Errorf("❌ 未成功加载任何配置文件，共失败 %d 个:\n  %s", len(loadErrors), strings.Join(loadErrors, "\n  "))
		} else {
			log.Warnf("部分配置文件加载失败（%d 个），共成功加载 %d 个", len(loadErrors), len(configs))
		}
	}

	return configs, nil
}
