package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// RootCmd is root command
var (
	RootCmd = &cobra.Command{
		Use:   "zenn-tool",
		Short: "自前 zenn CLI",
		Long:  "自前 zenn CLI",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("サブコマンドを入力してください")
		},
	}
)

// コマンド実行時に最初に呼ばれる初期化処理
func init() {
	// フラグの定義
	// 第1引数: フラグ名、第2引数: 省略したフラグ名
	// 第3引数: デフォルト値、第4引数: フラグの説明
}
