package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Database Database `yaml:"database" json:"database"`
	Settings Settings `yaml:"settings" json:"settings"`
}

type Database struct {
	Type     string `yaml:"type" json:"type"`
	Address  string `yaml:"address" json:"address"`
	Port     string `yaml:"port" json:"port"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	DBName   string `yaml:"dbname" json:"dbname"`
}

type Settings struct {
	SecretKey           string `yaml:"secretkey" json:"secretkey"`
	AgreementLetterHtml string `yaml:"agreementletterhtml" json:"agreementletterhtml"`
}

func New() (*Config, error) {
	cfg := &Config{}
	if err := LoadConfiguration(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func LoadConfiguration(bindCfg interface{}) error {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return err
		} else {
			return err
		}
	}
	if err := viper.Unmarshal(bindCfg); err != nil {
		return err
	}
	return nil
}
