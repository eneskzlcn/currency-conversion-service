package config

type Jwt struct {
	ATPrivateKey        string `mapstructure:"atPrivateKey"`
	ATExpirationSeconds int    `mapstructure:"atExpirationSec"`
}
