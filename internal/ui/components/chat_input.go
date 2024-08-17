package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type ChatInputReturnMsg struct {
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

	inner.KeyMap.InsertNewline.SetEnabled(true)
	inner.Cursor.SetMode(cursor.CursorStatic)

	return ChatInputModel{
		inner: inner,
	}
}

func (c *ChatInputModel) View() string {
	return c.inner.View()
}

func (c *ChatInputModel) Focus() tea.Cmd {
	return c.inner.Focus()
}

func (c *ChatInputModel) Blur() {
	c.inner.Blur()
}

func (c *ChatInputModel) Value() string {
	return c.inner.Value()
}

func (c *ChatInputModel) EnterCmd(value string) tea.Cmd {
	return func() tea.Msg {
		return ChatInputReturnMsg{Value: value}
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
		case tea.KeyTab:
			value := strings.TrimSpace(c.inner.Value())
			if value != "" {
				c.inner.Reset()
				return *c, c.EnterCmd(value)
			}
		case tea.KeyEscape:
			return *c, EscapeCmd
		}
	}
	return *c, cmd
}
