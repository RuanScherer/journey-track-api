package config

import "github.com/spf13/viper"

type AppConfig struct {
	Environment string `mapstructure:"ENVIRONMENT"`
	FrontendUrl string `mapstructure:"FRONTEND_URL"`

	DbDsn        string `mapstructure:"DB_DSN"`
	DbLogEnabled bool   `mapstructure:"DB_LOG_ENABLED"`

	RestApiPort uint `mapstructure:"REST_API_PORT"`

	JwtSecret string `mapstructure:"JWT_SECRET"`

	EmailSmtpHost     string `mapstructure:"EMAIL_SMTP_HOST"`
	EmailSmtpPort     uint   `mapstructure:"EMAIL_SMTP_PORT"`
	EmailSmtpUsername string `mapstructure:"EMAIL_SMTP_USERNAME"`
	EmailSmtpPassword string `mapstructure:"EMAIL_SMTP_PASSWORD"`
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
	viper.AddConfigPath("/app/")

	viper.SetDefault("DB_LOG_ENABLED", false)
	viper.SetDefault("REST_API_PORT", 3000)

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
