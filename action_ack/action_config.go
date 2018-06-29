package main

type config struct {
	ListenAddress string                `toml:"listen_address"`
	Zabbix        zabbixConfig          `toml:"zabbix"`
	Chat          map[string]chatConfig `toml:"chat"`
}

type zabbixConfig struct {
	ZabbixAPIURL   string `toml:"zabbix_api_url"`
	ZabbixAPIToken string `toml:"zabbix_api_token"`
}

type chatConfig struct {
	ChatAPIToken     string `toml:"chat_api_token"`
	ChatAPIURL       string `toml:"chat_api_url"`
	AttachmentsColor string `toml:"attachments_color"`
	AuthorMessage    string `toml:"author_message"`
	AuthorImageURL   string `toml:"author_image_url"`
}
