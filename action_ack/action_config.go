package main

import (
	"fmt"
	"os"
)

var defaultConfiguration = `

listen_address = "0.0.0.0:5666"

[zabbix]
zabbix_api_url = "http://localhost/api_jsonrpc.php"
zabbix_api_token = "token"

[messenger]
    [messenger.mattermost]
    messenger_api_token = "secret_user_token"
    messenger_api_url = "https://localhost/api/v4"
    attachments_color = "#000000"
    author_message = "Acknowledged by {{USERNAME}}"
    author_image_url = "http://localhost/image"

    [messenger.slack]
    messenger_api_token = "secret_user_token"
    messenger_api_url = "https://slack.com/api"
    attachments_color = "#000000"
    author_message = "Acknowledged by {{USERNAME}}"
    author_image_url = "http://localhost/image"
`

type config struct {
	ListenAddress string                     `toml:"listen_address"`
	Zabbix        zabbixConfig               `toml:"zabbix"`
	Messenger     map[string]messengerConfig `toml:"messenger"`
}

type zabbixConfig struct {
	ZabbixAPIURL   string `toml:"zabbix_api_url"`
	ZabbixAPIToken string `toml:"zabbix_api_token"`
}

type messengerConfig struct {
	MessengerAPIToken string `toml:"messenger_api_token"`
	MessengerAPIURL   string `toml:"messenger_api_url"`
	AttachmentsColor  string `toml:"attachments_color"`
	AuthorMessage     string `toml:"author_message"`
	AuthorImageURL    string `toml:"author_image_url"`
}

func parseEnvironmentVariables(
	config *config,
) {
	if value := os.Getenv("CHATTIX_PORT"); value != "" {
		config.ListenAddress = fmt.Sprintf(
			"0.0.0.0:%s",
			value,
		)
	}

	if value := os.Getenv("CHATTIX_ZABBIX_URL"); value != "" {
		config.Zabbix.ZabbixAPIURL = value
	}

	if value := os.Getenv("CHATTIX_ZABBIX_TOKEN"); value != "" {
		config.Zabbix.ZabbixAPIToken = value
	}

	if value := os.Getenv("CHATTIX_ZABBIX_USER"); value != "" {
		//
	}

	if value := os.Getenv("CHATTIX_ZABBIX_PASSWORD"); value != "" {
		//
	}

	if value := os.Getenv("CHATTIX_MESSENGER_API_TOKEN"); value != "" {
		messengerConfig := config.Messenger[definedMessenger]
		messengerConfig.MessengerAPIToken = value

		config.Messenger[definedMessenger] = messengerConfig

	}

	if value := os.Getenv("CHATTIX_MESSENGER_API_URL"); value != "" {
		messengerConfig := config.Messenger[definedMessenger]
		messengerConfig.MessengerAPIURL = value

		config.Messenger[definedMessenger] = messengerConfig

	}

	if value := os.Getenv("CHATTIX_MESSENGER_ATTACHMENT_COLOR"); value != "" {
		messengerConfig := config.Messenger[definedMessenger]
		messengerConfig.AttachmentsColor = value

		config.Messenger[definedMessenger] = messengerConfig

	}

	if value := os.Getenv("CHATTIX_MESSENGER_AUTHOR_MESSAGE"); value != "" {
		messengerConfig := config.Messenger[definedMessenger]
		messengerConfig.AuthorMessage = value

		config.Messenger[definedMessenger] = messengerConfig

	}

	if value := os.Getenv("CHATTIX_MESSENGER_AUTHOR_IMAGE_URL"); value != "" {
		messengerConfig := config.Messenger[definedMessenger]
		messengerConfig.AuthorImageURL = value

		config.Messenger[definedMessenger] = messengerConfig

	}
}
