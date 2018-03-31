package main

import (
	"os"

	"github.com/pjgg/slackbot/commands"
	"github.com/pjgg/slackbot/configuration"
	"github.com/spf13/viper"
)

func main() {
	configuration.New()
	commands.Execute()
}

func init() {

	viper.SetConfigName("config")
	configPath, exist := os.LookupEnv("CONFIG_PATH")
	if exist {
		viper.AddConfigPath(configPath)
	}
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

}
