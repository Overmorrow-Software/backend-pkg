package config

import "github.com/caarlos0/env/v11"

type EnvParser[T any] struct{}

func NewEnv[T any]() Parser[T] {
	return &EnvParser[T]{}
}

func (e *EnvParser[T]) Parse() (*T, error) {
	cfg, err := env.ParseAs[T]()
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
