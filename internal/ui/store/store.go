package store

import "github.com/aavshr/panda/internal/db"

type Store interface {
	ListLatestThreadsPaginated(offset, limit int) ([]*db.Thread, error)
}

type Mock struct {
	threads []*db.Thread
}

func NewMock(threads []*db.Thread) *Mock {
	return &Mock{
		threads: threads,
	}
}

func (m *Mock) ListLatestThreadsPaginated(offset, limit int) ([]*db.Thread, error) {
	return m.threads, nil
}