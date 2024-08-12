package components

type Component string

const (
	ComponentHistory   Component = "history"
	ComponentMessages  Component = "messages"
	ComponentChatInput Component = "chatInput"
	ComponentNone      Component = "none" // utility component
)

type MsgEscape struct{}
