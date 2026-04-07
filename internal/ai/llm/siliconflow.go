package llm

import (
	"context"
	"errors"
	"fmt"
	"io"
	"project_for_tmk_04_06/internal/config"

	"github.com/sashabaranov/go-openai"
)

// 大语言模型
// Large Language Model

type SiliconFlowLLM struct {
	client *openai.Client
	model  string
}

func NewSiliconFlowLLM() *SiliconFlowLLM {
	cfg := openai.DefaultConfig(config.AppConfig.AI.Key)
	cfg.BaseURL = config.AppConfig.AI.BaseURL

	return &SiliconFlowLLM{
		client: openai.NewClientWithConfig(cfg),
		model:  "Qwen/Qwen2.5-7B-Instruct",
	}
}

func (s *SiliconFlowLLM) buildPrompt(sourceText, sourceLang, targetLang string) []openai.ChatCompletionMessage {
	return []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: fmt.Sprintf("You are an expert bilingual simultaneous interpreter. Translate the user's input from %s to %s fluently. Only output the translated text without any conversational fillers, markdown formatting, or explanations.", sourceLang, targetLang),
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: sourceText,
		},
	}
}

func (s *SiliconFlowLLM) Translate(ctx context.Context, sourceText, sourceLang, targetLang string) (string, error) {
	req := openai.ChatCompletionRequest{
		Model:    s.model,
		Messages: s.buildPrompt(sourceText, sourceLang, targetLang),
	}

	resp, err := s.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("LLM translation failed: %w", err)
	}

	if len(resp.Choices) > 0 {
		return resp.Choices[0].Message.Content, nil
	}
	return "", errors.New("empty response from LLM")
}

func (s *SiliconFlowLLM) TranslateStream(ctx context.Context, sourceText, sourceLang, targetLang string) (<-chan string, error) {
	req := openai.ChatCompletionRequest{
		Model:    s.model,
		Messages: s.buildPrompt(sourceText, sourceLang, targetLang),
		Stream:   true,
	}

	stream, err := s.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("LLM stream setup failed: %w", err)
	}

	outChan := make(chan string)

	go func() {
		defer stream.Close()
		defer close(outChan)

		for {
			resp, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				return
			}
			if err != nil {
				return
			}
			if len(resp.Choices) > 0 {
				outChan <- resp.Choices[0].Delta.Content
			}
		}
	}()

	return outChan, nil
}
