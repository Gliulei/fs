/*
Copyright © 2025 SHEN YI <1245332635@qq.com>
*/

package types

// SshConfig 定义SSH配置结构
type SshConfig struct {
	Host               string `mapstructure:"host"`
	User               string `mapstructure:"username"`
	Password           string `mapstructure:"password"`
	Port               int    `mapstructure:"port"`
	DefaultUploadDir   string `mapstructure:"default_upload_dir"`
	DefaultDownloadDir string `mapstructure:"default_download_dir"`
	Name               string `mapstructure:"name"`
	PrivateKeyPath     string `yaml:"private_key_path" mapstructure:"private_key_path"` //私钥路径
	Passphrase         string `yaml:"passphrase" mapstructure:"passphrase"`             // 私钥加密密码
}
