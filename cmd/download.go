/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"github.com/ngaut/log"
	"github.com/spf13/cobra"
	"os"
	"path"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "A brief description of your command",
	Long: `download file`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Error("No upload file specified")
			return
		}
		client, err := establishScpClient()
		if err != nil {
			log.Error(err)
			return
		}
		// Close client connection after the file has been copied
		defer client.Close()

		remoteFile := args[0]
		if !path.IsAbs(remoteFile) {
			remoteFile = path.Join(cfg.UploadDir, remoteFile)
		}

		filename := path.Base(remoteFile)
		srcFile := path.Join(cfg.DownloadDir, filename)

		// Open a file
		f, err := os.OpenFile(srcFile, os.O_RDWR|os.O_CREATE, 0777)
		if err != nil {
			log.Error(err)
			return
		}

		// Close the file after it has been copied
		defer f.Close()

		//err = client.CopyFromRemote(context.Background(), f, remoteFile)
		err = client.CopyFromRemotePassThru(context.Background(), f, remoteFile, passThru)

		if bar != nil {
			bar.Finish()
		}
		if err != nil {
			log.Errorf("Error while copying file %s", err.Error())
			return
		}

		//记录history
		cmdLog := []string{"fs", "download"}
		cmdLog = append(cmdLog, args...)
		record(cmdLog)
		//log.Infof("download %s success, file in %s",remoteFile, srcFile)
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
