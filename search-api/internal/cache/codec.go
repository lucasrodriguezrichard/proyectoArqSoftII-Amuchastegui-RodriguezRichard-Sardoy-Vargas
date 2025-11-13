package cache

import (
	"encoding/json"
	"errors"
)

// Codec serializes values so they can be stored in distributed caches.
type Codec interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte) (any, error)
}

// JSONCodec marshals values with encoding/json. New returns a pointer where the JSON will be unmarshaled.
type JSONCodec struct {
	New func() any
}

func (c JSONCodec) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (c JSONCodec) Unmarshal(data []byte) (any, error) {
	if c.New == nil {
		return nil, errors.New("json codec: missing constructor")
	}
	target := c.New()
	if err := json.Unmarshal(data, target); err != nil {
		return nil, err
	}
	return target, nil
}
