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

func ParseTemplate(fsys fs.FS, htmlFilename string, cssFilename string) (*Template, error) {
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

	return &Template{std: tpl}, nil
}

type Template struct {
	std *template.Template
}

func (t *Template) Execute(w io.Writer, view *View) error {
	return t.std.ExecuteTemplate(w, "root", view)
}
