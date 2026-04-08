package asr

// 自动语音识别

import (
	"bytes"
	"context"
	"fmt"
	"project_for_tmk_04_06/internal/config"

	"github.com/sashabaranov/go-openai"
)

type SiliconFlowASR struct {
	client *openai.Client
	model  string
}

func NewSiliconFlowASR() *SiliconFlowASR {
	cfg := openai.DefaultConfig(config.AppConfig.AI.Key)
	cfg.BaseURL = config.AppConfig.AI.BaseURL
	return &SiliconFlowASR{
		client: openai.NewClientWithConfig(cfg),
		model:  "FunAudioLLM/SenseVoiceSmall",
	}
}

func (s *SiliconFlowASR) Transcribe(ctx context.Context, audioFilepath string) (string, error) {
	req := openai.AudioRequest{
		Model:    s.model,
		FilePath: audioFilepath,
	}
	resp, err := s.client.CreateTranscription(ctx, req)
	if err != nil {
		return "", fmt.Errorf("SiliconFlow ASR transcription failed: %w", err)
	}
	return resp.Text, nil
}

func (s *SiliconFlowASR) TranscribeBytes(ctx context.Context, audioData []byte, lang string) (string, error) {
	req := openai.AudioRequest{
		Model:    s.model,
		Reader:   bytes.NewReader(audioData),
		FilePath: "audio.wav", // Dummy filename required by API
	}
	resp, err := s.client.CreateTranscription(ctx, req)
	if err != nil {
		return "", fmt.Errorf("SiliconFlow ASR transcription (bytes) failed: %w", err)
	}
	return resp.Text, nil
}
