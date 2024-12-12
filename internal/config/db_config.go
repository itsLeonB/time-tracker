package config

type DB struct {
	Host     string
	User     string
	Password string
	Name     string
	Port     string
	SslMode  string
	TimeZone string
}

func loadDBConfig() *DB {
	return &DB{
		Host:     getRequiredEnv("DB_HOST"),
		User:     getRequiredEnv("DB_USER"),
		Password: getRequiredEnv("DB_PASSWORD"),
		Name:     getRequiredEnv("DB_NAME"),
		Port:     getRequiredEnv("DB_PORT"),
		SslMode:  getRequiredEnv("DB_SSL_MODE"),
		TimeZone: getRequiredEnv("DB_TIMEZONE"),
	}
}
