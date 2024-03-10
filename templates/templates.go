package templates

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
)

//go:embed root.html
var embedded embed.FS

func New() *Templates {
	return &Templates{
		fs:       MergeFS(embedded),
		registry: make(map[string]*template.Template),
	}
}

type Templates struct {
	fs       *mergedFS
	registry map[string]*template.Template
}

func (tpls Templates) Add(fs fs.FS, name string) error {
	tpls.fs.Add(fs)
	if tpls.registry[name] != nil {
		return fmt.Errorf("template %q already exists", name)
	}
	tpl, err := template.ParseFS(tpls.fs, "root.html", name)
	if err != nil {
		return err
	}
	tpls.registry[name] = tpl
	return nil
}

func (tpls Templates) ExecuteTemplate(w io.Writer, name string, page Page) error {
	tpl := tpls.registry[name]
	if tpl == nil {
		return fmt.Errorf("template %q not found", name)
	}
	return tpl.Execute(w, page)
}

type Page struct {
	Title string
	Data  any
}
