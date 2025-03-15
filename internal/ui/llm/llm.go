package llm

import (
	"context"
	"github.com/aavshr/panda/internal/db"
	"io"
	"strings"
)

// TODO: fix coupling with db message
type LLM interface {
	CreateChatCompletion(context.Context, string, []*db.Message) (string, error)
	CreateChatCompletionStream(context.Context, string, []*db.Message) (io.ReadCloser, error)
	SetAPIKey(string) error
}

type Mock struct{}

func NewMock() *Mock {
	return &Mock{}
}

func (m *Mock) CreateChatCompletion(ctx context.Context, model, input string) (string, error) {
	return "this is a mock AI response", nil
}

func (m *Mock) CreateChatCompletionStream(ctx context.Context, model, input string) (io.ReadCloser, error) {
	return io.NopCloser(strings.NewReader("this is a mock AI response.\nwith a new line.")), nil
}

func (m *Mock) SetAPIKey(string) error {
	return nil
}
