package components

import (
	"fmt"
	"strings"

	"github.com/aavshr/panda/internal/ui/styles"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Message struct {
	Content   string
	CreatedAt string
	IsUser    bool
}

type ChatModel struct {
	viewport       viewport.Model
	messages       []Message
	width          int
	height         int
	userStyle      lipgloss.Style
	assistantStyle lipgloss.Style
	timestampStyle lipgloss.Style
}

func NewChatModel(width, height int) ChatModel {
	vp := viewport.New(width, height)
	vp.KeyMap.PageDown.SetEnabled(true)
	vp.KeyMap.PageUp.SetEnabled(true)

	userStyle := styles.UserMessageStyle()
	assistantStyle := styles.AIMessageStyle()
	timestampStyle := styles.MetadataStyle()

	return ChatModel{
		viewport:       vp,
		messages:       []Message{},
		width:          width,
		height:         height,
		userStyle:      userStyle,
		assistantStyle: assistantStyle,
		timestampStyle: timestampStyle,
	}
}

func (m *ChatModel) SetMessage(index int, msg Message) {
	m.messages[index] = msg
	m.updateViewportContent()
}

func (m *ChatModel) AddMessage(msg Message) {
	m.messages = append(m.messages, msg)
	m.updateViewportContent()
	m.ScrollToBottom()
}

func (m *ChatModel) ResetMessages() {
	m.messages = []Message{}
	m.updateViewportContent()
}

func (m *ChatModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.viewport.Width = width
	m.viewport.Height = height
	m.updateViewportContent()
}

func (m *ChatModel) ScrollToBottom() {
	m.viewport.GotoBottom()
}

func (m *ChatModel) formatMessage(msg Message) string {
	if msg.Content == "" {
		return ""
	}
	var style lipgloss.Style
	var sender = "You"
	if msg.IsUser {
		style = m.userStyle
	} else {
		style = m.assistantStyle
		sender = "AI"
	}

	header := style.Render(sender) + m.timestampStyle.Render(msg.CreatedAt)
	contentWidth := m.width - 4
	wrappedContent := wrapText(msg.Content, contentWidth)
	indentedContent := strings.ReplaceAll(wrappedContent, "\n", "\n  ")
	return fmt.Sprintf("%s\n  %s\n", header, indentedContent)
}

// updateViewportContent updates the content in the viewport
func (m *ChatModel) updateViewportContent() {
	var sb strings.Builder

	for _, msg := range m.messages {
		sb.WriteString(m.formatMessage(msg))
		sb.WriteString("\n")
	}

	m.viewport.SetContent(sb.String())
	m.ScrollToBottom()
}

func wrapText(text string, width int) string {
	var result strings.Builder
	lines := strings.Split(text, "\n")

	for i, line := range lines {
		if i > 0 {
			result.WriteString("\n")
		}

		if len(line) <= width {
			result.WriteString(line)
			continue
		}

		words := strings.Fields(line)
		lineLength := 0

		for j, word := range words {
			if j > 0 {
				if lineLength+len(word)+1 > width {
					result.WriteString("\n")
					lineLength = 0
				} else {
					result.WriteString(" ")
					lineLength++
				}
			}

			result.WriteString(word)
			lineLength += len(word)
		}
	}

	return result.String()
}

func (m *ChatModel) View() string {
	return m.viewport.View()
}

func (m *ChatModel) Update(msg tea.Msg) (ChatModel, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEscape:
			return *m, EscapeCmd
		}

	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return *m, tea.Batch(cmds...)
}
