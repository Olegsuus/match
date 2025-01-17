package config

import (
	"github.com/spf13/viper"
)

type MongoSetting struct {
	Uri    string `mapstructure:"uri"    yaml:"uri"`
	DBName string `mapstructure:"dbname" yaml:"dbname"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret" yaml:"secret"`
}

type ServerConfig struct {
	Port string `mapstructure:"port" yaml:"port"`
}

type Log struct {
	LogFile string `mapstructure:"log_file" yaml:"log_file"`
}

type MovieAPI struct {
	Key string `mapstructure:"key" yaml:"key"`
	Url string `mapstructure:"url" yaml:"url"`
}

type Config struct {
	Env       string       `mapstructure:"env"       yaml:"env"`
	Mongo     MongoSetting `mapstructure:"mongo"     yaml:"mongo"`
	JWT       JWTConfig    `mapstructure:"jwt"       yaml:"jwt"`
	MoviesApi MovieAPI     `mapstructure:"movie_api" yaml:"movie_api"`
	Server    ServerConfig `mapstructure:"server"    yaml:"server"`
	Log       Log          `mapstructure:"log"       yaml:"log"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("local")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config/")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
