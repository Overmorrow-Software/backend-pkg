package db

type Config struct {
	Host     string `json:"host" env:"DB_HOST"`
	User     string `json:"user" env:"DB_USER"`
	Password string `json:"password" env:"DB_PASSWORD"`
	DBName   string `json:"db_name" env:"DB_NAME"`
	SSLMode  string `json:"ssl_mode" env:"DB_SSL_MODE"`
	Port     int    `json:"port" env:"DB_PORT"`
}
