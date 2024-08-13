package store

import "github.com/aavshr/panda/internal/db"

type Store interface {
	ListLatestThreadsPaginated(offset, limit int) ([]*db.Thread, error)
	ListMessagesByThreadIDPaginated(threadID string, offset, limit int) ([]*db.Message, error)
	CreateThread(thread *db.Thread) error
	UpdateThreadName(threadID, name string) error
	DeleteThread(threadID string) error
	DeleteAllThreads() error
	CreateMessage(message *db.Message) error
}

type Mock struct {
	threads  []*db.Thread
	messages []*db.Message
}

func NewMock(threads []*db.Thread, messages []*db.Message) *Mock {
	return &Mock{
		threads:  threads,
		messages: messages,
	}
}

func (m *Mock) ListLatestThreadsPaginated(offset, limit int) ([]*db.Thread, error) {
	return m.threads, nil
}

func (m *Mock) ListMessagesByThreadIDPaginated(threadID string, offset, limit int) ([]*db.Message, error) {
	return m.messages, nil
}

func (m *Mock) CreateThread(thread *db.Thread) error {
	m.threads = append(m.threads, thread)
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
	m.messages = append(m.messages, message)
	return nil
}
