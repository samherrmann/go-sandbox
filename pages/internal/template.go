package internal

import (
	"embed"
	"fmt"
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

type Page struct {
	Title  string
	Path   string
	Styles string
	Data   any
}

func (p *Page) AddStyleSheet(fs fs.FS, name string) error {
	f, err := fs.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	styles, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	p.Styles = fmt.Sprintf("%s\n%s", p.Styles, styles)
	return nil
}
