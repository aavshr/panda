package components

import (
	"strings"

	"github.com/aavshr/panda/internal/ui/styles"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type ListEnterMsg struct {
	Index int
}

func ListEnterCmd(selectedIndex int) func() tea.Msg {
	return func() tea.Msg {
		return ListEnterMsg{
			Index: selectedIndex,
		}
	}
}

type ListModel struct {
	inner list.Model
}

func NewListModel(title string, items []list.Item, width, height int) ListModel {
	model := list.New(items, NewMessageListItemDelegate(), width, height)
	model.Title = title
	model.Styles.Title = styles.DefaultListStyle()
	model.SetShowStatusBar(false)
	model.SetShowHelp(false)
	model.FilterInput.Blur()
	model.InfiniteScrolling = false
	// no item should be selected by default
	model.Select(-1)
	model.Styles.NoItems.Padding(0, 0, 1, 2)

	// TODO: what if title is not plural
	model.SetStatusBarItemName(strings.TrimSuffix(strings.ToLower(title), "s"), strings.ToLower(title))

	// disable quit key
	model.KeyMap.Quit.SetEnabled(false)

	return ListModel{inner: model}
}

func (m *ListModel) Focus() {
	m.inner.FilterInput.Focus()
	m.inner.Select(0)
}

func (m *ListModel) Blur() {
	m.inner.FilterInput.Blur()
	m.inner.Select(-1)
}

func (m *ListModel) Select(i int) {
	m.inner.Select(i)
}

func (m *ListModel) SetItems(items []list.Item) {
	m.inner.SetItems(items)
}

func (m *ListModel) View() string {
	return m.inner.View()
}

func (m *ListModel) Update(msg tea.Msg) (ListModel, tea.Cmd) {
	var cmd tea.Cmd
	m.inner, cmd = m.inner.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		// TODO: handle enter key
		case tea.KeyEscape:
			return *m, EscapeCmd
		case tea.KeyEnter:
			return *m, ListEnterCmd(m.inner.Index())
		}
	case error:
		// TODO: how to handle errors
		return *m, cmd
	}
	return *m, cmd
}
