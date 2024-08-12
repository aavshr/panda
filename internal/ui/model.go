package ui

import (
	"fmt"
	"time"

	"github.com/aavshr/panda/internal/db"
	"github.com/aavshr/panda/internal/ui/components"
	"github.com/aavshr/panda/internal/ui/store"
	"github.com/aavshr/panda/internal/ui/styles"

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

	componentsToContainer map[components.Component]lipgloss.Style
	focusedComponent      components.Component
	selectedComponent     components.Component

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

func (m *Model) Init() tea.Cmd {
	m.focusedComponent = components.ComponentChatInput
	m.selectedComponent = components.ComponentChatInput
	return tea.Batch(
		m.chatInputModel.Focus(),
	)
}

func (m *Model) View() string {
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

func (m *Model) handleKeyMsg(keyMsg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch keyMsg.String() {
	case "j", "up":
		switch m.selectedComponent {
		case components.ComponentChatInput:
			m.setSelectedComponent(components.ComponentMessages)
		}
	case "k", "down":
		switch m.selectedComponent {
		case components.ComponentMessages, components.ComponentHistory:
			m.setSelectedComponent(components.ComponentChatInput)
		}
	case "h", "left":
		switch m.selectedComponent {
		case components.ComponentMessages:
			m.setSelectedComponent(components.ComponentHistory)
		}
	case "l", "right":
		switch m.selectedComponent {
		case components.ComponentHistory:
			m.setSelectedComponent(components.ComponentMessages)
		}
	case "enter":
		m.setFocusedComponent(m.selectedComponent)
		return m, m.cmdFocusedComponent
	case "ctrl+c", "ctrl+d":
		return m, tea.Quit
	}
	return m, m.cmdSelectComponent
}

func (m *Model) setSelectedComponent(com components.Component) {
	if c, ok := m.componentsToContainer[com]; ok {
		if currentContainer, ok := m.componentsToContainer[m.selectedComponent]; ok {
			styles.SetNormalBorder(&currentContainer)
			m.componentsToContainer[m.selectedComponent] = currentContainer
		}
		m.selectedComponent = com
		styles.SetSelectedBorder(&c)
		m.componentsToContainer[com] = c
	}
}

func (m *Model) setFocusedComponent(com components.Component) {
	// focused component can be ComponentNone which won't be in the map
	m.focusedComponent = com
	if c, ok := m.componentsToContainer[com]; ok {
		styles.SetFocusedBorder(&c)
		m.componentsToContainer[com] = c

		switch com {
		case components.ComponentChatInput:
			m.chatInputModel.Focus()
		case components.ComponentMessages:
			components.FocusListModel(&m.messagesModel)
		case components.ComponentHistory:
			components.FocusListModel(&m.historyModel)
		}
	}
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.focusedComponent {
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
	case components.ChatInputEnterMsg:
		// TODO: store in db
		if msg.Value != "" {
			// TODO: send API request
			m.messages = append(m.messages, &db.Message{
				Role:      "user",
				Content:   msg.Value,
				CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
			})
			m.messagesModel.SetItems(components.NewMessageListItems(m.messages))
		}
	case components.MsgEscape:
		m.focusedComponent = components.ComponentNone
		switch m.focusedComponent {
		case components.ComponentChatInput:
			m.chatInputModel.Blur()
		case components.ComponentMessages:
			components.BlurListModel(&m.messagesModel)
		case components.ComponentHistory:
			components.BlurListModel(&m.historyModel)
		}
	}
	return m, cmd
}
