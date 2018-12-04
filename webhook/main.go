package main

import (
	"regexp"
	"strings"

	docopt "github.com/docopt/docopt-go"
	"github.com/fatih/structs"
	"github.com/kovetskiy/lorg"
	"github.com/kovetskiy/toml"
	karma "github.com/reconquest/karma-go"
	chat "github.com/zarplata/chattix/chat"
	"github.com/zarplata/chattix/context"
)

const (
	defaultAction     = "ACK"
	defaultActionType = "button"
	severityProblem   = "PROBLEM"

	messengerMattermost = "mattermost"
	messengerSlack      = "slack"
)

var (
	logger           *lorg.Log
	version          = "[manual build]"
	definedMessenger = messengerMattermost
	configPath       = "/etc/chattix/zabbix-to-" + definedMessenger + ".conf"
	eventIDExists    = false
	usage            = "zabbix-to-" + definedMessenger + " " + version + `

Usage:
  zabbix-to-` + definedMessenger + `  <channel> <severity> <message>

Options:
  <channel>   Channel in ` + definedMessenger + ` where message will be placed.  
                                                                   
  <severity>  Severity of event. Possible values are: OK or PROBLEM
                                                                   
  <message>   Message from Zabbix                                  
`
)

func main() {
	chatChooser := map[string]func() chat.Message{
		messengerMattermost: chat.NewMattermostMessage,
		messengerSlack:      chat.NewSlackMessage,
	}

	destiny := karma.Describe(
		"method", "main",
	).Describe(
		"version", version,
	).Describe(
		"chat type", definedMessenger,
	)

	logger = lorg.NewLog()
	conf := &config{}

	_, err := toml.DecodeFile(configPath, conf)
	if err != nil {
		logger.Fatal(
			destiny.Format(
				err,
				"can't read config file %s",
				configPath,
			),
		)
	}

	eventIDPattern, err := regexp.Compile(conf.EventIDRegexp)
	if err != nil {
		logger.Fatal(destiny.Reason(err))
	}

	args, err := docopt.Parse(usage, nil, true, version, false)
	if err != nil {
		logger.Fatal(destiny.Format(err, "can't parse args"))
	}

	channel, severity, message := parseArgs(args)

	eventID := ""
	fullEventIDMessage := ""

	matches := eventIDPattern.FindStringSubmatch(message)
	if len(matches) < 2 {
		logger.Warning(
			destiny.Describe(
				"message", message,
			).Describe(
				"error", err,
			).Reason(
				"can't find Event ID",
			),
		)
	} else {

		fullEventIDMessage = matches[0]
		eventID = matches[1]
		eventIDExists = true
	}

	request := chatChooser[definedMessenger]()

	icon := conf.getIconURL(severity)
	color := conf.getColor(severity)

	request.SetChannel(channel)
	request.SetIcon(icon)
	request.SetUsername(conf.Messengers[definedMessenger].MessengerUsername)

	attachment := request.CreateAttachment(message, color)
	attachment.SetTitle(severity)

	if eventIDExists {
		attachment.AddField(false, "Event ID", eventID)
		attachment.SetText(
			strings.Replace(message, fullEventIDMessage, "", -1),
		)
	}

	if severity != severityProblem {
		err = request.SendRequest(
			conf.Messengers[definedMessenger].MessengerAPIURL,
			conf.Messengers[definedMessenger].MessengerAPIToken,
		)

		if err != nil {
			logger.Fatal(
				destiny.Describe(
					"error", err,
				).Reason(
					"can't send message to chat",
				),
			)
		}

		return
	}

	if definedMessenger == messengerMattermost {

		actionContext := context.ContextActionACK{
			EventID:  eventID,
			Action:   defaultAction,
			Severity: severity,
			Message:  strings.Replace(message, fullEventIDMessage, "", -1),
			Channel:  channel,
			Username: conf.Messengers[messengerMattermost].MessengerUsername,
			IconURL:  icon,
		}

		attachment.AddAction(
			defaultAction,
			conf.Actions[defaultAction].ActionURL,
			defaultActionType,
			structs.Map(actionContext),
		)
	}

	if definedMessenger == messengerSlack {
		attachment.AddAction(
			defaultAction,
			defaultAction,
			defaultActionType,
			eventID,
		)
	}

	err = request.SendRequest(
		conf.Messengers[definedMessenger].MessengerAPIURL,
		conf.Messengers[definedMessenger].MessengerAPIToken,
	)
	if err != nil {
		logger.Fatal(
			destiny.Describe(
				"error", err,
			).Reason(
				"can't send message to chat",
			),
		)
	}

}

func parseArgs(
	args map[string]interface{},
) (channel, severity, message string) {

	if args["<channel>"] != nil {
		channel = args["<channel>"].(string)
	}

	if args["<severity>"] != nil {
		severity = args["<severity>"].(string)
	}

	if args["<message>"] != nil {
		message = args["<message>"].(string)
	}

	return
}
