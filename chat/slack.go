package chat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
)

// SlackMessage - represents Slack message
type SlackMessage struct {
	Text        string             `json:"text"`
	Username    string             `json:"username"`
	IconURL     string             `json:"icon_url"`
	ChannelName string             `json:"channel"`
	AsUser      bool               `json:"as_user"`
	Attachments []*SlackAttachment `json:"attachments"`
}

// SlackAttachment - represents an attachment in Slack message`
type SlackAttachment struct {
	ID         int64                   `json:"id"`
	Fallback   string                  `json:"fallback"`
	CallbackID string                  `json:"callback_id"`
	Color      string                  `json:"color"`
	Pretext    string                  `json:"pretext"`
	AuthorName string                  `json:"author_name"`
	AuthorLink string                  `json:"author_link"`
	AuthorIcon string                  `json:"author_icon"`
	Title      string                  `json:"title"`
	TitleLink  string                  `json:"title_link"`
	Text       string                  `json:"text"`
	Fields     []*SlackAttachmentField `json:"fields"`
	ImageURL   string                  `json:"image_url"`
	ThumbURL   string                  `json:"thumb_url"`
	Footer     string                  `json:"footer"`
	FooterIcon string                  `json:"footer_icon"`
	Timestamp  interface{}             `json:"ts"` // This is either a string or an int64
	Actions    []*SlackAction          `json:"actions,omitempty"`
}

// AddAction - add action to Slack attachment
func (attachment *SlackAttachment) AddAction(
	name string,
	text string,
	actionType string,
	value interface{},
) AttachmentAction {

	action := &SlackAction{
		Name:  name,
		Text:  text,
		Type:  actionType,
		Value: value,
	}

	attachment.Actions = append(
		attachment.Actions,
		action,
	)

	return action
}

// SetColor - set color to attachment
func (attachment *SlackAttachment) SetColor(
	color string,
) {
	attachment.Color = color
}

// SetText - set text to attachment
func (attachment *SlackAttachment) SetText(
	text string,
) {
	attachment.Text = text
}

// SetTitle - set title for attachment
func (attachment *SlackAttachment) SetTitle(
	title string,
) {
	attachment.Title = title
}

// AddField - add field to attachment
func (attachment *SlackAttachment) AddField(
	short bool,
	title string,
	value interface{},
) {
	field := &SlackAttachmentField{
		Title: title,
		Value: value,
		Short: short,
	}

	attachment.Fields = append(
		attachment.Fields,
		field,
	)
}

// SlackAttachmentField - represents a field for Slack attachment
type SlackAttachmentField struct {
	Title string      `json:"title"`
	Value interface{} `json:"value"`
	Short bool        `json:"short"`
}

// SlackAction - representation of Slack action
type SlackAction struct {
	Name  string      `json:"name"`
	Text  string      `json:"text"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

// NewSlackAction - creates a new action for Slack attachment
func NewSlackAction() AttachmentAction {
	return &SlackAction{}
}

// SetText - set text to action
func (action *SlackAction) SetText(
	text string,
) {
	action.Text = text
}

// SetName - set name to Slack attachment action
func (action *SlackAction) SetName(
	name string,
) {
	action.Name = name
}

// NewSlackMessage - creates new Slack message
func NewSlackMessage() Message {
	return &SlackMessage{}
}

// SetChannel - set channel where message will be posted
func (request *SlackMessage) SetChannel(
	name string,
) {
	request.ChannelName = name
}

// SetIcon - set icon URL to message
func (request *SlackMessage) SetIcon(
	icon string,
) {
	request.IconURL = icon
}

// SetUsername - set username which will post a messages.
// For Slack it will be random name because Slack glue
// messages which posted from one username and doesn't
// display an icon.
func (request *SlackMessage) SetUsername(
	name string,
) {
	request.Username = fmt.Sprintf(
		"%s-%s",
		name,
		getRandString(7),
	)
}

// SendRequest - send message to Slack
func (request *SlackMessage) SendRequest(
	url string, token string,
) error {
	payload, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(payload),
	)
	if err != nil {
		return err
	}

	req.Header.Add(
		"Content-Type",
		"application/json",
	)

	if len(token) != 0 {
		req.Header.Add(
			"Authorization",
			"Bearer "+token,
		)
	}

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"chat on %s returned %d status code",
			url,
			response.StatusCode,
		)
	}

	return nil
}

// CreateAttachment - create new message attachment and append it
func (request *SlackMessage) CreateAttachment(
	text string, color string,
) MessageAttachment {
	attachment := &SlackAttachment{
		Color:      color,
		Text:       text,
		CallbackID: "ack",
	}

	request.Attachments = append(
		request.Attachments,
		attachment,
	)

	return attachment
}

// GetAttachment - get attachment from message
// by their index
func (request *SlackMessage) GetAttachment(
	attachmentID int,
) (MessageAttachment, error) {
	if len(request.Attachments) < attachmentID+1 {
		return nil, fmt.Errorf(
			"attachement %d did not found",
			attachmentID,
		)
	}

	return request.Attachments[attachmentID], nil
}

func getRandString(n int) string {
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
