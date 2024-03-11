package internal

import (
	"embed"
	"io"
	"io/fs"
	"text/template"

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

func (t *Template) Execute(w io.Writer, page *Page) error {
	return t.std.Execute(w, page)
}
