package chat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// MattermostMessage - represents Mattermost message
type MattermostMessage struct {
	Text        string                  `json:"text"`
	Username    string                  `json:"username"`
	IconURL     string                  `json:"icon_url"`
	ChannelName string                  `json:"channel"`
	Props       map[string]interface{}  `json:"props"`
	Attachments []*MattermostAttachment `json:"attachments"`
}

// MattermostAttachment - represents Mattermost message attachment
type MattermostAttachment struct {
	ID         int64                        `json:"id"`
	Fallback   string                       `json:"fallback"`
	Color      string                       `json:"color"`
	Pretext    string                       `json:"pretext"`
	AuthorName string                       `json:"author_name"`
	AuthorLink string                       `json:"author_link"`
	AuthorIcon string                       `json:"author_icon"`
	Title      string                       `json:"title"`
	TitleLink  string                       `json:"title_link"`
	Text       string                       `json:"text"`
	Fields     []*MattermostAttachmentField `json:"fields"`
	ImageURL   string                       `json:"image_url"`
	ThumbURL   string                       `json:"thumb_url"`
	Footer     string                       `json:"footer"`
	FooterIcon string                       `json:"footer_icon"`
	Timestamp  interface{}                  `json:"ts"` // This is either a string or an int64
	Actions    []*MattermostAction          `json:"actions,omitempty"`
}

// AddAction - add an action to attachment
func (attachment *MattermostAttachment) AddAction(
	name string,
	text string,
	actionType string,
	context interface{},
) AttachmentAction {
	integration := &MattermostActionIntegration{
		URL:     text,
		Context: context.(map[string]interface{}),
	}

	action := &MattermostAction{
		Name:        name,
		Integration: integration,
	}

	attachment.Actions = append(
		attachment.Actions,
		action,
	)

	return action
}

// SetColor - set attachment color
func (attachment *MattermostAttachment) SetColor(
	color string,
) {
	attachment.Color = color
}

// SetText - set attachment text
func (attachment *MattermostAttachment) SetText(
	text string,
) {
	attachment.Text = text
}

// SetTitle - set title for attachment
func (attachment *MattermostAttachment) SetTitle(
	title string,
) {
	attachment.Title = title
}

// AddField - add field to attachment
func (attachment *MattermostAttachment) AddField(
	short bool,
	title string,
	value interface{},
) {
	field := &MattermostAttachmentField{
		Title:   title,
		Message: value,
		Short:   short,
	}

	attachment.Fields = append(
		attachment.Fields,
		field,
	)
}

// MattermostAttachmentField - represents field for message attachment
type MattermostAttachmentField struct {
	Title   string      `json:"title"`
	Message interface{} `json:"value"`
	Short   bool        `json:"short"`
}

// MattermostAction - represents of action for message attachment
type MattermostAction struct {
	ID          string                       `json:"id"`
	Name        string                       `json:"name"`
	Integration *MattermostActionIntegration `json:"integration,omitempty"`
}

// MattermostActionIntegration - represents an integration information
// for attachment action
type MattermostActionIntegration struct {
	URL     string                 `json:"url,omitempty"`
	Context map[string]interface{} `json:"context,omitempty"`
}

// SetText - set action text
func (action *MattermostAction) SetText(
	text string,
) {
	action.Integration.URL = text
}

// SetName - set name of action
func (action *MattermostAction) SetName(
	name string,
) {
	action.Name = name
}

// NewMattermostMessage - creates a new Mattermost message
func NewMattermostMessage() Message {
	return &MattermostMessage{}
}

// SetChannel - set channel where message will
// be delivered
func (request *MattermostMessage) SetChannel(
	name string,
) {
	request.ChannelName = name
}

// SetIcon - set icon URL to message
func (request *MattermostMessage) SetIcon(
	icon string,
) {
	request.IconURL = icon
}

// SetUsername - set username for message
func (request *MattermostMessage) SetUsername(
	name string,
) {
	request.Username = name
}

// CreateAttachment - creates attachment to Mattermost message
// with passed text and color
func (request *MattermostMessage) CreateAttachment(
	text string, color string,
) MessageAttachment {
	attachment := &MattermostAttachment{
		Color: color,
		Text:  text,
	}

	request.Attachments = append(
		request.Attachments,
		attachment,
	)

	return attachment
}

// GetAttachment - return attachment by index
// from message
func (request *MattermostMessage) GetAttachment(
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

// SendRequest - sending request to Mattermost
func (request *MattermostMessage) SendRequest(
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

func (attachment *MattermostAttachment) createAction(
	name string,
	actionURL string,
	context map[string]interface{},
) {
	integration := &MattermostActionIntegration{
		URL:     actionURL,
		Context: context,
	}

	action := &MattermostAction{
		Name:        name,
		Integration: integration,
	}

	attachment.Actions = append(attachment.Actions, action)
}
