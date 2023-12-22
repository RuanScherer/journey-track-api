package config

import "github.com/spf13/viper"

type AppConfig struct {
	DbDsn        string `mapstructure:"DB_DSN"`
	DbLogEnabled bool   `mapstructure:"DB_LOG_ENABLED"`
}

var config *AppConfig

func GetAppConfig() *AppConfig {
	if config == nil {
		config = &AppConfig{}
		config = loadConfig()
	}
	return config
}

func loadConfig() *AppConfig {
	viper.SetConfigType("env")
	viper.SetConfigName(".env")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	err = viper.Unmarshal(config)
	if err != nil {
		panic("failed to unmarshal config: " + err.Error())
	}
	return config
}
