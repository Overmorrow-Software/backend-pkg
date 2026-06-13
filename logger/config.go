package logger

type Config struct {
	Level string `json:"level" env:"LOG_LEVEL" envDefault:"info"`
	Env   string `json:"env"   env:"APP_ENV"   envDefault:"prod"`
}

func (c *Config) isDev() bool { return c.Env == "dev" || c.Env == "development" || c.Env == "local" }
