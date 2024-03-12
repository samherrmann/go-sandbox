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

func ParseTemplate(fsys fs.FS, name string) (*Template, error) {
	fsys = fsutil.MergeFS(embedded, fsys)

	tpl, err := template.ParseFS(fsys, "root.html", name)
	if err != nil {
		return nil, err
	}
	return &Template{std: tpl}, nil
}

type Template struct {
	std *template.Template
}

func (t *Template) Execute(w io.Writer, view *View) error {
	return t.std.Execute(w, view)
}
