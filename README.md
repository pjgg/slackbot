# Slackbot

This is an example about how to implement your own slackbot. If you want to test it, just create a bot in Slack and copy/paste your slackbot token in [config_example.yaml](https://github.com/pjgg/slackbot/blob/master/config_example.yaml) and in the makef[Makefile](https://github.com/pjgg/slackbot/blob/master/Makefile) file variable "SLACK_TOKEN"

## Config

- **slack.token**: This is your slackbot token.
- **help**: When you talk to your bot @botName you could configure some help command. For example, which is google url?
- **joneypot.jokeLevel**: This bot will make you yokes, based in some keyword defined in slackbot.triggers. This property will limit the number of jokes that the bot will make you. 10, means 10%. 
- **slackbot.triggers**: Key words that will trigger a joke.
- **slackbot.random-replies**: jokes database. 
