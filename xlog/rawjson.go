package xlog

import (
	"encoding/json"

	"github.com/rs/zerolog"
)

var nullBytes = []byte("null")

func RawJSON(key string, rawJson []byte) func(e *zerolog.Event) {
	return func(e *zerolog.Event) {
		if len(rawJson) == 0 {
			e.RawJSON(key, nullBytes)
			return
		}
		if !IsJSON(rawJson) {
			e.Str(key, string(rawJson))
			return
		}
		e.RawJSON(key, rawJson)
	}
}

func IsJSON(rawJson []byte) bool {
	var js json.RawMessage
	return json.Unmarshal(rawJson, &js) == nil
}
