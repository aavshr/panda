package components

import (
	"github.com/aavshr/panda/internal/db"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"io"
)

// ThreadListItem implements the list.Item, list.DefaultItem and list.ItemDelegate interface
type ThreadListItem struct {
	thread *db.Thread
}

func (t *ThreadListItem) Title() string {
	return t.thread.Name
}

func (t *ThreadListItem) Description() string {
	return t.thread.CreatedAt
}

func (t *ThreadListItem) FilterValue() string {
	return t.thread.Name
}

func (t *ThreadListItem) Render(w io.Writer, m list.Model, index int, item list.Item) {
	// TODO: implement
	return
}

func (t *ThreadListItem) Height() int {
	// TODO: implement
	return 1
}

func (t *ThreadListItem) Spacing() int{
	//TODO: implement
	return 1
}

func (t *ThreadListItem) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	// TODO: implement
	return nil
}

func NewThreadListItems(threads []*db.Thread) []list.Item{
	items := make([]list.Item, len(threads))
	for i, t := range threads {
		items[i] = &ThreadListItem{
			thread: t,
		}
	}
	return items
}