package consumer

import (
	"log"
)

type JSONReader interface {
	ReadJSON() ([]byte, error)
}

func DoStuff(p JSONReader) {
	b, err := p.ReadJSON()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(b))
}
