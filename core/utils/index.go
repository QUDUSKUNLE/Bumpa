package utils

import (
	"encoding/json"
	"log"
)

func MustJSON(v any) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

func LogInfo(msg string, args ...any) {
	log.Printf(msg, args...)
}

func LogError(msg string, args ...any) {
	log.Printf("ERROR: "+msg, args...)
}
