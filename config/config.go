package config

type Config struct {
	TablePath string
}

func New() (Config, error) {
	return Config{}, nil
}
