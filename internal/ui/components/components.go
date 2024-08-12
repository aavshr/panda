package components

import (
	"github.com/charmbracelet/bubbles/list"
)

type Component string

const (
	ComponentHistory   Component = "history"
	ComponentMessages  Component = "messages"
	ComponentChatInput Component = "chatInput"
	ComponentNone      Component = "none" // utility component
)

type MsgEscape struct{}

func FocusListModel(model *list.Model) {
	model.FilterInput.Focus()
	model.Select(0)
}

func BlurListModel(model *list.Model) {
	model.FilterInput.Blur()
	model.Select(-1)
}
