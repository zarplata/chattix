package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kovetskiy/lorg"
	karma "github.com/reconquest/karma-go"
	chat "github.com/zarplata/chattix/chat"
)

const (
	acknowledgedStatus  = "ACKNOWLEDGED"
	usernamePlaceholder = "{{USERNAME}}"
	chatSlack           = "slack"
	chatMattermost      = "mattermost"
)

type actionACKService struct {
	config   *config
	gin      *gin.Engine
	logger   *lorg.Log
	chatType string
}

func newActionACKService(
	config *config,
	logger *lorg.Log,
	chatType string,
) *actionACKService {

	service := &actionACKService{
		config:   config,
		gin:      gin.Default(),
		logger:   logger,
		chatType: chatType,
	}

	return service
}

func (service *actionACKService) setRoute() {
	if service.chatType == chatSlack {
		service.gin.POST("/", service.handleACKSlack)
		service.gin.GET("/", service.handleACKSlack)
	}

	if service.chatType == chatMattermost {
		service.gin.POST("/", service.handleACKMattermost)
		service.gin.GET("/", service.handleACKMattermost)
	}
}

func (service *actionACKService) run() {
	service.setRoute()
	service.gin.Run(service.config.ListenAddress)
}

func (service *actionACKService) handleACKSlack(
	context *gin.Context,
) {

	destiny := karma.Describe(
		"method", "handleACKSlack",
	)

	var payload slackActionRequest

	rawPayload := context.Request.FormValue("payload")

	err := json.Unmarshal([]byte(rawPayload), &payload)
	if err != nil {
		service.logger.Error(
			destiny.Describe(
				"error", err,
			).Reason(
				"can't unmarshal payload from Slack",
			),
		)
		context.JSON(sendInternalServerError(destiny))
		return
	}

	if len(payload.Actions) < 1 {
		service.logger.Error(
			destiny.Describe(
				"count of actions", len(payload.Actions),
			).Reason(
				"request from Slack should contains an action",
			),
		)

		context.JSON(
			sendInternalServerError(destiny),
		)
		return

	}

	// EventID which send by webhook binary as action value
	eventID := payload.Actions[0].Value.(string)

	username, err := fetchUserFromSlack(
		service.config.Chat[service.chatType].ChatAPIURL,
		service.config.Chat[service.chatType].ChatAPIToken,
		payload.User.ID,
	)

	if err != nil {
		service.logger.Error(
			destiny.Describe(
				"error", err,
			).Reason(
				"can't fetch user from Slack",
			),
		)
		context.JSON(sendInternalServerError(destiny))
		return
	}

	newColor := service.config.Chat[service.chatType].AttachmentsColor
	authorMessage := strings.Replace(
		service.config.Chat[service.chatType].AuthorMessage,
		usernamePlaceholder,
		username,
		-1,
	)

	message := payload.OriginalMessage
	if len(message.Attachments) < 1 {
		service.logger.Error(
			destiny.Describe(
				"count of attachment", len(message.Attachments),
			).Reason(
				"original message should contains an attachment",
			),
		)

		context.JSON(
			sendInternalServerError(destiny),
		)
		return
	}

	message.Attachments[0].Color = newColor
	message.Attachments[0].Title = acknowledgedStatus
	message.Attachments[0].Actions = []*chat.SlackAction{}

	zabbixAttachment := &chat.SlackAttachment{
		AuthorName: authorMessage,
		AuthorIcon: service.config.Chat[service.chatType].AuthorImageURL,
		Color:      newColor,
	}

	message.Attachments = append(
		message.Attachments,
		zabbixAttachment,
	)

	err = acknowledgeZabbixEvent(
		service.config.Zabbix.ZabbixAPIURL,
		service.config.Zabbix.ZabbixAPIToken,
		eventID,
		authorMessage,
	)

	if err != nil {

		service.logger.Error(
			destiny.Describe(
				"error", err,
			).Reason(
				"can't acknowledge Zabbix event",
			),
		)

		context.JSON(sendInternalServerError(destiny))
		return
	}

	context.JSON(http.StatusOK, message)

}

func (service *actionACKService) handleACKMattermost(
	context *gin.Context,
) {
	destiny := karma.Describe(
		"method", "handleACKMattermost",
	)

	var request actionRequest

	err := json.NewDecoder(context.Request.Body).Decode(&request)
	if err != nil {
		service.logger.Error(
			destiny.Describe(
				"error", err,
			).Reason(
				"can't unmarshal payload from Mattermost",
			),
		)
		context.JSON(sendInternalServerError(destiny))
		return
	}

	username, err := fetchUserFromMattermost(
		service.config.Chat[service.chatType].ChatAPIURL,
		service.config.Chat[service.chatType].ChatAPIToken,
		request.UserID,
	)
	if err != nil {
		service.logger.Error(
			destiny.Describe(
				"error", err,
			).Reason(
				"can't fetch user from Mattermost",
			),
		)
		context.JSON(sendInternalServerError(destiny))
		return
	}

	mattermostMessage := chat.MattermostMessage{
		ChannelName: request.Context.Channel,
		Username:    request.Context.Username,
		IconURL:     request.Context.IconURL,
	}

	attachment := &chat.MattermostAttachment{
		Color: service.config.Chat[service.chatType].AttachmentsColor,
		Text:  request.Context.Message,
		Title: acknowledgedStatus,
	}

	attachment.AddField(
		false,
		"Event ID",
		request.Context.EventID,
	)

	authorMessage := strings.Replace(
		service.config.Chat[service.chatType].AuthorMessage,
		usernamePlaceholder,
		username,
		-1,
	)

	attachmentZabbix := &chat.MattermostAttachment{
		Color:      service.config.Chat[service.chatType].AttachmentsColor,
		AuthorName: authorMessage,
		AuthorIcon: service.config.Chat[service.chatType].AuthorImageURL,
	}

	mattermostMessage.Attachments = append(
		mattermostMessage.Attachments,
		attachment,
		attachmentZabbix,
	)

	response := map[string]interface{}{
		"update": map[string]interface{}{
			"props": mattermostMessage,
		},
	}

	err = acknowledgeZabbixEvent(
		service.config.Zabbix.ZabbixAPIURL,
		service.config.Zabbix.ZabbixAPIToken,
		request.Context.EventID,
		authorMessage,
	)

	if err != nil {
		service.logger.Error(
			destiny.Describe(
				"error", err,
			).Reason(
				"can't acknowledge Zabbix event",
			),
		)

		context.JSON(sendInternalServerError(destiny))
		return
	}

	context.JSON(http.StatusOK, response)
}

func sendInternalServerError(
	responseData interface{},
) (int, interface{}) {
	return http.StatusInternalServerError, responseData
}
