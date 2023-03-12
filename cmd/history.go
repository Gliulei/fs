/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bufio"
	"fmt"
	"github.com/ngaut/log"
	"io"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

// historyCmd represents the history command
var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "history command",
	Long: `history command`,
	Run: func(cmd *cobra.Command, args []string) {
		file := GetHistoryFile()
		fi, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		r := bufio.NewReader(fi)
		i := 0
		for {
			line, err := r.ReadString('\n')
			if err != nil && err != io.EOF {
				log.Fatal(err)
			}

			if err == io.EOF {
				break
			}
			i++
			line = fmt.Sprintf("[%d] %s", i, line)
			fmt.Print(line)

		}
	},
}

func init() {
	rootCmd.AddCommand(historyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// historyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// historyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func GetHistoryFile() string {
	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	file := path.Join(home, fsDir, "history.txt")
	return file
}

func record(cmdLog []string) {
	cmdString := strings.Join(cmdLog, " ")
	cmdString = cmdString + "\n"
	file := GetHistoryFile()
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Error(err)
		return
	}

	_, err = f.WriteString(cmdString)

	if err != nil {
		log.Error(err)
	}
}
