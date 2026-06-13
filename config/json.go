package config

import (
	"fmt"
	"os"

	"github.com/goccy/go-json"
)

type JSONParser[T any] struct {
	path string
}

func NewJSON[T any](path string) Parser[T] {
	return &JSONParser[T]{path: path}
}

func (j *JSONParser[T]) Parse() (*T, error) {
	data, err := os.ReadFile(j.path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	var cfg T
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return &cfg, nil
}
