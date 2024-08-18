package ui

import (
	"context"
	"errors"
	"fmt"
	"io"
	"slices"
	"time"

	"github.com/aavshr/panda/internal/config"
	"github.com/aavshr/panda/internal/db"
	"github.com/aavshr/panda/internal/ui/components"
	"github.com/aavshr/panda/internal/ui/styles"
	"github.com/aavshr/panda/internal/utils"
	tea "github.com/charmbracelet/bubbletea"
	"golang.design/x/clipboard"
)

func (m *Model) handleKeyMsg(keyMsg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch keyMsg.String() {
	case "k", "up":
		switch m.selectedComponent {
		case components.ComponentChatInput:
			m.setSelectedComponent(components.ComponentHistory)
		}
	case "j", "down":
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
			// TODO: check why we need to reslect the activethreadindex
			m.historyModel.Select(m.activeThreadIndex)
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

func (m *Model) handleSettingsSubmitMsg(msg components.SettingsSubmitMsg) tea.Cmd {
	savedConfig, err := config.Save(config.Config{
		LLMAPIKey: msg.APIKey,
	})
	if err != nil {
		return m.cmdError(fmt.Errorf("config.Save: %w", err))
	}
	if err := m.llm.SetAPIKey(msg.APIKey); err != nil {
		return m.cmdError(fmt.Errorf("llm.SetAPIKey: %w", err))
	}
	m.showSettings = false
	m.userConfig = savedConfig
	return m.Init()
}

func (m *Model) handleChatInputReturnMsg(msg components.ChatInputReturnMsg) tea.Cmd {
	if m.activeThreadIndex >= len(m.threads) {
		return m.cmdError(fmt.Errorf("invalid active thread index"))
	}

	// TODO: use llm to generate thread name as well

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
			return m.cmdError(fmt.Errorf("createNewThread: %w", err))
		}
		m.setThreads(slices.Insert(m.threads, 1, newThread))
		m.setActiveThreadIndex(1)
	}
	activeThread := m.threads[m.activeThreadIndex]
	userMessage := &db.Message{
		Role:      roleUser,
		ThreadID:  activeThread.ID,
		Content:   msg.Value,
		CreatedAt: time.Now().Format(timeFormat),
	}
	if err := m.store.CreateMessage(userMessage); err != nil {
		return m.cmdError(fmt.Errorf("store.CreateMessage: %w", err))
	}
	m.setMessages(append(m.messages, userMessage))
	// TODO: history?
	reader, err := m.llm.CreateChatCompletionStream(context.Background(),
		m.userConfig.LLMModel, msg.Value)
	if err != nil {
		return m.cmdError(fmt.Errorf("llm.CreateChatCompletionStream: %w", err))
	}
	m.activeLLMStream = reader

	// placeholder empty llm message for stream to update as data rolls in
	// message will only be saved to db when stream is done
	m.setMessages(append(m.messages, &db.Message{
		Role:     roleAssistant,
		ThreadID: activeThread.ID,
	}))
	return m.cmdForwardChatCompletionStream
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

func (m *Model) handleListEnterMsg(msg components.ListEnterMsg) tea.Cmd {
	// TODO: handle entering a message for messages
	// either copy the entire message or let user copy only specific parts
	switch m.focusedComponent {
	case components.ComponentHistory:
		// first item is always for new thread
		if msg.Index == 0 {
			m.setMessages([]*db.Message{})
			m.setSelectedComponent(components.ComponentChatInput)
			m.setFocusedComponent(components.ComponentChatInput)
			return nil
		}
		m.setSelectedComponent(components.ComponentMessages)
		m.setFocusedComponent(components.ComponentMessages)
	case components.ComponentMessages:
		// TODO: how to convey to the user that the message is copied or there was an error
		if err := clipboard.Init(); err == nil {
			if msg.Index >= len(m.messages) {
				return m.cmdError(fmt.Errorf("invalid message index"))
			}
			message := m.messages[msg.Index].Content
			clipboard.Write(clipboard.FmtText, []byte(message))
		}
	}
	return nil
}

func (m *Model) handleListSelectMsg(msg components.ListSelectMsg) tea.Cmd {
	switch m.focusedComponent {
	case components.ComponentHistory:
		m.setActiveThreadIndex(msg.Index)

		if len(m.threads) == 0 || msg.Index >= len(m.threads) {
			return m.cmdError(fmt.Errorf("invalid thread index"))
		}
		threadId := m.threads[msg.Index].ID
		messages, err := m.store.ListMessagesByThreadIDPaginated(threadId, 0, m.conf.MessagesLimit)
		if err != nil {
			return m.cmdError(err)
		}
		m.setMessages(messages)
	// TODO: how should we handle selecting a message?
	case components.ComponentMessages:
	}
	return nil
}

func (m *Model) handleForwardChatCompletionStreamMsg(msg ForwardChatCompletionStreamMsg) tea.Cmd {
	if m.activeThreadIndex >= len(m.threads) {
		return m.cmdError(fmt.Errorf("invalid active thread index"))
	}
	activeThreadId := m.threads[m.activeThreadIndex].ID
	llmMessageIndex := len(m.messages) - 1
	if llmMessageIndex == 0 {
		return m.cmdError(fmt.Errorf("bad llm message index: 0, should be at least 1"))
	}
	content := m.messages[llmMessageIndex].Content
	createdAt := m.messages[llmMessageIndex].CreatedAt

	// TODO: what buffer size makes it look smooth?
	buffer := make([]byte, 16)

	streamDone := false
	n, err := m.activeLLMStream.Read(buffer)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			return m.cmdError(fmt.Errorf("activeLLMStream.Read: %w", err))
		}
		streamDone = true
		m.activeLLMStream.Close()
	}
	/*
		if n == 0 && !streamDone {
			return m.cmdError(fmt.Errorf("activeLLMStream.Read: no bytes read"))
		}
	*/
	if n > 0 {
		content = fmt.Sprintf("%s%s", content, string(buffer[:n]))
		// upate created at as soon as first bytes are read
		if createdAt == "" {
			createdAt = time.Now().Format(timeFormat)
		}
	}
	updatedLLMMessage := &db.Message{
		Role:      roleAssistant,
		Content:   content,
		CreatedAt: createdAt,
		ThreadID:  activeThreadId,
	}
	m.messages[llmMessageIndex] = updatedLLMMessage
	setItemCmd := m.messagesModel.SetItem(
		llmMessageIndex,
		components.NewMessageListItem(updatedLLMMessage),
	)
	if streamDone {
		if err := m.store.CreateMessage(updatedLLMMessage); err != nil {
			return m.cmdError(fmt.Errorf("store.CreateMessage: %w", err))
		}
		return setItemCmd
	}
	return tea.Batch(setItemCmd, m.cmdForwardChatCompletionStream)
}
