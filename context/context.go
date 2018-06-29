package context

// ContextActionACK - context for Mattermost action
type ContextActionACK struct {
	EventID  string
	Action   string
	Severity string
	Message  string
	Channel  string
	Username string
	IconURL  string
}
