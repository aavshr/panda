package store

import "github.com/aavshr/panda/internal/db"

type Store interface {
	ListLatestThreadsPaginated(offset, limit int) ([]*db.Thread, error)
	ListLatestMessagesPaginated(threadID string, offset, limit int) ([]*db.Message, error)
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

func (m *Mock) ListLatestMessagesPaginated(threadID string, offset, limit int) ([]*db.Message, error) {
	return m.messages, nil
}
