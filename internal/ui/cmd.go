package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type SelectComponentMsg struct{}
type FocusComponentMsg struct{}
type ForwardChatCompletionStreamMsg struct{}

func (m *Model) cmdSelectComponent() tea.Msg {
	return SelectComponentMsg{}
}

func (m *Model) cmdFocusedComponent() tea.Msg {
	return FocusComponentMsg{}
}

func (m *Model) cmdForwardChatCompletionStream() tea.Msg {
	return ForwardChatCompletionStreamMsg{}
}

func (m *Model) cmdError(err error) func() tea.Msg {
	return func() tea.Msg {
		return err
	}
}
