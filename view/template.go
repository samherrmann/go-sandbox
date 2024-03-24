package view

import (
	"embed"
	"html/template"
	"io"
	"io/fs"

	"github.com/samherrmann/go-sandbox/fsutil"
)

//go:embed root.html
var embedded embed.FS

func ParseTemplate[T any](fsys fs.FS, htmlFilename string, cssFilename string) (*Template[T], error) {
	fsys = fsutil.MergeFS(embedded, fsys)

	filenames := map[string]string{
		"body":   htmlFilename,
		"styles": cssFilename,
	}

	tpl, err := template.ParseFS(fsys, "root.html")
	if err != nil {
		return nil, err
	}

	for name, filename := range filenames {
		b, err := fs.ReadFile(fsys, filename)
		if err != nil {
			return nil, err
		}
		tpl, err = tpl.New(name).Parse(string(b))
		if err != nil {
			return nil, err
		}
	}

	return &Template[T]{std: tpl}, nil
}

type Template[T any] struct {
	std *template.Template
}

func (t *Template[T]) Execute(w io.Writer, view *ViewData[T]) error {
	return t.std.ExecuteTemplate(w, "root", view)
}
