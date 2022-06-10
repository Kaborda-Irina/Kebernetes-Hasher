package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Port         int `mapstructure:"port"`
	PostgreSQL   `mapstructure:"postgres"`
	LoggerConfig `mapstructure:"logger"`
}

type PostgreSQL struct {
	URI      string `mapstructure:"uri"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database_name"`
}
type LoggerConfig struct {
	Level       int    `mapstructure:"level"`
	InfoLogFile string `mapstructure:"info_log_file"`
}

func LoadConfig() (config *Config, logger *logrus.Logger, err error) {
	// Initialize properties config
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	logger = InitLogger(&config.LoggerConfig)
	return config, logger, err
}
