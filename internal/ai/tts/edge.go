package tts

import (
	"context"
	"fmt"
	"os"

	"github.com/surfaceyu/edge-tts-go/edgeTTS"
)

/* ########
出于某种特殊的不可抗力，edge 的 TTS 不可使用
 ######## */

type EdgeTTS struct{}

func NewEdgeTTS() *EdgeTTS {
	return &EdgeTTS{}
}
func (e *EdgeTTS) GetVoiceForLang(lang string) string {
	switch lang {
	case "zh":
		return "zh-CN-XiaoxiaoNeural"
	case "en":
		return "en-US-AriaNeural"
	case "ja":
		return "ja-JP-NanamiNeural"
	case "es":
		return "es-ES-ElviraNeural"
	default:
		return "en-US-AriaNeural"
	}
}

func (e *EdgeTTS) Synthesize(ctx context.Context, text string, lang string) ([]byte, error) {
	voice := e.GetVoiceForLang(lang)

	tempFile, err := os.CreateTemp("", "tts_*.mp3")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	tempName := tempFile.Name()
	tempFile.Close()
	defer os.Remove(tempName)
	args := edgeTTS.Args{
		Text:       text,
		Voice:      voice,
		Volume:     "+0%",
		Rate:       "+0%",
		WriteMedia: tempName,
	}

	ttsCli := edgeTTS.NewTTS(args)
	if ttsCli == nil {
		return nil, fmt.Errorf("failed to init Edge TTS client")
	}
	ttsCli.AddTextDefault(text).Speak()
	audioData, err := os.ReadFile(tempName)
	if err != nil {
		return nil, fmt.Errorf("failed to read synthesized audio: %w", err)
	}
	return audioData, nil
}
