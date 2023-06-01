package rbac

import (
	jsoniter "github.com/json-iterator/go"
)

type codec []string

func (codec) Name() string {
	return "json"
}

func (c codec) Marshal(v interface{}) ([]byte, error) {
	return jsoniter.Config{TagKey: []string(c)}.Froze().Marshal(v)
}

func (c codec) Unmarshal(data []byte, v interface{}) error {
	return jsoniter.Config{TagKey: []string(c)}.Froze().Unmarshal(data, v)
}
