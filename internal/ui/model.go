package ui

import (
	"fmt"

	"github.com/aavshr/panda/internal/db"
	"github.com/aavshr/panda/internal/ui/components"
	"github.com/aavshr/panda/internal/ui/store"
	"github.com/aavshr/panda/internal/ui/styles"

	//"github.com/aavshr/panda/internal/ui/styles"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	/*
		| history  |  messages  |
		|          |            |
		|          |            |
		|----------|------------|
		| settings |  prompt    |
	*/
	widthSeparationRatio  = 0.3
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
	sidebarWidth     int
	historyHeight    int
	messagesWidth    int
	messagesHeight   int
	settingsHeight   int
	chatInputWidth   int
	chatInputHeight  int
}

type Model struct {
	conf *Config

	messagesModel  list.Model
	historyModel   list.Model
	chatInputModel textarea.Model

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
	conf.sidebarWidth = int(float64(conf.Width) * widthSeparationRatio)
	conf.historyHeight = conf.Height - int(float64(conf.Height)*heightSeparationRatio)
	conf.messagesWidth = conf.Width - conf.sidebarWidth
	conf.messagesHeight = conf.Height - conf.chatInputHeight
	conf.settingsHeight = conf.Height - conf.historyHeight
	conf.chatInputWidth = conf.messagesWidth
	conf.chatInputHeight = conf.settingsHeight

	m := &Model{
		conf:  conf,
		store: store,
	}
	threads, err := m.store.ListLatestThreadsPaginated(0, m.conf.InitThreadsLimit)
	if err != nil {
		return m, fmt.Errorf("store.ListLatestThreadsPaginated %w", err)
	}
	m.threads = threads
	m.historyModel = components.NewHistoryModel(m.threads, conf.sidebarWidth, conf.historyHeight)
	m.messagesModel = components.NewMessagesModel(m.messages, conf.messagesWidth, conf.messagesHeight)
	m.initChatInputModel()
	return m, nil
}

func (m *Model) initChatInputModel() {
	chatInputModel := textarea.New()
	chatInputModel.Placeholder = "Type..."
	chatInputModel.SetWidth(m.conf.chatInputWidth)
	chatInputModel.SetHeight(m.conf.chatInputHeight)
	chatInputModel.Focus()
	//chatInputModel.Prompt = "> "
	chatInputModel.FocusedStyle.CursorLine = lipgloss.NewStyle()
	chatInputModel.ShowLineNumbers = false
	m.chatInputModel = chatInputModel
}

func (m *Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m *Model) View() string {
	mainContainer := styles.MainContainerStyle()

	return mainContainer.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				m.historyModel.View(),
				m.messagesModel.View(),
			),
			m.chatInputModel.View(),
		),
	)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO: implement
	return m, nil
}
