package services

import (
	"encoding/json"
	"github.com/QUDUSKUNLE/Bumpa/core/ports"
)

type ServicesHandler struct {
	ports ports.RepositoryPorts
}

func NewServiceAdapter(repositoryPort ports.RepositoryPorts) *ServicesHandler {
	return &ServicesHandler{
		ports: repositoryPort,
	}
}

func MustJSON(v any) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}
