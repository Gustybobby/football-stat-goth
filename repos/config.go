package repos

import "fmt"

type Config struct {
	Host     string
	Port     string
	Password string
	User     string
	DBName   string
	SSLMode  string
}

func Dsn(config *Config) string {
	dsn := fmt.Sprintf(
		"host=%s port=%s password=%s user=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.Password, config.User, config.DBName, config.SSLMode,
	)
	return dsn
}
