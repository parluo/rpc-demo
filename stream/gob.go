package stream

import (
	"bytes"
	"encoding/gob"
)

type GobEncoder struct {
}

func (g *GobEncoder) Encode(data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (g *GobEncoder) Decode(data []byte, obj interface{}) error {
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return nil
}
