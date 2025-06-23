package configs

import (
	"github.com/spf13/viper"
)

type Config struct {
	DB struct {
		Host     string
		Port     int
		User     string
		Password string
		Name     string
		SSLMode  string
	}
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AutomaticEnv() // Read from environment variables

	if err := viper.ReadInConfig(); err != nil {
		// It's okay if the config file doesn't exist, we can rely on env vars
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
