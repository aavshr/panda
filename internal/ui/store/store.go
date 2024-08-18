package store

import (
	"github.com/aavshr/panda/internal/db"
)

type Store interface {
	ListLatestThreadsPaginated(offset, limit int) ([]*db.Thread, error)
	ListMessagesByThreadIDPaginated(threadID string, offset, limit int) ([]*db.Message, error)
	UpsertThread(thread *db.Thread) error
	UpdateThreadName(threadID, name string) error
	DeleteThread(threadID string) error
	DeleteAllThreads() error
	CreateMessage(message *db.Message) error
}

type Mock struct {
	threads  []*db.Thread
	messages map[string][]*db.Message
}

func NewMock(threads []*db.Thread, messages []*db.Message) *Mock {
	msgs := make(map[string][]*db.Message)
	for _, message := range messages {
		msgs[message.ThreadID] = append(msgs[message.ThreadID], message)
	}
	return &Mock{
		threads:  threads,
		messages: msgs,
	}
}

func (m *Mock) ListLatestThreadsPaginated(offset, limit int) ([]*db.Thread, error) {
	return m.threads, nil
}

func (m *Mock) ListMessagesByThreadIDPaginated(threadID string, offset, limit int) ([]*db.Message, error) {
	messages, _ := m.messages[threadID]
	return messages, nil
}

func (m *Mock) UpsertThread(thread *db.Thread) error {
	for i, t := range m.threads {
		if t.ID == thread.ID {
			m.threads[i] = thread
			return nil
		}
	}
	m.threads = append(m.threads, thread)
	m.messages[thread.ID] = []*db.Message{}
	return nil
}

func (m *Mock) UpdateThreadName(threadID, name string) error {
	for _, thread := range m.threads {
		if thread.ID == threadID {
			thread.Name = name
			return nil
		}
	}
	return nil
}

func (m *Mock) DeleteThread(threadID string) error {
	for i, thread := range m.threads {
		if thread.ID == threadID {
			m.threads = append(m.threads[:i], m.threads[i+1:]...)
			return nil
		}
	}
	return nil
}

func (m *Mock) DeleteAllThreads() error {
	m.threads = []*db.Thread{}
	return nil
}

func (m *Mock) CreateMessage(message *db.Message) error {
	if msgs, ok := m.messages[message.ThreadID]; ok {
		m.messages[message.ThreadID] = append(msgs, message)
	}
	return nil
}
