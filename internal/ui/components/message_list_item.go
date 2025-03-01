package components

import (
	"fmt"
	"io"

	"github.com/aavshr/panda/internal/db"
	"github.com/aavshr/panda/internal/ui/styles"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

// MessageListItem implements the list.Item, list.DefaultItem interface
type MessageListItem struct {
	message *db.Message
}

func (t *MessageListItem) Title() string {
	return t.message.Content
}

func (t *MessageListItem) Description() string {
	if t.message.Role == "" || t.message.CreatedAt == "" {
		return ""
	}
	role := "You"
	if t.message.Role == "assistant" {
		role = "AI"
	}
	return fmt.Sprintf("%s at %s", role, t.message.CreatedAt)
}

func (t *MessageListItem) FilterValue() string {
	return t.Title()
}

// MessageListItemDelegate implements list.ItemDelegate interface
type MessageListItemDelegate struct {
	userStyle lipgloss.Style
	aiStyle   lipgloss.Style
	metaStyle lipgloss.Style
}

func (d *MessageListItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	messageItem, ok := item.(*MessageListItem)
	if !ok {
		return
	}

	contentWidth := m.Width() - 4 // account for padding

	var contentStyle lipgloss.Style
	if messageItem.message.Role == "assistant" {
		contentStyle = d.aiStyle
	} else {
		contentStyle = d.userStyle
	}

	content := wordwrap.String(messageItem.Title(), contentWidth)
	meta := d.metaStyle.Render(messageItem.Description())

	fmt.Fprintln(w, contentStyle.Render(content))
	fmt.Fprintln(w, meta)

	fmt.Fprintln(w)
}

func (d *MessageListItemDelegate) Height() int {
	return 2 + d.Spacing() // Minimum for content + metadata + spacing
}

func (d *MessageListItemDelegate) Spacing() int {
	return 0
}

func (d *MessageListItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func NewMessageListItem(message *db.Message) *MessageListItem {
	return &MessageListItem{
		message: message,
	}
}

func NewMessageListItems(messages []*db.Message) []list.Item {
	items := make([]list.Item, len(messages))
	for i, m := range messages {
		items[i] = &MessageListItem{
			message: m,
		}
	}
	return items
}

// NewMessageListItemDelegate maintains the existing API while enhancing functionality
func NewMessageListItemDelegate() list.ItemDelegate {
	return &MessageListItemDelegate{
		userStyle: styles.UserMessageStyle(),
		aiStyle:   styles.AIMessageStyle(),
		metaStyle: styles.MetadataStyle(),
	}
}
