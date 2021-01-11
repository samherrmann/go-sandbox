package producer

import (
	"encoding/json"
)

func New() *Producer {
	return new(Producer)
}

type Producer struct{}

func (p *Producer) ReadJSON() ([]byte, error) {
	return json.Marshal(new(Data))
}

type Data struct {
	A string `json:"a"`
	B int    `json:"b"`
}
