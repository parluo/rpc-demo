package stream

import (
	"bytes"
	"encoding/json"
)

// todo 序列化float64的问题
type JsonEncoder struct {
}

func (j *JsonEncoder) Encode(obj interface{}) ([]byte, error) {
	return json.Marshal(obj)
}

func (j *JsonEncoder) Decode(data []byte, obj interface{}) error {
	return json.Unmarshal(data, obj)
}

func (j *JsonEncoder) Decode1(data []byte, obj interface{}) error {
	buf := json.NewDecoder(bytes.NewBuffer(data))
	buf.UseNumber()
	err := buf.Decode(obj)
	if err != nil {
		return err
	}
	return err
}
