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
	return tea.Batch(
		m.chatInputModel.Init(),
		m.chatInputModel.Focus(),
	)
}

func (m *Model) View() string {
	mainContainer := styles.MainContainerStyle()

	container := styles.ContainerStyle()
	historyContainer := container.Copy().Width(m.conf.historyWidth).Height(m.conf.historyHeight)
	messagesContainer := container.Copy().Width(m.conf.messagesWidth).Height(m.conf.messagesHeight)
	chatInputContainer := container.Copy().Width(m.conf.chatInputWidth).Height(m.conf.chatInputHeight)

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

// TODO: how does the Update loop work?
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var historyCmd, messagesCmd, chatInputCmd tea.Cmd
	m.chatInputModel, chatInputCmd = m.chatInputModel.Update(msg)

	// TODO: store in db
	if msg, ok := msg.(components.ChatInputEnterMsg); ok {
		// TODO: not use SetItems directly?
		// is there a better way to do this? bubble tea recommends not to use Cmd like this
		if msg.Value != "" {
			// TODO: send API request
			m.messages = append(m.messages, &db.Message{Content: msg.Value})
			m.messagesModel.SetItems(components.NewMessageListItems(m.messages))
		}
	}
	m.messagesModel, messagesCmd = m.messagesModel.Update(msg)
	m.historyModel, historyCmd = m.historyModel.Update(msg)
	return m, tea.Batch(historyCmd, messagesCmd, chatInputCmd)
}
