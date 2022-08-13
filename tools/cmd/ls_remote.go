package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"kmtym1998/zenn-tools/service"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func init() {
	RootCmd.AddCommand(lsCmd)

	lsCmd.Flags().StringP("articles-path", "p", "", "articles ディレクトリまでのパス")
}

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls-remote",
	Short: "zenn に投稿済みの記事リスト",
	Long:  "zenn に投稿済みの記事リスト",
	Run: func(cmd *cobra.Command, args []string) {
		dirPath, err := cmd.Flags().GetString("articles-path")
		if err != nil {
			log.Fatalf("error getting articles-path: %v", err)
		}
		if dirPath == "" {
			dirPath = "./articles"
		}

		articles, err := os.ReadDir(dirPath)
		if err != nil {
			log.Fatalf("%s を読み取ってエラー。フラグで読み取るディレクトリの指定ができます", dirPath)
		}

		var fileNames []string
		for _, info := range articles {
			if info.IsDir() {
				continue
			}

			fileNames = append(fileNames, info.Name())
		}

		var eg errgroup.Group
		var displayInfoList []map[string]string
		for _, f := range fileNames {
			f := f

			eg.Go(func() error {
				if filepath.Ext(f) != ".md" {
					return nil
				}

				metadata, err := service.ParseMDMetadata(dirPath + "/" + f)
				if err != nil {
					return err
				}

				baseName := filepath.Base(f[:len(f)-len(filepath.Ext(f))])

				status := color.RedString("未公開")
				if metadata.Published {
					status = color.GreenString("公開中")
				}

				displayInfoList = append(
					displayInfoList,
					map[string]string{
						"status": status,
						"emoji":  metadata.Emoji,
						"title":  metadata.Title,
						"url":    "https://zenn.dev/kmtym1998/articles/" + baseName,
					},
				)

				return nil
			})

			if err := eg.Wait(); err != nil {
				panic(err)
			}
		}

		for _, d := range displayInfoList {
			fmt.Println("=======================================")
			fmt.Println(d["status"])
			fmt.Printf("%s %s\n", d["emoji"], d["title"])
			fmt.Println(d["url"])
		}
	},
}
