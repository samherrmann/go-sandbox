package internal

import (
	"fmt"
	"io"
	"io/fs"
)

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
