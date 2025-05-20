package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBHost            string `mapstructure:"DB_HOST"`
	DBPort            string `mapstructure:"DB_PORT"`
	DBUser            string `mapstructure:"DB_USER"`
	DBPassword        string `mapstructure:"DB_PASSWORD"`
	DBName            string `mapstructure:"DB_NAME"`
	APIPort           string `mapstructure:"API_PORT"`
	AgifyAPIURL       string `mapstructure:"AGIFY_API_URL"`
	GenderizeAPIURL   string `mapstructure:"/erize_API_URL"`
	NationalizeAPIURL string `mapstructure:"NATIONALIZE_API_URL"`
	HTTPClientTimeout string `mapstructure:"HTTP_CLIENT_TIMEOUT"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
