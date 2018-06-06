package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "go-talend-sqlextractor",
	Short: "go-talend-sqlextractor は Talendの各コンポーネントからSQLを抽出するツールです。",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("go-talend-sqlextractor v0.1")
	},
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
}
