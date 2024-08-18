package openai

import (
	"context"
	"errors"
	"io"

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

func (o *OpenAI) CreateChatCompletion(ctx context.Context, model, input string) (string, error) {
	if o.apiKey == "" {
		return "", ErrAPIKeyNotSet
	}
	resp, err := o.client.CreateChatCompletion(
		ctx,
		client.ChatCompletionRequest{
			Model: model,
			Messages: []client.ChatCompletionMessage{
				{
					Role:    client.ChatMessageRoleUser,
					Content: input,
				},
			},
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

func (o *OpenAI) CreateChatCompletionStream(ctx context.Context, model, input string) (io.ReadCloser, error) {
	if o.apiKey == "" {
		return nil, ErrAPIKeyNotSet
	}
	req := client.ChatCompletionRequest{
		Model: model,
		Messages: []client.ChatCompletionMessage{
			{
				Role:    client.ChatMessageRoleUser,
				Content: input,
			},
		},
		Stream: true,
	}
	stream, err := o.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return nil, err
	}
	return OpenAIStream{stream: stream}, nil
}
