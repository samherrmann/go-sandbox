package main

import (
	"log"
	"os"
	"path/filepath"
	"plugin"
)

func main() {
	dir, err := executableDir()
	if err != nil {
		log.Fatalln(err)
	}

	p, err := plugin.Open(filepath.Join(dir, "plugin.so"))
	if err != nil {
		log.Fatalln(err)
	}

	f, err := p.Lookup("SayHello")
	if err != nil {
		log.Fatalln(err)
	}
	f.(func())()
}

func executableDir() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(execPath), nil
}
