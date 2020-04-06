package main

import (
	"encoding/json"
	"github.com/mattfenwick/upanddowntheriver/pkg/game"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	LogLevel string
	Host     string
	Port     int
}

// GetLogLevel ...
func (config *Config) GetLogLevel() (log.Level, error) {
	return log.ParseLevel(config.LogLevel)
}

// GetConfig ...
func GetConfig(configPath string) (*Config, error) {
	var config *Config

	viper.SetConfigFile(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to ReadInConfig at %s", configPath)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal config at %s", configPath)
	}

	return config, nil
}

func doOrDie(err error) {
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func main() {
	configPath := os.Args[1]
	config, err := GetConfig(configPath)
	doOrDie(err)

	logLevel, err := config.GetLogLevel()
	doOrDie(err)

	log.SetLevel(logLevel)

	client := game.NewClient(config.Host, config.Port)
	log.Infof("client: %+v", client)

	mod, err := client.GetModel()
	doOrDie(err)
	//modBytes, err := json.MarshalIndent(mod, "", "  ")
	//doOrDie(err)
	log.Infof("model: \n%s\n\n\n", mod)

	myModel, err := client.GetMyModel("abc")
	doOrDie(err)
	myModelBytes, err := json.MarshalIndent(myModel, "", "  ")
	doOrDie(err)
	log.Infof("my model: \n%s\n", string(myModelBytes))
}
