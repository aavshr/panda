package openai

import (
	"context"
	"errors"
	"io"

	"github.com/aavshr/panda/internal/db"
	client "github.com/sashabaranov/go-openai"
)

const (
	defaultBaseURL = "https://api.openai.com/v1"
)

var (
	ErrAPIKeyNotSet      = errors.New("API key not set")
	ErrNoChoicesReturned = errors.New("no completion returned")
	ErrBufferTooSmall    = errors.New("buffer too small")
)

type OpenAIStream struct {
	stream *client.ChatCompletionStream
}

// TODO: verify that I can return EOF before copying the content
// TODO: what if len(p) is too small?
func (s OpenAIStream) Read(p []byte) (int, error) {
	resp, err := s.stream.Recv()
	if err != nil {
		return 0, err
	}
	if len(resp.Choices) == 0 {
		return 0, ErrNoChoicesReturned
	}
	content := resp.Choices[0].Delta.Content
	n := copy(p, content)
	if n < len(content) {
		return n, ErrBufferTooSmall
	}
	return n, nil
}

func (s OpenAIStream) Close() error {
	return s.stream.Close()
}

type OpenAI struct {
	baseURL string
	apiKey  string
	client  *client.Client
}

func New(baseURL string) *OpenAI {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	return &OpenAI{
		baseURL: baseURL,
	}
}

func (o *OpenAI) SetAPIKey(apiKey string) error {
	o.apiKey = apiKey
	o.client = client.NewClient(o.apiKey)
	return nil
}

func (o *OpenAI) dbMessagesToClientMessage(messages []*db.Message) []client.ChatCompletionMessage {
	var clientMessages []client.ChatCompletionMessage
	for _, m := range messages {
		m := m
		clientMessages = append(clientMessages, client.ChatCompletionMessage{
			Role:    m.Role,
			Content: m.Content,
		})
	}
	return clientMessages
}

// TODO: fix coupling with db message
func (o *OpenAI) CreateChatCompletion(ctx context.Context, model string, messages []*db.Message) (string, error) {
	if o.apiKey == "" {
		return "", ErrAPIKeyNotSet
	}
	resp, err := o.client.CreateChatCompletion(
		ctx,
		client.ChatCompletionRequest{
			Model:    model,
			Messages: o.dbMessagesToClientMessage(messages),
		},
	)
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", ErrNoChoicesReturned
	}
	return resp.Choices[0].Message.Content, nil
}

// TODO: coupling with db.Message
func (o *OpenAI) CreateChatCompletionStream(ctx context.Context, model string, messages []*db.Message) (io.ReadCloser, error) {
	if o.apiKey == "" {
		return nil, ErrAPIKeyNotSet
	}
	req := client.ChatCompletionRequest{
		Model:    model,
		Messages: o.dbMessagesToClientMessage(messages),
		Stream:   true,
	}
	stream, err := o.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return nil, err
	}
	return OpenAIStream{stream: stream}, nil
}
