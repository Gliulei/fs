/*
Copyright © 2025 SHEN YI <1245332635@qq.com>
*/

package utils

import (
	"fmt"
	"fs/internal/types"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v3"

	"github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
)

// GetAuthMethod 根据配置返回 SSH 认证方法
func GetAuthMethod(cfg *types.SshConfig) ([]ssh.AuthMethod, error) {
	var authMethods []ssh.AuthMethod

	// 1. 如果配置了私钥
	if cfg.PrivateKeyPath != "" {
		var keyData []byte
		var err error

		// 优先从文件读取
		if cfg.PrivateKeyPath != "" {
			path := ExpandHome(cfg.PrivateKeyPath) // 处理 ~ 符号
			keyData, err = os.ReadFile(path)
			if err != nil {
				return nil, fmt.Errorf("❌ 读取私钥文件失败 %s: %v", cfg.PrivateKeyPath, err)
			}
		}

		// 解密私钥（如果被加密）
		var signer ssh.Signer
		if cfg.Passphrase != "" {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(keyData, []byte(cfg.Passphrase))
			if err != nil {
				return nil, fmt.Errorf("❌ 解密私钥失败: %v", err)
			}
		} else {
			signer, err = ssh.ParsePrivateKey(keyData)
			if err != nil {
				return nil, fmt.Errorf("❌ 解析私钥失败: %v", err)
			}
		}

		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	// 2. 如果有密码，也加入密码认证（可选）
	if cfg.Password != "" {
		authMethods = append(authMethods, ssh.Password(cfg.Password))
	}

	if len(authMethods) == 0 {
		return nil, fmt.Errorf("❌ 未配置任何认证方式（请设置 private_key_path 或 password）")
	}

	return authMethods, nil
}

// LoadConfigByName 加载指定名称的主机配置
func LoadConfigByName(fsHostsConfigDir, name string) (*types.SshConfig, error) {
	configFile := filepath.Join(fsHostsConfigDir, name+".yaml")

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("配置文件不存在: %s", configFile)
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	cfg := &types.SshConfig{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("解析配置失败: %v", err)
	}

	// 补全 Name 字段
	if cfg.Name == "" {
		cfg.Name = name
	}

	return cfg, nil
}

// establishScpClient 建立 SCP 客户端（支持密钥认证）
func EstablishScpClient(cfg *types.SshConfig) (scp.Client, error) {
	// we ignore the host key in this example, please change this if you use this library
	var clientConfig ssh.ClientConfig
	var err error
	if cfg.PrivateKeyPath != "" {
		clientConfig, err = auth.PrivateKeyWithPassphrase(cfg.User, []byte(cfg.Passphrase), cfg.PrivateKeyPath, ssh.InsecureIgnoreHostKey())
	} else {
		clientConfig, err = auth.PasswordKey(cfg.User, cfg.Password, ssh.InsecureIgnoreHostKey())
	}

	// For other authentication methods see ssh.ClientConfig and ssh.AuthMethod

	// Create a new SCP client
	client := scp.NewClient(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), &clientConfig)

	// Connect to the remote server
	err = client.Connect()
	if err != nil {
		fmt.Println("Couldn't establish a connection to the remote server ", err)
	}

	return client, nil
}
