package configs

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
	"os"
)

// once lv-fetch configs is imported, this init() function will be run and fill the Config Struct
func init() {
	readFile(&cfg)
	readEnv(&cfg)
}

// GetConfig return an instance of Config Struct
func GetConfig() *Config {
	return &cfg
}

// readFile reads config.yml file at first place of call.
func readFile(cfg *Config) {
	f, err := os.Open("./configs/config.yml")
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(f)

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}

// readEnv is for reading production level configs from os env.
func readEnv(cfg *Config) {
	err := envconfig.Process("", cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}
