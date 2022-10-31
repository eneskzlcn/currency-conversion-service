package config

type DB struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbName"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
}
