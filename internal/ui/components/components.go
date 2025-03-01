package components

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Component string

const (
	ComponentHistory       Component = "history"
	ComponentMessages      Component = "messages"
	ComponentInnerMessages Component = "innerMessages"
	ComponentChatInput     Component = "chatInput"
	ComponentSettings      Component = "settings"
	ComponentNone          Component = "none" // utility component
)

type EscapeMsg struct{}

func EscapeCmd() tea.Msg {
	return EscapeMsg{}
}
