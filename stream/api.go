package stream

type Encoder interface {
	Encode(obj interface{}) ([]byte, error)
	Decode(data []byte, obj interface{}) error
}
