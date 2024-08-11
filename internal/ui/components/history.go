package components

import (
	"github.com/aavshr/panda/internal/db"
	"github.com/aavshr/panda/internal/ui/styles"
	"github.com/charmbracelet/bubbles/list"
)

func NewHistoryModel(threads []*db.Thread, width, height int) list.Model {
	items := NewThreadListItems(threads)
	model := list.New(items, NewThreadListItemDelegate(), width, height)
	model.Title = "History"
	model.Styles.Title = styles.DefaultListStyle()
	model.SetShowStatusBar(false)
	model.SetShowHelp(false)
	model.SetFilteringEnabled(false)
	return model
}
