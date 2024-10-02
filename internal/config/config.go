package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Addr string `yaml:"address" env-required:"true"`
}

// env-default:"production" -> for default settings
type Config struct {
	Env         string `yaml:"env" env:"ENV" env-required:"true"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

func MustLoad() *Config {
	// Loading the config path
	configPath := os.Getenv("CONFIG_PATH")
	// checking is it passed
	if configPath == "" {

		// if it is passed along with run command
		flags := flag.String("config", "", "Path to the config file")
		flag.Parse()

		configPath = *flags

		if configPath == "" {
			log.Fatal("Config path is not set")
		}
	}

	// if file does not exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("Config file does not exist.", configPath)
	}

	// Reading file
	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("Cannot read config file %s", err.Error())
	}

	return &cfg
}
