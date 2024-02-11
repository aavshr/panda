package ui

import (
	"fmt"
	"github.com/aavshr/panda/internal/ui/components"
	"github.com/aavshr/panda/internal/ui/store"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Config struct {
	InitThreadsLimit int
	MaxThreadsLimit int
}

type Model struct {
	conf *Config

	threadsListModel list.Model

	threads []list.Item
	threadsOffset int

	store store.Store
}

func New(conf *Config, store store.Store) *Model {
	return &Model{
		conf: conf,
		store: store,
	}
}

func (m *Model) Init() tea.Cmd{
	threads, err := m.store.ListLatestThreadsPaginated(0, m.conf.InitThreadsLimit)
	if err != nil {
		fmt.Println("failed to list threads: ", err)
	}

	m.threads = components.NewThreadListItems(threads)
	m.threadsListModel = list.New(m.threads, list.NewDefaultDelegate(), 100, 100)
	return nil
}

func (m *Model) View() string {
	return m.threadsListModel.View()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO: implement
	return m, nil
}
