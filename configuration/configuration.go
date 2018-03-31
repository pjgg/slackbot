package configuration

import (
	"fmt"
	"regexp"
	"strconv"
	"sync"

	"github.com/spf13/viper"
)

type ConfigurationManager struct {
	SlackToken string
	Replies    []string
	Triggers   []*regexp.Regexp
	JokeLevel  int
	Commands   map[*regexp.Regexp]string
}

var once sync.Once
var ConfigurationManagerInstance *ConfigurationManager

func New() *ConfigurationManager {

	once.Do(func() {
		viper.BindEnv("slack.token", "SLACK_TOKEN")
		slackToken := viper.GetString("slack.token")
		randomReplies := viper.GetStringSlice("easteregg.random-replies")
		triggersList := viper.GetStringSlice("easteregg.triggers")

		viper.BindEnv("easteregg.jokeLevel", "JOKE_LEVEL")
		jokeLevel, _ := strconv.Atoi(viper.GetString("easteregg.jokeLevel"))
		t := make([]*regexp.Regexp, len(triggersList))
		for k, v := range triggersList {
			t[k] = regexp.MustCompile(v)
		}

		commands := viper.GetStringMapString("help")
		c := make(map[*regexp.Regexp]string)
		for action, response := range commands {
			c[regexp.MustCompile(action)] = response
		}

		ConfigurationManagerInstance = &ConfigurationManager{
			SlackToken: slackToken,
			Replies:    randomReplies,
			Triggers:   t,
			JokeLevel:  jokeLevel,
			Commands:   c,
		}

		fmt.Println("Token" + slackToken)
		fmt.Println("random-replies size" + strconv.Itoa(len(randomReplies)))
	})

	return ConfigurationManagerInstance
}
