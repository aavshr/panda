package components

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type ChatInputEnterMsg struct {
	Value string
}

type ChatInputModel struct {
	inner textarea.Model
}

func NewChatInputModel(width, height int) ChatInputModel {
	inner := textarea.New()
	inner.SetWidth(width)
	inner.SetHeight(height)

	inner.Placeholder = "Send message..."
	inner.ShowLineNumbers = false
	inner.Prompt = ""

	return ChatInputModel{
		inner: inner,
	}
}

func (c *ChatInputModel) Init() tea.Cmd {
	return textarea.Blink
}

func (c *ChatInputModel) View() string {
	return c.inner.View()
}

func (c *ChatInputModel) Focus() tea.Cmd {
	return c.inner.Focus()
}

func (c *ChatInputModel) Value() string {
	return c.inner.Value()
}

func (c *ChatInputModel) EnterCmd(value string) tea.Cmd {
	return func() tea.Msg {
		return ChatInputEnterMsg{Value: value}
	}
}

func (c *ChatInputModel) Update(msg tea.Msg) (ChatInputModel, tea.Cmd) {
	var cmd tea.Cmd
	c.inner, cmd = c.inner.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyCtrlD:
			return *c, tea.Quit
		case tea.KeyEnter:
			value := c.inner.Value()
			c.inner.Reset()
			return *c, c.EnterCmd(value)
		}
	case error:
		// TODO: how to handle errors
		return *c, cmd
	}
	return *c, cmd
}
