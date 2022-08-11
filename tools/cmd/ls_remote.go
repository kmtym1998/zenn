package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"kmtym1998/zenn-tools/service"

	"github.com/spf13/cobra"
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

		rg1 := regexp.MustCompile("<title.*?</title>")
		rg2 := regexp.MustCompile("<.*?>")

		var displayInfoList []map[string]string
		for _, f := range fileNames {
			if filepath.Ext(f) != ".md" {
				continue
			}

			metadata, err := service.ParseMDMetadata(dirPath + "/" + f)
			if err != nil {
				panic(err)
			}

			baseName := filepath.Base(f[:len(f)-len(filepath.Ext(f))])
			resp, err := service.SendRequest(http.MethodGet, "https://zenn.dev/kmtym1998/articles/"+baseName, nil)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()

			b, err := io.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}

			match := rg1.FindAllString(string(b), -1)
			if err != nil {
				panic(err)
			}

			var title string
			if len(match) != 0 {
				title = rg2.ReplaceAllString(match[0], "")
			} else {
				title = "タイトル不明"
			}

			displayInfoList = append(
				displayInfoList,
				map[string]string{
					"status": resp.Status,
					"emoji":  metadata.Emoji,
					"title":  title,
					"url":    "https://zenn.dev/kmtym/articles/" + baseName,
				},
			)
		}

		for _, d := range displayInfoList {
			fmt.Println("=======================================")
			fmt.Println(d["status"])
			fmt.Printf(
				"%s %s\n",
				d["emoji"], d["title"],
			)
			fmt.Println(d["url"])
			fmt.Printf("=======================================\n")
		}
	},
}
