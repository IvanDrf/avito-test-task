package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" `
	LoggerLevel string `yaml:"level"`

	Host string `yaml:"host"`
	Port string `yaml:"port"`

	StoragePath string `yaml:"storage_path" `
}

const defaultPath = "./config/config.yaml"

func MustLoad() *Config {
	cfgPath := getPathFromFlag()
	if cfgPath == "" {
		cfgPath = defaultPath
	}

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		log.Fatal("cant find config file " + cfgPath)
	}

	cfg := new(Config)

	if err := cleanenv.ReadConfig(cfgPath, cfg); err != nil {
		log.Fatal("cant parse config file " + err.Error())
	}

	return cfg
}

func getPathFromFlag() string {
	cfgPath := ""

	flag.StringVar(&cfgPath, "cfg", "", "")
	flag.Parse()

	return cfgPath
}
