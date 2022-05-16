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
		JwtKey string `yaml:"jwtkey"`
	} `yaml:"auth"`
}

func Load() (*Config, error) {
	appConfig := &Config{}
	configFile := "local.yaml"
	if _, err := os.Stat(configFile); err != nil {
		log.Printf("could not find local.yaml in directory: %s", err)
		configFile = "config.yaml"
	}

	bytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Printf("error reading config file: %s", err)
		return nil, err
	}
	if err = yaml.Unmarshal(bytes, &appConfig); err != nil {
		log.Printf("error unmarshalling config: %s", err)
		return nil, err
	}

	return appConfig, nil

}
