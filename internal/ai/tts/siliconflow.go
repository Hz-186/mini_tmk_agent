package tts

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"project_for_tmk_04_06/internal/config"
)

type SiliconFlowTTS struct{}

func NewSiliconFlowTTS() *SiliconFlowTTS {
	return &SiliconFlowTTS{}
}

func (s *SiliconFlowTTS) GetVoiceForLang(lang string) string {
	return "FunAudioLLM/CosyVoice2-0.5B:alex"
}

func (s *SiliconFlowTTS) Synthesize(ctx context.Context, text string, lang string) ([]byte, error) {
	key := config.AppConfig.AI.Key
	if key == "" {
		return nil, fmt.Errorf("SiliconFlow API key is required for TTS")
	}

	payload := map[string]interface{}{
		"model":           "FunAudioLLM/CosyVoice2-0.5B",
		"input":           text,
		"voice":           s.GetVoiceForLang(lang),
		"response_format": "mp3",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.siliconflow.cn/v1/audio/speech", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("TTS API error, status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	audioData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return audioData, nil
}
