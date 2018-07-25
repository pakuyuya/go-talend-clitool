package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"../jobitem"
	"../sqlserialize"
	"../util/fileutils"
	"../util/stringutils"
	"github.com/spf13/cobra"
)

type _Options struct {
	Source string
	Output string
	Format string
	Tag1   string
	Tag2   string
	Tag3   string
	Bundle bool
}

var (
	o = &_Options{}
)

const (
	JSON string = "json"
	CSV  string = "csv"
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
	gensqlCmd.Flags().StringVarP(&o.Output, "output", "o", "{source}.{ext}", "出力ファイル名。（デフォルト：解析したファイル名.拡張子")
	gensqlCmd.Flags().StringVarP(&o.Format, "format", "f", "json", "出力するファイルのフォーマット。")
	gensqlCmd.Flags().StringVarP(&o.Tag1, "tag1", "", "{source}", "出力ファイルのTag1に設定する内容のテンプレート")
	gensqlCmd.Flags().StringVarP(&o.Tag2, "tag2", "", "{uniquename}", "出力ファイルのTag2に設定する内容のテンプレート")
	gensqlCmd.Flags().StringVarP(&o.Tag3, "tag3", "", "", "出力ファイルのTag3に設定する内容のテンプレート")
	gensqlCmd.Flags().BoolVarP(&o.Bundle, "bundle", "b", false, "出力ファイルを1つに固めます。")

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

	if o.Bundle {
	} else {

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

func getGensqlEntries(path string, talendFile *jobitem.TalendFile) ([]*sqlserialize.SqlEntry, error) {
	entries := make([]*sqlserialize.SqlEntry, 0, 0)

	links, err := jobitem.GetNodeLinks(talendFile)

	if err != nil {
		return nil, err
	}

	for _, l := range links {
		t := jobitem.GetComponentType(&l.Node)

		sql := ""
		switch t {
		case jobitem.ComponentELTOutput:
			sql, _ = jobitem.TELTOutput2InsertSQL(l)
		case jobitem.ComponentDBRow:
			sql, _ = jobitem.DBRow2SQL(l)
		default:
			continue
		}

		uniquename, _ := jobitem.GetUniqueName(&l.Node)

		fctx := fmtContext{path, uniquename}

		e := &sqlserialize.SqlEntry{
			Sql:  sql,
			Tag1: fmtInContext(o.Tag1, fctx),
			Tag2: fmtInContext(o.Tag2, fctx),
			Tag3: fmtInContext(o.Tag3, fctx),
		}
		entries = append(entries, e)
	}
	return entries, nil
}

type fmtContext struct {
	Filename  string
	Component string
}

func fmtInContext(fmt string, ctx fmtContext) string {
	base := filepath.Base(ctx.Filename)
	idot := strings.LastIndex(base, ".")
	var fname, ext string
	if idot != -1 {
		fname = base[0:idot]
		ext = base[idot+1:]
	} else {
		fname = base
		ext = ""
	}

	s := strings.Replace(fmt, "{source}", fname, -1)
	s = strings.Replace(fmt, "{ext}", ext, -1)
	s = strings.Replace(fmt, "{component}", ctx.Component, -1)
	return s
}

func validateOptions() bool {
	isvalid := true

	if stringutils.EqualsAny(o.Format, JSON, CSV) {
		fmt.Println("output format `%s` is not supported.", o.Format)
		isvalid = false
	}

	return isvalid
}

func writeBundle(jobs []*jobitemInfo) error {

	bundledFilename := fmtInContext(o.Output, fmtContext{"bundle", "all"})
	fp, err := os.OpenFile(bundledFilename, os.O_WRONLY, 0666)

	if err != nil {
		return fmt.Errorf("try write to %s, but access denied. reason:%s", bundledFilename, err.Error())
	}

	entries := make([]*sqlserialize.SqlEntry, 0, 0)
	for _, job := range jobs {
		es, _ := getGensqlEntries(job.FilePath, job.XMLElem)
		entries = append(entries, es...)
	}

	switch o.Format {
	case JSON:
		err = sqlserialize.JsonAry(entries, fp)
	case CSV:
		err = sqlserialize.CsvAry(entries, fp)
	}
	fp.Close()

	if err != nil {
		return err
	}

	return nil
}

func writeEach(jobs []*jobitemInfo) error {
	for _, job := range jobs {
		entries, _ := getGensqlEntries(job.FilePath, job.XMLElem)

		basename := fileutils.Basename(job.FilePath)
		ext := filepath.Ext(job.FilePath)

		bundledFilename := fmtInContext(o.Output, fmtContext{basename, ext})
		fp, err := os.OpenFile(bundledFilename, os.O_WRONLY, 0666)

		if err != nil {
			// error, but continue routien
			fmt.Println("try write to %s, but access denied. reason:%s", bundledFilename, err.Error())
			continue
		}

		switch o.Format {
		case JSON:
			err = sqlserialize.JsonAry(entries, fp)
		case CSV:
			err = sqlserialize.CsvAry(entries, fp)
		}
		fp.Close()
	}

	return nil
}
