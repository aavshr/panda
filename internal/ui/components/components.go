package components

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Component string

const (
	ComponentHistory   Component = "history"
	ComponentMessages  Component = "messages"
	ComponentChatInput Component = "chatInput"
	ComponentNone      Component = "none" // utility component
)

type EscapeMsg struct{}

func EscapeCmd() tea.Msg {
	return EscapeMsg{}
}
