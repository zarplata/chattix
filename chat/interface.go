package chat

// Message - interface which must be implemets by all
// chats requests representation
type Message interface {
	SetChannel(name string)
	SetUsername(name string)
	SetIcon(icon string)
	CreateAttachment(text string, color string) MessageAttachment
	GetAttachment(attachmentID int) (MessageAttachment, error)
	SendRequest(url string, token string) error
}

type MessageAttachment interface {
	SetText(text string)
	SetTitle(title string)
	SetColor(color string)
	AddAction(
		name string,
		text string,
		actionType string,
		value interface{},
	) AttachmentAction
	AddField(
		short bool,
		title string,
		value interface{},
	)
}

type AttachmentAction interface {
	SetText(text string)
	SetName(name string)
}
