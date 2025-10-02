/*
Copyright © 2025 SHEN YI <1245332635@qq.com>
*/
package command

import (
	"fs/internal/types"
	"fs/internal/version"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/cheggaaa/pb/v3"
	"github.com/ngaut/log"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const (
	fsDir      = ".fs"     // 隐藏目录名
	envRootDir = "FS_HOME" // 主目录
)

var (
	fsRootDir        string           // fs程序主路径
	fsBinDir         string           // fs 二进制路径
	fsHostsConfigDir string           // 存放 host 配置文件的目录
	bar              *pb.ProgressBar  // 进度条
	cfg              *types.SshConfig // ssh 配置
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fs",
	Short: "⚡ 轻量级命令行文件同步工具，让 `scp` 成为过去式",
	Long: `fs 是一款专为开发者与运维人员设计的轻量级文件传输工具，旨在替代传统 scp 命令，提供更简洁、高效、智能的跨主机文件操作体验。

无需记忆复杂参数，一键上传下载，支持多环境切换与 Shell 自动补全，大幅提升远程文件管理效率`,
	Version: version.Version, // 来自 internal/version

}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// 定义需要加载当前选中配置的命令（白名单）
var requiresSelectedConfig = []*cobra.Command{
	uploadCmd,
	downloadCmd,
	sshCmd,
	// 添加其他需要当前配置的命令
}

// init 初始化路径并注册配置加载
func init() {
	// 根据环境决定日志格式
	if os.Getenv("FS_DEBUG") != "" || os.Getenv("DEBUG") != "" {
		// 调试模式：显示文件名和行号
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	} else {
		// 生产模式：只显示时间和消息
		log.SetFlags(0)
	}

	// 1. 初始化根目录
	if home := os.Getenv(envRootDir); home != "" {
		fsRootDir = filepath.Join(home, fsDir)
	} else {
		userHome, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("❌ 无法获取用户主目录，且 %s 未设置: %v", envRootDir, err)
		}
		fsRootDir = filepath.Join(userHome, fsDir)
	}

	// 2. 初始化子目录
	fsBinDir = filepath.Join(fsRootDir, "bin")
	fsHostsConfigDir = filepath.Join(fsRootDir, "hosts")

	// 3. 注册配置初始化（在命令执行前运行）
	cobra.OnInitialize(loadSelectedConfig)
}

// loadSelectedConfig 加载当前选中的 SSH 配置到全局变量 cfg
func loadSelectedConfig() {
	// 获取当前被调用的命令
	cmd, _, err := rootCmd.Find(os.Args[1:])
	if err != nil {
		// 处理错误（可选）
		return
	}

	// 检查当前命令是否需要加载配置
	if !containsCommand(requiresSelectedConfig, cmd) {
		return // 不需要，跳过初始化
	}

	// 1. 读取 current 文件，获取当前选中的主机名
	currentFile := getUsedConfigFile()
	data, err := os.ReadFile(currentFile)
	if err != nil {
		log.Fatalf("❌ 未找到当前选中的主机配置，请先使用 'fs use <name>' 进行设置: %v", err)
	}

	name := strings.TrimSpace(string(data))
	if name == "" {
		log.Fatal("❌ 当前选中主机名为空，请使用 'fs use <name>' 选择一个主机")
	}

	// 2. 构造配置文件路径
	configFile := filepath.Join(fsHostsConfigDir, name+".yaml")

	// 3. 检查文件是否存在
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		log.Fatalf("❌ 主机配置文件不存在: %s", configFile)
	}

	// 4. 读取并解析 YAML
	yamlData, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("❌ 读取配置文件失败 %s: %v", configFile, err)
	}

	// 5. 初始化 cfg 并解析
	if cfg == nil {
		cfg = &types.SshConfig{}
	}

	if err := yaml.Unmarshal(yamlData, cfg); err != nil {
		log.Fatalf("❌ 解析 YAML 配置失败: %v", err)
	}

	// 6. 补全 Name 字段
	if cfg.Name == "" {
		cfg.Name = name
	}

	// 可选：打印调试信息
	// log.Printf("✅ 已加载主机配置: %s (%s@%s:%d)", cfg.Name, cfg.User, cfg.Host, cfg.Port)
}

// containsCommand 检查 cmd 是否在需要加载配置的命令列表中
func containsCommand(cmds []*cobra.Command, cmd *cobra.Command) bool {
	for _, c := range cmds {
		if c == cmd {
			return true
		}
	}
	return false
}

func passThru(r io.Reader, total int64) io.Reader {
	// start new bar
	reader := io.LimitReader(r, total)

	tmpl := `{{counters . }}  {{ bar . "[" "=" ">" "_" "|"}} {{rtime . "%s ]"}} {{speed . "%s/s" | rndcolor }} {{percent . | green}}`
	bar = pb.ProgressBarTemplate(tmpl).Start64(total)
	//bar := pb.Full.Start64(total)
	bar.Set(pb.SIBytesPrefix, true)
	bar.SetMaxWidth(100)

	// set custom bar template
	//bar.SetTemplateString(myTemplate)

	// create proxy reader
	barReader := bar.NewProxyReader(reader)

	return barReader

}

func getUsedConfigFile() string {
	file := path.Join(fsRootDir, "current")
	return file
}
