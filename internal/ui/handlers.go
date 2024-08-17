package ui

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/aavshr/panda/internal/db"
	"github.com/aavshr/panda/internal/ui/components"
	"github.com/aavshr/panda/internal/ui/styles"
	"github.com/aavshr/panda/internal/utils"
	tea "github.com/charmbracelet/bubbletea"
)

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
			m.messagesModel.Focus()
		case components.ComponentHistory:
			m.historyModel.Focus()
		}
	}
}

func (m *Model) createNewThread(name string) (*db.Thread, error) {
	newThreadID, err := utils.RandomID()
	if err != nil {
		return nil, fmt.Errorf("utils.RandomID: %w", err)
	}
	thread := &db.Thread{
		ID:        newThreadID,
		Name:      name,
		CreatedAt: time.Now().Format(timeFormat),
		UpdatedAt: time.Now().Format(timeFormat),
	}
	if err := m.store.UpsertThread(thread); err != nil {
		return thread, err
	}
	return thread, nil
}

func (m *Model) handleChatInputReturnMsg(msg components.ChatInputReturnMsg) error {
	msg.Value = strings.TrimSpace(msg.Value)
	if msg.Value == "" {
		return nil
	}
	if m.activeThreadIndex >= len(m.threads) {
		return fmt.Errorf("invalid active thread index")
	}

	// TODO: send API request
	// TODO: use thread name from api response
	// TODO: more robust behavior for thread creation
	// first thread is always for new thread
	if m.activeThreadIndex == 0 {
		n := 20
		if len(msg.Value) < n {
			n = len(msg.Value)
		}
		name := fmt.Sprintf("%s..", msg.Value[:n])
		newThread, err := m.createNewThread(name)
		if err != nil {
			return err
		}
		m.setThreads(slices.Insert(m.threads, 1, newThread))
		m.setActiveThreadIndex(1)
	}
	activeThread := m.threads[m.activeThreadIndex]
	message := &db.Message{
		Role:      "user",
		ThreadID:  activeThread.ID,
		Content:   msg.Value,
		CreatedAt: time.Now().Format(timeFormat),
	}
	if err := m.store.CreateMessage(message); err != nil {
		return err
	}
	m.setMessages(append(m.messages, message))
	return nil
}

func (m *Model) handleEscapeMsg() {
	m.focusedComponent = components.ComponentNone
	switch m.focusedComponent {
	case components.ComponentChatInput:
		m.chatInputModel.Blur()
	case components.ComponentMessages:
		m.messagesModel.Blur()
	case components.ComponentHistory:
		m.historyModel.Blur()
	}
}

func (m *Model) handleListEnterMsg(msg components.ListEnterMsg) error {
	switch m.focusedComponent {
	case components.ComponentHistory:
		m.setActiveThreadIndex(msg.Index)
		// first item is always for new thread
		if msg.Index == 0 {
			m.setMessages([]*db.Message{})
			m.setSelectedComponent(components.ComponentChatInput)
			m.setFocusedComponent(components.ComponentChatInput)
			return nil
		}
		if len(m.threads) == 0 || msg.Index >= len(m.threads) {
			return fmt.Errorf("invalid thread index")
		}
		threadId := m.threads[msg.Index].ID
		messages, err := m.store.ListMessagesByThreadIDPaginated(threadId, 0, m.conf.MessagesLimit)
		if err != nil {
			return err
		}
		m.setMessages(messages)
	// TODO: how should we handle selecting a message?
	case components.ComponentMessages:
	}
	return nil
}
