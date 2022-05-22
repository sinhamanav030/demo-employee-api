package config

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Port int    `yaml:"port"`
		Host string `yaml:"host"`
		Cors string `yaml:"cors"`
	} `yaml:"server"`

	Database struct {
		Host     string `yaml:"host"`
		Name     string `yaml:"name"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Port     string `yaml:"port"`
	} `yaml:"database"`

	Auth struct {
		JwtKey    string `yaml:"jwtkey"`
		JwtKeyLen int    `yaml:"jwtkeysize"`
	} `yaml:"auth"`

	Pagination struct {
		PerPage     int `yaml:"perpage"`
		DefaultPage int `yaml:"defaultpage"`
	} `yaml:"pagination"`
}

func Load(logger *log.Logger) (*Config, error) {
	appConfig := &Config{}
	configFile := "local.yaml"
	if _, err := os.Stat(configFile); err != nil {
		logger.Printf("could not find local.yaml in directory: %s\n", err)
		configFile = "config.yaml"
	}

	bytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		logger.Printf("error reading config file: %s\n", err)
		return nil, err
	}
	if err = yaml.Unmarshal(bytes, &appConfig); err != nil {
		logger.Printf("error unmarshalling config: %s\n", err)
		return nil, err
	}

	return appConfig, nil

}
