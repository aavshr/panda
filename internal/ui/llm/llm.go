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

func (m *Mock) CreateChatCompletion(input string, model string) (string, error) {
	return "this is a mock AI response", nil
}

func (m *Mock) CreateChatCompletionStream(input string, model string) (io.Reader, error) {
	return strings.NewReader("this is a mock AI response"), nil
}
