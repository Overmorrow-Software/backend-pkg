package config

import (
	"errors"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

var ErrUndefinedConfigType = errors.New(`undefined config type (set CONFIG_TYPE="json" or "env")`)

type Parser[T any] interface {
	Parse() (*T, error)
}

type configType struct {
	CType string `env:"CONFIG_TYPE"`
}

func (t configType) isJSON() bool { return strings.ToLower(t.CType) == "json" }
func (t configType) isEnv() bool  { return strings.ToLower(t.CType) == "env" }

func New[T any](path string) (*T, error) {
	_ = godotenv.Load()

	t, err := env.ParseAs[configType]()
	if err != nil {
		return nil, err
	}

	if t.isJSON() {
		return NewJSON[T](path).Parse()
	}

	if t.isEnv() {
		return NewEnv[T]().Parse()
	}

	return nil, ErrUndefinedConfigType
}
