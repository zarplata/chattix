package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	karma "github.com/reconquest/karma-go"
	chat "github.com/zarplata/chattix/chat"
)

type slackActionRequest struct {
	Type       string              `json:"type"`
	Actions    []*chat.SlackAction `json:"actions"`
	CallbackID string              `json:"callback_id"`

	Team struct {
		ID     string `json:"id"`
		Domain string `json:"domain"`
	} `json:"team"`

	Channel struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"channel"`

	User struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"user"`

	OriginalMessage *chat.SlackMessage `json:"original_message"`
}

type slackUserResponse struct {
	Ok   bool `json:"ok"`
	User struct {
		RealName string `json:"real_name"`
	} `json:"user"`
	Error string `json:"error"`
}

func fetchUserFromSlack(
	chatAPIURL string,
	chatAPIToken string,
	userID string,
) (string, error) {

	destiny := karma.Describe(
		"method", "fetchUserFromSlack",
	).Describe(
		"url", chatAPIURL,
	).Describe(
		"authToken", chatAPIToken,
	).Describe(
		"userID", userID,
	)

	requestURL := fmt.Sprintf(
		"%s/users.info?user=%s",
		chatAPIURL,
		userID,
	)

	request, err := http.NewRequest(
		"GET",
		requestURL,
		nil,
	)
	if err != nil {
		return "", destiny.Describe(
			"request url", requestURL,
		).Describe(
			"error", err,
		).Reason(
			"can't create HTTP request",
		)
	}

	request.Header.Set(
		"Authorization",
		fmt.Sprintf("Bearer %s", chatAPIToken),
	)

	request.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", destiny.Describe(
			"error", err,
		).Reason(
			"can't execute HTTP request",
		)
	}

	defer response.Body.Close()

	var body slackUserResponse

	err = json.NewDecoder(response.Body).Decode(&body)
	if err != nil {
		return "", destiny.Describe(
			"error", err,
		).Reason(
			"can't decode Slack user response",
		)
	}

	if response.StatusCode != http.StatusOK {
		return "", destiny.Describe(
			"status code", response.StatusCode,
		).Describe(
			"error", body.Error,
		).Reason(
			"unexpected status code from Slack API",
		)
	}

	if !body.Ok {
		return "", destiny.Describe(
			"error", body.Error,
		).Reason(
			"Non ok anwer from Slack API",
		)
	}

	return body.User.RealName, nil

}
