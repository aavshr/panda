package components

import (
	"github.com/aavshr/panda/internal/db"
	"github.com/aavshr/panda/internal/ui/styles"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

func NewHistoryModel(threads []*db.Thread, width, height int) list.Model {
	items := NewThreadListItems(threads)
	listPadding := width / 4
	model := list.New(items, NewThreadListItemDelegate(listPadding), width, height)
	model.Title = "History"
	model.Styles.Title = styles.DefaultListStyle().
		PaddingRight(listPadding).
		Border(lipgloss.DoubleBorder(), false, true, false, false)
	model.SetShowStatusBar(false)
	model.SetShowHelp(false)
	model.SetFilteringEnabled(false)
	return model
}
