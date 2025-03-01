package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type SettingsSubmitMsg struct {
	APIKey string
}

type SettingsModel struct {
	inner textinput.Model
}

func SettingsSubmitCmd(msg SettingsSubmitMsg) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}

func NewSettingsModel() SettingsModel {
	inner := textinput.New()
	inner.Placeholder = "Enter your API key..."
	return SettingsModel{inner: inner}
}

func (m *SettingsModel) Focus() tea.Cmd {
	return m.inner.Focus()
}

func (m *SettingsModel) Blur() {
	m.inner.Blur()
}

func (m *SettingsModel) View() string {
	return m.inner.View()
}

func (m *SettingsModel) Update(msg interface{}) (SettingsModel, tea.Cmd) {
	var cmd tea.Cmd
	m.inner, cmd = m.inner.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			value := strings.TrimSpace(m.inner.Value())
			if value != "" {
				return *m, SettingsSubmitCmd(SettingsSubmitMsg{APIKey: m.inner.Value()})
			}
		case tea.KeyEscape, tea.KeyCtrlC, tea.KeyCtrlD:
			return *m, tea.Quit
		}
	}
	return *m, cmd
}
