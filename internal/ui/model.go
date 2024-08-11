package ui

import (
	"fmt"

	"github.com/aavshr/panda/internal/db"
	"github.com/aavshr/panda/internal/ui/components"
	"github.com/aavshr/panda/internal/ui/store"
	"github.com/aavshr/panda/internal/ui/styles"

	//"github.com/aavshr/panda/internal/ui/styles"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	/*
		| history  |  messages  |
		|          |            |
		|          |            |
		|----------|------------|
		|  chat input           |
	*/
	widthSeparationRatio  = 0.2
	heightSeparationRatio = 0.1
)

const (
	titleMessages = "Messages"
	titleSettings = "Settings"
	titlePrompt   = "Prompt"
)

type Config struct {
	InitThreadsLimit int
	MaxThreadsLimit  int
	Width            int
	Height           int
	historyWidth     int
	historyHeight    int
	messagesWidth    int
	messagesHeight   int
	chatInputWidth   int
	chatInputHeight  int
}

type Model struct {
	conf *Config

	messagesModel  list.Model
	historyModel   list.Model
	chatInputModel components.ChatInputModel
	focusedModel   string
	selectedModel  string

	threads       []*db.Thread
	threadsOffset int

	messages       []*db.Message
	messagesOffset int

	store store.Store
}

func New(conf *Config, store store.Store) (*Model, error) {
	if conf.Width == 0 || conf.Height == 0 {
		return nil, fmt.Errorf("invalid config: width and height must be greater than 0")
	}
	conf.historyWidth = int(float64(conf.Width) * widthSeparationRatio)
	conf.historyHeight = conf.Height - int(float64(conf.Height)*heightSeparationRatio)
	conf.messagesWidth = conf.Width - conf.historyWidth
	conf.messagesHeight = conf.historyHeight
	conf.chatInputWidth = conf.Width + 2 // to account for border and padding
	conf.chatInputHeight = conf.Height - conf.historyHeight

	m := &Model{
		conf:  conf,
		store: store,
	}
	threads, err := m.store.ListLatestThreadsPaginated(0, m.conf.InitThreadsLimit)
	if err != nil {
		return m, fmt.Errorf("store.ListLatestThreadsPaginated %w", err)
	}
	m.threads = threads
	m.historyModel = components.NewHistoryModel(m.threads, conf.historyWidth, conf.historyHeight)
	m.messagesModel = components.NewMessagesModel(m.messages, conf.messagesWidth, conf.messagesHeight)
	m.chatInputModel = components.NewChatInputModel(conf.chatInputWidth, conf.chatInputHeight)
	return m, nil
}

func (m *Model) Init() tea.Cmd {
	m.focusedModel = "chat_input"
	m.selectedModel = "chat_input"
	return tea.Batch(
		m.chatInputModel.Focus(),
	)
}

func (m *Model) View() string {
	mainContainer := styles.MainContainerStyle()

	container := styles.ContainerStyle()
	historyContainer := container.Copy().Width(m.conf.historyWidth).Height(m.conf.historyHeight)
	messagesContainer := container.Copy().Width(m.conf.messagesWidth).Height(m.conf.messagesHeight)
	chatInputContainer := container.Copy().Width(m.conf.chatInputWidth).Height(m.conf.chatInputHeight)

	switch m.selectedModel {
	case "history":
		styles.SetActiveBorder(&historyContainer)
	case "messages":
		styles.SetActiveBorder(&messagesContainer)
	case "chat_input":
		styles.SetActiveBorder(&chatInputContainer)
	}

	return mainContainer.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				historyContainer.Render(m.historyModel.View()),
				messagesContainer.Render(m.messagesModel.View()),
			),
			lipgloss.JoinVertical(
				lipgloss.Left,
				chatInputContainer.Render(m.chatInputModel.View()),
			),
		),
	)
}

func (m *Model) handleKeyMsg(keyMsg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch keyMsg.String() {
	case "j", "up":
		switch m.selectedModel {
		case "chat_input":
			m.selectedModel = "messages"
		}
	case "k", "down":
		switch m.selectedModel {
		case "messages", "history":
			m.selectedModel = "chat_input"
		}
	case "h", "left":
		switch m.selectedModel {
		case "messages":
			m.selectedModel = "history"
		}
	case "l", "right":
		switch m.selectedModel {
		case "history":
			m.selectedModel = "messages"
		}
	case "ctrl+c", "ctrl+d":
		return m, tea.Quit
	}
	return m, nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.focusedModel {
	case "history":
		m.historyModel, cmd = m.historyModel.Update(msg)
	case "messages":
		m.messagesModel, cmd = m.messagesModel.Update(msg)
	case "chat_input":
		m.chatInputModel, cmd = m.chatInputModel.Update(msg)
	case "none":
		switch msg := msg.(type) {
		case tea.KeyMsg:
			return m.handleKeyMsg(msg)
		}
	}

	switch msg := msg.(type) {
	case components.ChatInputEnterMsg:
		// TODO: store in db
		if msg.Value != "" {
			// TODO: send API request
			m.messages = append(m.messages, &db.Message{Content: msg.Value})
			m.messagesModel.SetItems(components.NewMessageListItems(m.messages))
		}
	case components.EscapeMsg:
		m.focusedModel = "none"
		m.chatInputModel.Blur()
	}
	return m, cmd
}
