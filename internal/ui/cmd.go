package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type SelectComponentMsg struct{}
type FocusComponentMsg struct{}
type CreateChatCompletionStreamMsg struct {
	history []string
	prompt  string
}

func (m *Model) cmdSelectComponent() tea.Msg {
	return SelectComponentMsg{}
}

func (m *Model) cmdFocusedComponent() tea.Msg {
	return FocusComponentMsg{}
}

func (m *Model) cmdCreateChatCompletionStream(prompt string, history []string) tea.Cmd {
	return func() tea.Msg {
		return CreateChatCompletionStreamMsg{
			prompt:  prompt,
			history: history,
		}
	}
}

func (m *Model) cmdError(err error) func() tea.Msg {
	return func() tea.Msg {
		return err
	}
}
