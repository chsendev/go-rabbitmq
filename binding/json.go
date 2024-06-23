package binding

import (
	"bytes"
	"encoding/json"
	"io"
)

type jsonBinding struct{}

func (jsonBinding) Name() string {
	return "json"
}

func (jsonBinding) BindBody(body []byte, obj any) error {
	return decodeJSON(bytes.NewBuffer(body), obj)
}

func decodeJSON(r io.Reader, obj any) error {
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return nil
}
