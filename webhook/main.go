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

	chatMattermost = "mattermost"
	chatSlack      = "slack"
)

var (
	logger        *lorg.Log
	version       = "[manual build]"
	definedChat   = chatMattermost
	configPath    = "/etc/chattix/zabbix-to-" + definedChat + ".conf"
	eventIDExists = false
	usage         = "zabbix-to-" + definedChat + " " + version + `

Usage:
  zabbix-to-` + definedChat + `  <channel> <severity> <message>

Options:
  <channel>   Channel in ` + definedChat + ` where message will be placed.  
                                                                   
  <severity>  Severity of event. Possible values are: OK or PROBLEM
                                                                   
  <message>   Message from Zabbix                                  
`
)

func main() {
	chatChooser := map[string]func() chat.Message{
		chatMattermost: chat.NewMattermostMessage,
		chatSlack:      chat.NewSlackMessage,
	}

	destiny := karma.Describe(
		"method", "main",
	).Describe(
		"version", version,
	).Describe(
		"chat type", definedChat,
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

	request := chatChooser[definedChat]()

	icon := conf.getIconURL(severity)
	color := conf.getColor(severity)

	request.SetChannel(channel)
	request.SetIcon(icon)
	request.SetUsername(conf.Chats[definedChat].ChatUsername)

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
			conf.Chats[definedChat].ChatAPIURL,
			conf.Chats[definedChat].ChatAPIToken,
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

	if definedChat == chatMattermost {

		actionContext := context.ContextActionACK{
			EventID:  eventID,
			Action:   defaultAction,
			Severity: severity,
			Message:  strings.Replace(message, fullEventIDMessage, "", -1),
			Channel:  channel,
			Username: conf.Chats[chatMattermost].ChatUsername,
			IconURL:  icon,
		}

		attachment.AddAction(
			defaultAction,
			conf.Actions[defaultAction].ActionURL,
			defaultActionType,
			structs.Map(actionContext),
		)
	}

	if definedChat == chatSlack {
		attachment.AddAction(
			defaultAction,
			defaultAction,
			defaultActionType,
			eventID,
		)
	}

	err = request.SendRequest(
		conf.Chats[definedChat].ChatAPIURL,
		conf.Chats[definedChat].ChatAPIToken,
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
