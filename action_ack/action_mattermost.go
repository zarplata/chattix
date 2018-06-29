package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	karma "github.com/reconquest/karma-go"
	"github.com/zarplata/chattix/context"
)

type actionRequest struct {
	UserID  string                   `json:"user_id"`
	Context context.ContextActionACK `json:"context"`
}

func fetchUserFromMattermost(
	chatURL string,
	authToken string,
	userID string,
) (string, error) {
	destiny := karma.Describe(
		"method", "fetchUserFromMattermost",
	).Describe(
		"url", chatURL,
	).Describe(
		"authToken", authToken,
	).Describe(
		"userID", userID,
	)

	payload := []string{
		userID,
	}

	body := new(bytes.Buffer)

	err := json.NewEncoder(body).Encode(payload)
	if err != nil {
		return "", destiny.Describe(
			"error", err,
		).Reason("can't marshal request payload")
	}

	requestURL := fmt.Sprintf(
		"%s/users/ids",
		chatURL,
	)

	destiny.Describe(
		"request URL", requestURL,
	)

	request, err := http.NewRequest(
		"POST",
		requestURL,
		body,
	)

	if err != nil {
		return "", destiny.Describe(
			"error", err,
		).Reason("can't create HTTP request")
	}

	request.Header.Set(
		"Authorization",
		fmt.Sprintf("Bearer %s", authToken),
	)

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", destiny.Describe(
			"error", err,
		).Reason("can't execute HTTP request")
	}

	defer func() {
		err = response.Body.Close()
		if err != nil {
			logger.Error(
				destiny.Describe(
					"error", err,
				).Reason(
					"can't close response body",
				),
			)
		}
	}()

	if response.StatusCode != http.StatusOK {
		answer := map[string]interface{}{}

		err = json.NewDecoder(response.Body).Decode(&answer)
		if err != nil {
			return "", destiny.Describe(
				"error", err,
			).Reason(
				"can't decode JSON response from Mattermost",
			)
		}

		return "", destiny.Describe(
			"status code", answer["status_code"].(int),
		).Describe(
			"request id", answer["request_id"].(string),
		).Reason(answer["message"].(string))
	}

	answer := []interface{}{}
	err = json.NewDecoder(response.Body).Decode(&answer)
	if err != nil {
		return "", destiny.Describe(
			"error", err,
		).Reason(
			"can't decode JSON response from Mattermost",
		)
	}

	if len(answer) != 1 {
		return "", destiny.Describe(
			"user count", len(answer),
		).Reason(
			"unexpected count of user requrned, expected 1",
		)
	}

	fullName := strings.TrimSpace(
		fmt.Sprintf(
			"%s %s",
			answer[0].(map[string]interface{})["first_name"].(string),
			answer[0].(map[string]interface{})["last_name"].(string),
		),
	)

	return fullName, nil
}
