package llm

import (
	"io"
	"strings"
)

type LLM interface {
	CreateChatCompletion(string, string) (string, error)
	CreateChatCompletionStream(string, string) (io.Reader, error)
}

type Mock struct{}

func NewMock() *Mock {
	return &Mock{}
}

func (m *Mock) CreateChatCompletion(model, input string) (string, error) {
	return "this is a mock AI response", nil
}

func (m *Mock) CreateChatCompletionStream(model, input string) (io.Reader, error) {
	return strings.NewReader("this is a mock AI response.\nwith a new line."), nil
}
