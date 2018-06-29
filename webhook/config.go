package main

import (
	"math/rand"
	"time"
)

type config struct {
	Chats         map[string]chatConfig     `toml:"chat"`
	EventIDRegexp string                    `toml:"event_id_regexp"`
	Severities    map[string]severityConfig `toml:"severities"`
	Actions       map[string]actionConfig   `toml:"actions"`
}

type chatConfig struct {
	ChatAPIURL   string `toml:"chat_api_url"`
	ChatAPIToken string `toml:"chat_api_token"`
	ChatUsername string `toml:"chat_username"`
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
