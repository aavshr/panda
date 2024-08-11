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

type ThreadListItemDelegate struct {
	inner list.DefaultDelegate
}

func (d *ThreadListItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	d.inner.Render(w, m, index, item)
}

func (d *ThreadListItemDelegate) Height() int {
	// TODO: implement
	return d.inner.Height()
}

func (d *ThreadListItemDelegate) Spacing() int {
	//TODO: implement
	return d.inner.Spacing()
}

func (t *ThreadListItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	// TODO: implement
	return nil
}

func NewThreadListItems(threads []*db.Thread) []list.Item {
	items := make([]list.Item, len(threads))
	for i, t := range threads {
		items[i] = &ThreadListItem{
			thread: t,
		}
	}
	return items
}

func NewThreadListItemDelegate() list.ItemDelegate {
	return &ThreadListItemDelegate{
		inner: list.NewDefaultDelegate(),
	}
}

