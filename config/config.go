package config

import "github.com/spf13/viper"

// Resource:
// https://dev.to/techschoolguru/load-config-from-file-environment-variables-in-golang-with-viper-2j2d

// Config - Hold configuration values from environment variables
type Config struct {
	CSV_FILENAME    string `mapstructure:"CSV_FILENAME"`
	POKEMON_API_URL string `mapstructure:"POKEMON_API_URL"`
	// DEFAULT_FILTER_NUM_WORKERS      int    `mapstructure:"DEFAULT_NUM_WORKERS"`
	// DEFAULT_FILTER_ITEMS            int    `mapstructure:"DEFAULT_FILTER_ITEMS"`
	// DEFAULT_FILTER_ITEMS_PER_WORKER int    `mapstructure:"DEFAULT_FILTER_ITEMS_PER_WORKER"`
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
