package connectors

import (
	"bufio"
	"fmt"
	"math/rand"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/nlopes/slack"
	"github.com/pjgg/slackbot/configuration"
	"github.com/sirupsen/logrus"
)

// SlackConnector ...struct contains a slack client pointer and a real time protocol pointer
type SlackConnector struct {
	Client              *slack.Client
	RealTimeMsgProtocol *slack.RTM
	amountMsg           int64
}

type SlackConnectorBehavior interface {
	SlackBotListener()
}

var onceSlack sync.Once

// SlackConnectorInstance is a slack client singleton instance
var SlackConnectorInstance SlackConnector

// New ...given a slack bot token, returns slack connector single instance.
func Instance(token string) *SlackConnector {
	onceSlack.Do(func() {
		SlackConnectorInstance.Client = slack.New(token)
		SlackConnectorInstance.RealTimeMsgProtocol = SlackConnectorInstance.Client.NewRTM()

		go SlackConnectorInstance.RealTimeMsgProtocol.ManageConnection()
	})

	return &SlackConnectorInstance
}
func (sc *SlackConnector) retrieveCommandResponse(ev *slack.MessageEvent) (any bool, response string) {
	for action, resp := range configuration.ConfigurationManagerInstance.Commands {
		if action.MatchString(ev.Text) {
			any = true
			response = resp
			break
		}
	}

	return
}

func (sc *SlackConnector) SlackBotListener() {

	for {
		select {
		case msg := <-sc.RealTimeMsgProtocol.IncomingEvents:
			logrus.Info("Event Received: ")
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				logrus.Info("Connection counter:", ev.ConnectionCount)

			case *slack.MessageEvent:
				sc.amountMsg++
				// Like taxi driver movie!. 'are you talking to me?'
				if sc.isTalkingToMe(ev) {
					//cmd := exec.Command(os.Args[0], sc.getUserExecCommand(ev))
					//sc.replyCommandStd(cmd, ev)
					go func() {
						if any, resp := sc.retrieveCommandResponse(ev); any {
							sc.RealTimeMsgProtocol.SendMessage(sc.RealTimeMsgProtocol.NewOutgoingMessage(resp, ev.Channel))
						}
					}()
				} else {
					if !sc.isMeWhoIsTalking(ev) {
						sc.asyncRandomMsg(ev)
						inboundMsg := ev.Text
						for _, trigger := range configuration.ConfigurationManagerInstance.Triggers {
							if (trigger.MatchString(inboundMsg) || sc.amountMsg > 150) && makeAJoke() {
								go func() {
									sc.RealTimeMsgProtocol.SendMessage(sc.RealTimeMsgProtocol.NewOutgoingMessage(randomReply(), ev.Channel))
								}()
								sc.amountMsg = 0
							}
						}
					}
				}

			case *slack.RTMError:
				logrus.Error("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				logrus.Error("Invalid credentials")

			default:
				// do nothing
			}
		}
	}
}

func (sc *SlackConnector) asyncRandomMsg(ev *slack.MessageEvent) {
	go func() {
		randomTextInterval := randomNum(60, 120)
		select {
		case <-time.Tick(time.Duration(randomTextInterval) * time.Minute):
			if makeAJoke() {
				sc.RealTimeMsgProtocol.SendMessage(sc.RealTimeMsgProtocol.NewOutgoingMessage(randomReply(), ev.Channel))
			}
		}

	}()
}

func (sc *SlackConnector) getEventUserIDTag(info *slack.Info) string {
	return fmt.Sprintf("<@%s> ", info.User.ID)
}

func (sc *SlackConnector) getUserExecCommand(ev *slack.MessageEvent) string {
	info := sc.RealTimeMsgProtocol.GetInfo()
	return strings.Replace(ev.Text, sc.getEventUserIDTag(info), "", -1)
}

func (sc *SlackConnector) isTalkingToMe(ev *slack.MessageEvent) bool {
	info := sc.RealTimeMsgProtocol.GetInfo()
	return ev.User != info.User.ID && strings.HasPrefix(ev.Text, sc.getEventUserIDTag(info))
}

func (sc *SlackConnector) isMeWhoIsTalking(ev *slack.MessageEvent) bool {
	info := sc.RealTimeMsgProtocol.GetInfo()
	return ev.User == info.User.ID
}

func (sc *SlackConnector) replyCommandStd(cmd *exec.Cmd, ev *slack.MessageEvent) {
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		logrus.Error("Error creating StdoutPipe: %v", err)
		sc.RealTimeMsgProtocol.SendMessage(sc.RealTimeMsgProtocol.NewOutgoingMessage(err.Error(), ev.Channel))
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			sc.RealTimeMsgProtocol.SendMessage(sc.RealTimeMsgProtocol.NewOutgoingMessage(scanner.Text(), ev.Channel))
		}
	}()

	if err := cmd.Start(); err != nil {
		sc.RealTimeMsgProtocol.SendMessage(sc.RealTimeMsgProtocol.NewOutgoingMessage(err.Error(), ev.Channel))
		logrus.Error(err)
	}

	if err := cmd.Wait(); err != nil {
		logrus.Error(err)
	}
}

func randomReply() string {
	position := randomNum(0, len(configuration.ConfigurationManagerInstance.Replies))
	return configuration.ConfigurationManagerInstance.Replies[position]
}

func randomNum(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func makeAJoke() (required bool) {
	seed := randomNum(0, 100)
	if configuration.ConfigurationManagerInstance.JokeLevel <= seed {
		required = true
	}

	return
}
