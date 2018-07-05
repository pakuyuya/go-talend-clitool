package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type _Options struct {
	Source string
	Output string
}

var (
	o = &_Options{}
)

var gensqlCmd = &cobra.Command{
	Use:   "gensql",
	Short: "SQLファイルを抽出して出力します。",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("go-talend-sqlextractor v0.1")
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.AddCommand(gensqlCmd)

	gensqlCmd.Flags().StringVarP(&o.Source, "source", "s", "", "必須。対象のファイルパスです。")
}

func initConfig() {
}
