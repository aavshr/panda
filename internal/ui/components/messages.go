package components

import (
	"github.com/aavshr/panda/internal/db"
	"github.com/aavshr/panda/internal/ui/styles"
	"github.com/charmbracelet/bubbles/list"
)

func NewMessagesModel(messages []*db.Message, width, height int) list.Model {
	items := NewMessageListItems(messages)
	model := list.New(items, NewMessageListItemDelegate(), width, height)
	model.Title = "Messages"
	model.Styles.Title = styles.DefaultListStyle()
	model.SetShowStatusBar(false)
	model.SetShowHelp(false)
	model.DisableQuitKeybindings()
	return model
}
