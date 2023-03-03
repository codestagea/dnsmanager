package config

import (
	"os"

	yaml "gopkg.in/yaml.v2"
)

type HttpConfig struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	ContextPath string `yaml:"context-path"`
	Token       string `yaml:"token"`
}

type Database struct {
	Driver string `yaml:"driver"`
	Dsn    string `yaml:"dsn"`
	Debug  bool   `yaml:"debug"`
}

type Redis struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
}

type Config struct {
	LogLevel   string     `yaml:"log_level"`
	HttpConfig HttpConfig `yaml:"http"`
	Envs       []string   `yaml:"envs"`
	Database   Database   `yaml:"database"`
	//Redis      Redis      `yaml:"redis"`
	Jwt Jwt `yaml:"jwt"`
	Rsa Rsa `yaml:"rsa"`
}

type Jwt struct {
	Secret    string `yaml:"secret"`
	Timeout   int    `yaml:"timeout"`
	CacheType string `yaml:"cacheType"`
}

type Rsa struct {
	PrivateKey string `yaml:"privateKey"`
}

func LoadConfig(filepath string) (*Config, error) {
	yamlFile, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	} else {
		conf := &Config{}

		err = yaml.Unmarshal(yamlFile, conf)
		if err != nil {
			return nil, err
		}
		return conf, nil
	}
}
