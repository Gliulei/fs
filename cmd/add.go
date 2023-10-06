/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/ngaut/log"
	"github.com/spf13/viper"
	"os"
	"path"
	"strconv"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add command",
	Long:  `add command`,
	Run: func(cmd *cobra.Command, args []string) {
		user := cmd.Flag("user").Value
		pwd := cmd.Flag("pwd").Value
		host := cmd.Flag("host").Value
		upload := cmd.Flag("ud").Value
		download := cmd.Flag("dd").Value

		portVal := cmd.Flag("port").Value
		port, _ := strconv.Atoi(portVal.String())
		config := SshConfig{
			User:        user.String(),
			Password:    pwd.String(),
			Host:        host.String(),
			Port:        port,
			UploadDir:   upload.String(),
			DownloadDir: download.String(),
		}

		instance := fmt.Sprintf("%s_%d", host, port)
		viper.Set(instance, config)
		viper.WriteConfig()
		log.Infof("%s add called", instance)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	home, _ := os.UserHomeDir()
	defaultUploadDir = home
	defaultDownloadDir = path.Join(home, "download")

	addCmd.Flags().StringP("user", "u", "", "username")
	addCmd.Flags().StringP("pwd", "p", "", "password")
	addCmd.Flags().StringP("host", "H", "", "choose group")
	addCmd.Flags().IntP("port", "P", 0, "port")
	addCmd.Flags().String("ud", defaultUploadDir, "upload file dir")
	addCmd.Flags().String("dd", defaultDownloadDir, "download file dir")
}
