package view

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
)

type View struct {
	Title  string
	Path   string
	Styles template.CSS
	Data   any
}

func (p *View) AddStyleSheet(fs fs.FS, name string) error {
	f, err := fs.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	styles, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	p.Styles = template.CSS(fmt.Sprintf("%s\n%s", p.Styles, styles))
	return nil
}
