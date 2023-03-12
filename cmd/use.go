/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/ngaut/log"
	"github.com/spf13/cobra"
	"os"
)

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use",
	Short: "which group to use",
	Long:  "which group to use",
	Run: func(cmd *cobra.Command, args []string) {
		groupName := args[0]
		if len(groupName) < 1 {
			log.Fatal("group can't is empty")
		}
		// Open a file
		f, err := os.OpenFile(getUsedConfigFile(), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
		if err != nil {
			log.Error(err)
			return
		}

		f.Write([]byte(groupName))

		fmt.Println("set success!")
	},
}

func init() {
	rootCmd.AddCommand(useCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// useCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// useCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
