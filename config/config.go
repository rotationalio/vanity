package config

const Prefix = "vanity"

type Config struct{}

func New() (Config, error) {
	return Config{}, nil
}
