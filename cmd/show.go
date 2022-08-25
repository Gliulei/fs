/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "查看配置文件",
	Long: `查看配置文件`,
	Run: func(cmd *cobra.Command, args []string) {
		b0,_ := ioutil.ReadFile(getUsedConfigFile())
		fmt.Println(fmt.Sprintf("key `%s` is used", string(b0)))

		file := viper.ConfigFileUsed()
		b,_ := ioutil.ReadFile(file)
		fmt.Println(string(b))
	},
}

func init() {
	rootCmd.AddCommand(showCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
