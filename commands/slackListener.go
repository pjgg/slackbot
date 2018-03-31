package commands

import (
	"github.com/pjgg/slackbot/configuration"
	"github.com/pjgg/slackbot/connectors"
	"github.com/spf13/cobra"
)

var (
	listenerCmd = &cobra.Command{
		Use: "slack-listener",
		Run: slackListenerHandler,
	}
)

func init() {
	BaseCmd.AddCommand(listenerCmd)
}

func slackListenerHandler(cmd *cobra.Command, args []string) {
	slackConnector := connectors.Instance(configuration.ConfigurationManagerInstance.SlackToken)
	slackConnector.SlackBotListener()
}
