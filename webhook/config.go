package main

import (
	"math/rand"
	"time"
)

type config struct {
	Messengers    map[string]messengerConfig `toml:"messenger"`
	EventIDRegexp string                     `toml:"event_id_regexp"`
	Severities    map[string]severityConfig  `toml:"severities"`
	Actions       map[string]actionConfig    `toml:"actions"`
}

type messengerConfig struct {
	MessengerAPIURL   string `toml:"messenger_api_url"`
	MessengerAPIToken string `toml:"messenger_api_token"`
	MessengerUsername string `toml:"messenger_username"`
}

type severityConfig struct {
	ImageURLs []string `toml:"image_urls"`
	Color     string   `toml:"color"`
}

type actionConfig struct {
	ActionName string `toml:"action_name"`
	ActionURL  string `toml:"action_url"`
}

func (c *config) getIconURL(
	zabbixSeverity string,
) string {
	if severity, exists := c.Severities[zabbixSeverity]; exists {
		if len(severity.ImageURLs) == 1 {
			return severity.ImageURLs[0]
		}

		if len(severity.ImageURLs) > 1 {
			rand.Seed(time.Now().UTC().UnixNano())
			return severity.ImageURLs[rand.Intn(len(severity.ImageURLs))]
		}
	}

	return ""
}

func (c *config) getColor(
	zabbixSeverity string,
) string {
	if severity, exists := c.Severities[zabbixSeverity]; exists {
		return severity.Color
	}

	return ""
}
