package components

import (
	"fmt"
	"io"

	"github.com/aavshr/panda/internal/db"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// MessageListItem implements the list.Item, list.DefaultItem and list.ItemDelegate interface
type MessageListItem struct {
	message *db.Message
}

func (t *MessageListItem) Title() string {
	return t.message.Content
}

func (t *MessageListItem) Description() string {
	role := "You"
	if t.message.Role == "assistant" {
		role = "AI"
	}
	return fmt.Sprintf("%s at %s", role, t.message.CreatedAt)
}

func (t *MessageListItem) FilterValue() string {
	return t.Title()
}

type MessageListItemDelegate struct {
	inner list.DefaultDelegate
}

// TODO: different rendering for different roles
func (d *MessageListItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	d.inner.Render(w, m, index, item)
}

func (d *MessageListItemDelegate) Height() int {
	// TODO: implement
	return 1
}

func (d *MessageListItemDelegate) Spacing() int {
	//TODO: implement
	return 1
}

func (t *MessageListItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	// TODO: implement
	return nil
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

func NewMessageListItemDelegate() list.ItemDelegate {
	return &MessageListItemDelegate{
		inner: list.NewDefaultDelegate(),
	}
}
