package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type MsgSelectComponent struct{}
type MsgFocusComponent struct{}

func (m *Model) cmdSelectComponent() tea.Msg {
	return MsgSelectComponent{}
}

func (m *Model) cmdFocusedComponent() tea.Msg {
	return MsgFocusComponent{}
}
