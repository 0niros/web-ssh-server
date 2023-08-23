package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

type AppConfig struct {
	Port     string `json:"port" yaml:"port"`
	LogLevel string `json:"level" yaml:"logLevel"`
	AuthKey  string `json:"AuthKey" yaml:"authKey"`
	RootPath string `json:"rootPath" yaml:"rootPath"`
}

var GlobalConfig = new(AppConfig)

func ParseConfig() {
	// 1. Declare config path.
	// 1.1 Default config path: ./config.yaml .
	configPath := "config.yaml"
	// 1.2 If config env exists, cover default.
	if configEnv := os.Getenv("WEB_SSH_CONFIG"); configEnv != "" {
		configPath = configEnv
	}

	// 2. Initial viper.
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("read config failed: %s \n", err))
	}

	// 3. Unmarshal config.
	if err := v.Unmarshal(&GlobalConfig); err != nil {
		logrus.Error("[Config] error unmarshal config: ", err)
	}

	logrus.Info("[Config] config parsed: ", GlobalConfig)
}
