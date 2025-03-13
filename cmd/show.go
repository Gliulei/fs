/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "查看配置文件",
	Long:  `查看配置文件`,
	Run: func(cmd *cobra.Command, args []string) {
		usedId, _ := ioutil.ReadFile(getUsedConfigFile())
		// 提取所有的 key 到一个切片中
		keys := make([]string, 0, len(cfgs))
		for k := range cfgs {
			keys = append(keys, k)
		}
		// 对 key 切片进行排序
		sort.Strings(keys)

		table := tablewriter.NewWriter(os.Stdout)
		table.SetAlignment(tablewriter.ALIGN_CENTER)
		table.SetHeader([]string{"id", "host","port"})
		for _, key := range keys {
			cfg := cfgs[key]
			row := []string{key, cfg.Host, strconv.Itoa(cfg.Port)}
			if key == string(usedId) {
				table.Rich(row, []tablewriter.Colors{tablewriter.Colors{tablewriter.Bold, tablewriter.BgRedColor}, tablewriter.Colors{tablewriter.Bold, tablewriter.BgRedColor}, tablewriter.Colors{tablewriter.Bold, tablewriter.BgRedColor}})
			} else {
				table.Append(row)
			}
		}
		table.Render() // Send output
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
