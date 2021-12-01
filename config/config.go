package config

import "github.com/spf13/viper"

// Config - Hold configuration values from environment variables
type Config struct {
	CSV_FILENAME    string `mapstructure:"CSV_FILENAME"`
	POKEMON_API_URL string `mapstructure:"POKEMON_API_URL"`
}

// New - Config factory method
func New(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
