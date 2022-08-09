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

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "upload file",
	Long: `upload file to dest host`,
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

		srcFile := args[0]
		if !path.IsAbs(srcFile) {
			srcFile = path.Join(cfg.DownloadDir, srcFile)
		}
		filename := path.Base(srcFile)
		// Open a file
		f, err := os.Open(srcFile)
		if err != nil {
			log.Error(err)
			return
		}

		// Close the file after it has been copied
		defer f.Close()

		// Finaly, copy the file over
		// Usage: CopyFromFile(context, file, remotePath, permission)

		// the context can be adjusted to provide time-outs or inherit from other contexts if this is embedded in a larger application.
		remoteFile := path.Join(cfg.UploadDir, filename)
		err = client.CopyFromFile(context.Background(), *f, remoteFile, "0655")

		if err != nil {
			log.Errorf("Error while copying file %s", err.Error())
			return
		}

		log.Infof("upload %s success, file in %s", srcFile, remoteFile)
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uploadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uploadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
