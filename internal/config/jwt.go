package config

type Jwt struct {
	ATPrivateKey string `mapstructure:"atPrivateKey"`
	ATExpiration string `mapstructure:"atExpiration"`
}
