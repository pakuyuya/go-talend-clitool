package cmd

import (
	"fmt"
	"os"

	"../jobitem"
	"../util/fileutils"
	"../util/stringutils"
	"github.com/spf13/cobra"
)

type _Options struct {
	Source string
	Output string
	Format string
	//	Bundle bool // not implemented
}

var (
	o = &_Options{}
)

var gensqlCmd = &cobra.Command{
	Use:   "gensql",
	Short: "SQLファイルを抽出して出力します。",
	Run: func(cmd *cobra.Command, args []string) {
		runapp()
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.AddCommand(gensqlCmd)

	gensqlCmd.Flags().StringVarP(&o.Source, "source", "s", "", "必須。解析対象のファイルパス。（例：project/path/**/*.item")
	gensqlCmd.Flags().StringVarP(&o.Output, "output", "o", "{source}.{ext}", "出力ファイル名。（デフォルト：解析したファイル名.出力ファイルフォーマット")
	gensqlCmd.Flags().StringVarP(&o.Format, "format", "f", "json", "出力するファイルのフォーマット。")

	gensqlCmd.MarkFlagRequired("source")
}
func runapp() {
	if !validateOptions() {
		return
	}

	paths := fileutils.FindMatchPathes(o.Source)

	if len(paths) == 0 {
		fmt.Println("no file matched path like `%s`", o.Source)
	}

	jobs := make([]*jobitemInfo, 0, 0)
	for _, path := range paths {
		s, err := os.Stat(path)
		if err != nil {
			fmt.Println("%s is invalid file path.", path)
			continue
		}
		if s.IsDir() {
			continue
		}

		job, err := parseJobitemFile(path)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		jobs = append(jobs, job)
	}

}

type jobitemInfo struct {
	FilePath string
	XMLElem  *jobitem.TalendFile
}

func parseJobitemFile(path string) (*jobitemInfo, error) {

	fp, err := os.OpenFile(path, os.O_RDONLY, 0444)
	if err != nil {
		return nil, fmt.Errorf("path `%s` cannot open. [Reason] %s", path, err.Error())
	}

	talendFile, err := jobitem.Parse(fp)
	fp.Close()

	if err != nil {
		return nil, fmt.Errorf("failed to parse `%s` as XML. [Reason] %s", path, err.Error())
	}

	if talendFile.DefaultContext == "" || talendFile.JobType == "" {
		return nil, fmt.Errorf("XML file `%s` is invalid format as talend job file.", path)
	}

	return &jobitemInfo{path, talendFile}, nil
}

func validateOptions() bool {
	isvalid := true

	if stringutils.EqualsAny(o.Format, "json", "csv") {
		fmt.Println("output format `%s` is not supported.", o.Format)
		isvalid = false
	}

	return isvalid
}
