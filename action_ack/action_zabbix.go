package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	karma "github.com/reconquest/karma-go"
)

type zabbixRequestPayload struct {
	JSONRPC string                 `json:"jsonrpc"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
	Auth    string                 `json:"auth,omitempty"`
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

type apiVersionZabbixResponse struct {
	JSONRPC string               `json:"jsonrpc"`
	Error   *zabbixResponseError `json:"error"`
	Result  string               `json:"result"`
	ID      int                  `json:"id"`
}

func getZabbixVersion(zabbixURL string) (string, error) {
	destiny := karma.Describe("method", "getZabbixVersion")

	payload := zabbixRequestPayload{
		JSONRPC: "2.0",
		Method:  "apiinfo.version",
		Params:  map[string]interface{}{},
		ID:      1,
	}

	body := new(bytes.Buffer)

	err := json.NewEncoder(body).Encode(payload)
	if err != nil {
		return "", destiny.Describe(
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
		return "", destiny.Describe(
			"zabbix URL", zabbixURL,
		).Describe(
			"error", err,
		).Reason(
			"can't send request to Zabbix",
		)
	}

	defer response.Body.Close()

	answer := apiVersionZabbixResponse{}

	err = json.NewDecoder(response.Body).Decode(&answer)
	if err != nil {
		return "", destiny.Reason(err)
	}

	if answer.Error != nil {
		return "", destiny.Describe(
			"error code", answer.Error.Code,
		).Describe(
			"error data", answer.Error.Data,
		).Reason(answer.Error.Message)
	}

	return answer.Result, nil
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

	params := map[string]interface{}{
		"eventids": eventID,
		"message":  acknowledgeMessage,
	}

	zabbixVersion, err := getZabbixVersion(zabbixURL)
	if err != nil {
		return destiny.Describe(
			"error", err,
		).Reason(
			"can't get Zabbix version",
		)
	}

	if len(strings.Split(zabbixVersion, ".")) < 1 {
		return destiny.Reason("can't parse zabbix version")
	}

	majorZabbixVersion := strings.Split(zabbixVersion, ".")[0]

	switch majorZabbixVersion {
	case "3":
		//https://www.zabbix.com/documentation/3.4/manual/api/reference/event/acknowledge
		params["action"] = 1

		//default:
		//https://www.zabbix.com/documentation/1.8/api/event/acknowledge
		//https://www.zabbix.com/documentation/2.0/manual/appendix/api/event/acknowledge
	default:
		//https://www.zabbix.com/documentation/4.0/manual/api/reference/event/acknowledge
		//https://www.zabbix.com/documentation/5.0/manual/api/reference/event/acknowledge
		params["action"] = 6
	}

	payload := zabbixRequestPayload{
		JSONRPC: "2.0",
		Method:  "event.acknowledge",
		Params:  params,
		Auth:    zabbixAPIToken,
		ID:      1,
	}

	body := new(bytes.Buffer)

	err = json.NewEncoder(body).Encode(payload)
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
