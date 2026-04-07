package asr

import "context"

// 自动语音识别
// Automatic Speech Recognition

type ASR interface {
	Transcribe(ctx context.Context, audioFilepath string) (string, error)
	TranscribeBytes(ctx context.Context, audioData []byte, lang string) (string, error)
}
