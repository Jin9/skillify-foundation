package codec

import (
	"bytes"
	"encoding/json"
)

type JSONCoder interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
}

var _ JSONCoder = (*jsonCoder)(nil)

type jsonCoder struct{}

func NewJSONCoder() JSONCoder {
	return jsonCoder{}
}

func (c jsonCoder) Marshal(v any) ([]byte, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c jsonCoder) Unmarshal(data []byte, v any) error {
	buf := bytes.NewReader(data)
	dec := json.NewDecoder(buf)
	return dec.Decode(v)
}
