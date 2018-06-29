package main

import (
	"bytes"
	"encoding/json"
	"net/http"

	karma "github.com/reconquest/karma-go"
)

type zabbixRequestPayload struct {
	JSONRPC string                 `json:"jsonrpc"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
	Auth    string                 `json:"auth"`
	ID      int                    `json:"id"`
}

type zabbixResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

type zabbixResponse struct {
	JSONRPC string                 `json:"jsonrpc"`
	Error   *zabbixResponseError   `json:"error"`
	Result  map[string]interface{} `json:"result"`
	ID      int                    `json:"id"`
}

func acknowledgeZabbixEvent(
	zabbixURL string,
	zabbixAPIToken string,
	eventID string,
	acknowledgeMessage string,
) error {
	destiny := karma.Describe(
		"method", "acknowledgeZabbixEvent",
	).Describe(
		"eventID", eventID,
	).Describe(
		"API Token", zabbixAPIToken,
	)

	payload := zabbixRequestPayload{
		JSONRPC: "2.0",
		Method:  "event.acknowledge",
		Params: map[string]interface{}{
			"eventids": eventID,
			"message":  acknowledgeMessage,
		},
		Auth: zabbixAPIToken,
		ID:   1,
	}

	body := new(bytes.Buffer)

	err := json.NewEncoder(body).Encode(payload)
	if err != nil {
		return destiny.Describe(
			"error", err,
		).Reason(
			"can't encode payload for Zabbix",
		)
	}

	response, err := http.Post(
		zabbixURL,
		"application/json",
		body,
	)

	if err != nil {
		return destiny.Describe(
			"zabbix URL", zabbixURL,
		).Describe(
			"error", err,
		).Reason(
			"can't send request to Zabbix",
		)
	}

	defer response.Body.Close()

	answer := zabbixResponse{}

	err = json.NewDecoder(response.Body).Decode(&answer)
	if err != nil {
		return destiny.Reason(err)
	}

	if answer.Error != nil {
		return destiny.Describe(
			"error code", answer.Error.Code,
		).Describe(
			"error data", answer.Error.Data,
		).Reason(answer.Error.Message)
	}

	return nil
}
