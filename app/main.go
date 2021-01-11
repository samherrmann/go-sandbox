package main

import (
	"github.com/samherrmann/go-sanbox/app/consumer"
	"github.com/samherrmann/go-sanbox/app/producer"
)

func main() {
	p := producer.New()
	consumer.DoStuff(p)
}
