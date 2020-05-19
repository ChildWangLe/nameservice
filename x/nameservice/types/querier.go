package types

import "encoding/json"

type QueryResResolve struct {
	Value string `json:"value"`
}

func (object QueryResResolve) String() string {
	return object.Value
}

type QueryResNames []string

func (object QueryResNames) String() string {
	raw, _ := json.Marshal(object)
	return string(raw)
}
