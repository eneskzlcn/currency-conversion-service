package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Db       DB
	Server   Server
	Jwt      Jwt
	AppEnv   string
	Rabbitmq RabbitMQ
}

func New() *Config {
	config := &Config{}
	viper.AutomaticEnv()
	appEnv := viper.GetString("APP_ENV")

	if appEnv == "" {
		appEnv = "local"
	}
	config.AppEnv = appEnv
	if appEnv == "local" {
		viper.AddConfigPath(".config/")
		viper.SetConfigName("local")
		viper.SetConfigType("yaml")
		err := viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("fatal error on reading config file %w", err))
		}
	}
	config.Server = Server{
		Port: viper.GetString("SERVER_PORT"),
	}
	config.Jwt = Jwt{
		ATPrivateKey:        viper.GetString("JWT_AT_PRIVATE_KEY"),
		ATExpirationMinutes: viper.GetInt("JWT_AT_EXPIRATION_MIN"),
	}
	config.Db = DB{
		Username: viper.GetString("DB_USERNAME"),
		Password: viper.GetString("DB_PASSWORD"),
		DBName:   viper.GetString("DB_NAME"),
		Host:     viper.GetString("DB_HOST"),
		Port:     viper.GetString("DB_PORT"),
	}

	config.Rabbitmq = RabbitMQ{
		Username:               viper.GetString("RABBITMQ_USERNAME"),
		Password:               viper.GetString("RABBITMQ_PASSWORD"),
		Host:                   viper.GetString("RABBITMQ_HOST"),
		Port:                   viper.GetString("RABBITMQ_PORT"),
		CurrencyConvertedQueue: viper.GetString("RABBITMQ_CURRENCY_CONVERTED_QUEUE"),
	}
	return config
}
func (c *Config) Print() {
	fmt.Println(*c)
}
