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
	Error  string
	Data   any
}

func (v *View) AddStyleSheet(fs fs.FS, name string) error {
	f, err := fs.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	styles, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	v.Styles = template.CSS(fmt.Sprintf("%s\n%s", v.Styles, styles))
	return nil
}
