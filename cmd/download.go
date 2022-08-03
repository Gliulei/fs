/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"github.com/cheggaaa/pb/v3/termutil"
	"github.com/ngaut/log"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
	"io"
	"os"
	"path"
	"time"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

		passThru := func(r io.Reader, total int64) io.Reader {
			width, _:= termutil.TerminalWidth()
			reader := io.LimitReader(r, total)

			p := mpb.New(
				mpb.WithWidth(width - 40),
				mpb.WithRefreshRate(180*time.Millisecond),
			)

			bar := p.New(total,
				mpb.BarStyle().Rbound("|"),
				mpb.PrependDecorators(
					decor.CountersKibiByte("% .2f / % .2f"),
				),
				mpb.AppendDecorators(
					decor.EwmaETA(decor.ET_STYLE_GO, 90),
					decor.Name(" ] "),
					decor.EwmaSpeed(decor.UnitKiB, "% .2f", 60),
				),
			)
			barReader := bar.ProxyReader(reader)

			return barReader

		}

		err = client.CopyFromRemotePassThru(context.Background(), f, remoteFile, passThru)

		if err != nil {
			log.Errorf("Error while copying file %s", err.Error())
			return
		}

		log.Infof("download %s success, file in %s",remoteFile, srcFile)
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
