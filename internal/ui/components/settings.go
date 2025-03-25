package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type SettingsMode int

const (
	SettingsModeAPIKey SettingsMode = iota
	SettingsModeLLMModel
)

type SettingsSubmitMsg struct {
	APIKey   string
	LLMModel string
}

type SettingsModel struct {
	inner textinput.Model
	mode  SettingsMode
	msg   SettingsSubmitMsg
}

func SettingsSubmitCmd(msg SettingsSubmitMsg) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}

func NewSettingsModel() SettingsModel {
	inner := textinput.New()
	inner.Placeholder = "Enter your API key..."
	return SettingsModel{
		inner: inner,
		mode:  SettingsModeAPIKey,
	}
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
				if m.mode == SettingsModeAPIKey {
					m.mode = SettingsModeLLMModel
					m.inner.Placeholder = "Enter your LLM model (default o3-mini)..."
					m.msg.APIKey = value
					m.inner.SetValue("")
					m.View()
				} else {
					m.msg.LLMModel = value
					return *m, SettingsSubmitCmd(m.msg)
				}
			}
		case tea.KeyEscape, tea.KeyCtrlC, tea.KeyCtrlD:
			return *m, tea.Quit
		}
	}
	return *m, cmd
}
