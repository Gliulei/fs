/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"github.com/cheggaaa/pb/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var bar *pb.ProgressBar
var cfg *SshConfig

var cfgs map[string]*SshConfig

type SshConfig struct {
	Host        string `mapstructure:"host"`
	User        string `mapstructure:"user"`
	Password    string `mapstructure:"password"`
	Port        int    `mapstructure:"port"`
	UploadDir   string `mapstructure:"upload_dir"`
	DownloadDir string `mapstructure:"download_dir"`
}

var cfgFile string
var group string
var defaultUploadDir string
var defaultDownloadDir string

const fsDir = ".fs"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fs",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.fs/fs.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.Flags().StringVar(&group, "g", "", "choose group")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".fs" (without extension).
		viper.AddConfigPath(home + "/.fs")
		viper.SetConfigType("yaml")
		viper.SetConfigName("fs")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		//fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	if err := viper.Unmarshal(&cfgs); err != nil {
		fmt.Fprintln(os.Stderr, "Unmarshal err:", err.Error())
	}

	//fmt.Printf("%+v\n", cfgs["default"])
	checkSshErr(cfgs)

	//read in used config
	if len(group) < 1 {
		b, _ := ioutil.ReadFile(getUsedConfigFile())
		group = strings.TrimSpace(string(b))
	}
	if c, ok := cfgs[group]; ok {
		cfg = c
	}
}

func checkSshErr(cfgs map[string]*SshConfig) {
	for _, cfg := range cfgs {
		if cfg.User == "" || cfg.Password == "" {
			fmt.Println("[error]user or password is empty")
			os.Exit(1)
		}

		if cfg.Port == 0 {
			cfg.Port = 22
		}

		if cfg.UploadDir == "" {
			cfg.UploadDir = "/home/" + cfg.User
		}

		if cfg.DownloadDir == "" {
			cfg.DownloadDir = defaultDownloadDir
		}
	}

}

func establishScpClient() (scp.Client, error) {
	// we ignore the host key in this example, please change this if you use this library
	clientConfig, _ := auth.PasswordKey(cfg.User, cfg.Password, ssh.InsecureIgnoreHostKey())

	// For other authentication methods see ssh.ClientConfig and ssh.AuthMethod

	// Create a new SCP client
	client := scp.NewClient(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), &clientConfig)

	// Connect to the remote server
	err := client.Connect()
	if err != nil {
		fmt.Println("Couldn't establish a connection to the remote server ", err)
	}

	return client, err
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
	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	file := path.Join(home, fsDir, "use.txt")
	return file
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}
