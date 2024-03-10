package templates

import (
	"fmt"
	"io"
	"text/template"
)

func New() *Templates {
	return &Templates{
		registry: make(map[string]*template.Template),
	}
}

type Templates struct {
	registry map[string]*template.Template
}

func (tpls Templates) Add(name string) error {
	if tpls.registry[name] != nil {
		return fmt.Errorf("template %q already exists", name)
	}
	tpl, err := template.ParseFiles("templates/root.html", name)
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
