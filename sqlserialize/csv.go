package sqlserialize

import (
	"io"
	"strings"
)

type textOption struct {
	RowFormat string
	ApplyFunc func(string) string
}

type TextOption func(*textOption)

func WithRowFormat(f string) TextOption {
	return func(opt *textOption) {
		opt.RowFormat = f
	}
}
func WithApplyFunc(f func(string) string) TextOption {
	return func(opt *textOption) {
		opt.ApplyFunc = f
	}
}

func Csv(entry *SqlEntry, w io.Writer, options ...TextOption) error {
	opt := textOption{
		RowFormat: "\"@Tag1@_@Tag2@_@Tag3@\",\"@Sql@\"",
		ApplyFunc: func(s string) string { return s },
	}

	for _, o := range options {
		o(&opt)
	}

	row := opt.RowFormat

	row = strings.Replace(row, "@Tag1@", opt.ApplyFunc(entry.Tag1), -1)
	row = strings.Replace(row, "@Tag2@", opt.ApplyFunc(entry.Tag2), -1)
	row = strings.Replace(row, "@Tag3@", opt.ApplyFunc(entry.Tag3), -1)
	row = strings.Replace(row, "@Sql@", opt.ApplyFunc(entry.Sql), -1)

	_, err := w.Write([]byte(row))

	return err
}

func CsvAry(entries []*SqlEntry, w io.Writer, options ...TextOption) error {
	for _, entry := range entries {
		if err := Csv(entry, w, options...); err != nil {
			return err
		}
	}
	return nil
}
