package config

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/yudai/pp"
)

type Config struct {
	Server    Server
	MongoDB   MongoConfig
	LogLevel  string
	SecretKey string
	Aud       string
	Iss       string
}

type Server struct {
	Port int
}

func New(configPath, configName string) (*Config, error) {
	v, err := readConfig(configPath, configName)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	config := &Config{}
	_ = v.Unmarshal(config)

	return config, nil
}

func readConfig(configPath, configName string) (*viper.Viper, error) {
	v := viper.New()
	v.AddConfigPath(configPath)
	v.SetConfigName(configName)
	err := v.ReadInConfig()

	return v, err
}

func (c *Config) Print() {
	pp.Println(c)
}
