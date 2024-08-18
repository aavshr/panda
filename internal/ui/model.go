package ui

import (
	"errors"
	"fmt"
	"io"

	"github.com/aavshr/panda/internal/config"
	"github.com/aavshr/panda/internal/db"
	"github.com/aavshr/panda/internal/ui/components"
	"github.com/aavshr/panda/internal/ui/llm"
	"github.com/aavshr/panda/internal/ui/store"
	"github.com/aavshr/panda/internal/ui/styles"

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
	titleMessages         = "Messages"
	titleHistory          = "History"
	timeFormat            = "2006-01-02 15:04:05"
	newThreadName         = "New"
	roleUser              = "user"
	roleAssistant         = "assistant"
	roleSystem            = "system"
)

type Config struct {
	InitThreadsLimit int
	MaxThreadsLimit  int
	MessagesLimit    int
	Width            int
	Height           int
	LLMModel         string

	historyWidth    int
	historyHeight   int
	messagesWidth   int
	messagesHeight  int
	chatInputWidth  int
	chatInputHeight int
}

type Model struct {
	conf         *Config
	userConfig   *config.Config
	showSettings bool

	messagesModel  components.ListModel
	historyModel   components.ListModel
	chatInputModel components.ChatInputModel
	settingsModel  components.SettingsModel

	threads           []*db.Thread
	threadsOffset     int
	activeThreadIndex int

	messages        []*db.Message
	messagesOffset  int
	activeLLMStream io.ReadCloser

	componentsToContainer map[components.Component]lipgloss.Style
	focusedComponent      components.Component
	selectedComponent     components.Component

	store store.Store
	llm   llm.LLM

	errorState error
}

func New(conf *Config, store store.Store, llm llm.LLM) (*Model, error) {
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
		llm:   llm,
	}
	m.settingsModel = components.NewSettingsModel()
	userConfig, err := config.Load()
	if err != nil {
		if !errors.Is(err, config.ErrConfigNotFound) {
			return m, fmt.Errorf("config.Load %w", err)
		}
		m.showSettings = true
	} else {
		m.userConfig = userConfig
		if err := m.llm.SetAPIKey(m.userConfig.LLMAPIKey); err != nil {
			return m, fmt.Errorf("llm.SetAPIKey %w", err)
		}
	}

	m.activeThreadIndex = 0
	m.threads = []*db.Thread{
		{
			Name: newThreadName,
			// TODO: not this hack lmao
			CreatedAt: "Create a new thread..",
		},
	}
	threads, err := m.store.ListLatestThreadsPaginated(0, m.conf.InitThreadsLimit)
	if err != nil {
		return m, fmt.Errorf("store.ListLatestThreadsPaginated %w", err)
	}
	m.threads = append(m.threads, threads...)

	m.messages = []*db.Message{}

	containerPaddingHeight := 18
	containerPaddingWidth := 10
	m.historyModel = components.NewListModel(
		titleHistory,
		components.NewThreadListItems(m.threads),
		conf.historyWidth-containerPaddingWidth,
		conf.historyHeight-containerPaddingHeight,
	)
	m.historyModel.Select(0) // New Thread is selected by default

	m.messagesModel = components.NewListModel(
		titleMessages,
		components.NewMessageListItems(m.messages),
		conf.messagesWidth-containerPaddingWidth,
		conf.messagesHeight-containerPaddingHeight,
	)
	m.chatInputModel = components.NewChatInputModel(conf.chatInputWidth, conf.chatInputHeight)

	container := styles.ContainerStyle()
	historyContainer := container.Copy().Width(m.conf.historyWidth).Height(m.conf.historyHeight)
	messagesContainer := container.Copy().Width(m.conf.messagesWidth).Height(m.conf.messagesHeight)
	chatInputContainer := container.Copy().Width(m.conf.chatInputWidth).Height(m.conf.chatInputHeight)
	m.componentsToContainer = map[components.Component]lipgloss.Style{
		components.ComponentHistory:   historyContainer,
		components.ComponentMessages:  messagesContainer,
		components.ComponentChatInput: chatInputContainer,
	}
	return m, nil
}

func (m *Model) setThreads(threads []*db.Thread) {
	m.threads = threads
	m.historyModel.SetItems(components.NewThreadListItems(threads))
}

func (m *Model) setMessages(messages []*db.Message) {
	m.messages = messages
	m.messagesModel.SetItems(components.NewMessageListItems(messages))
	m.messagesModel.GoToLastPage()
}

func (m *Model) setActiveThreadIndex(index int) {
	m.activeThreadIndex = index
	m.historyModel.Select(index)
}

func (m *Model) Init() tea.Cmd {
	if m.showSettings {
		m.focusedComponent = components.ComponentSettings
		m.selectedComponent = components.ComponentSettings
		return m.settingsModel.Focus()
	}
	m.settingsModel.Blur()
	m.focusedComponent = components.ComponentChatInput
	m.selectedComponent = components.ComponentChatInput
	return tea.Batch(
		m.chatInputModel.Focus(),
	)
}

func (m *Model) View() string {
	if m.errorState != nil {
		return fmt.Sprintf("Error: %v", m.errorState)
	}
	if m.showSettings {
		return m.settingsModel.View()
	}

	mainContainer := styles.MainContainerStyle()

	if container, ok := m.componentsToContainer[m.selectedComponent]; ok {
		styles.SetSelectedBorder(&container)
	}

	if container, ok := m.componentsToContainer[m.focusedComponent]; ok {
		styles.SetFocusedBorder(&container)
	}

	return mainContainer.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				m.componentsToContainer[components.ComponentHistory].Render(m.historyModel.View()),
				m.componentsToContainer[components.ComponentMessages].Render(m.messagesModel.View()),
			),
			lipgloss.JoinVertical(
				lipgloss.Left,
				m.componentsToContainer[components.ComponentChatInput].Render(m.chatInputModel.View()),
			),
		),
	)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.focusedComponent {
	case components.ComponentSettings:
		m.settingsModel, cmd = m.settingsModel.Update(msg)
	case components.ComponentHistory:
		m.historyModel, cmd = m.historyModel.Update(msg)
	case components.ComponentMessages:
		m.messagesModel, cmd = m.messagesModel.Update(msg)
	case components.ComponentChatInput:
		m.chatInputModel, cmd = m.chatInputModel.Update(msg)
	case components.ComponentNone:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			return m.handleKeyMsg(msg)
		}
	}

	switch msg := msg.(type) {
	case components.SettingsSubmitMsg:
		cmd = m.handleSettingsSubmitMsg(msg)
	case components.ChatInputReturnMsg:
		cmd = m.handleChatInputReturnMsg(msg)
	case components.EscapeMsg:
		m.handleEscapeMsg()
	case components.ListEnterMsg:
		cmd = m.handleListEnterMsg(msg)
	case components.ListSelectMsg:
		cmd = m.handleListSelectMsg(msg)
	case ForwardChatCompletionStreamMsg:
		cmd = m.handleForwardChatCompletionStreamMsg(msg)
	case error:
		m.errorState = msg
	}
	return m, cmd
}
