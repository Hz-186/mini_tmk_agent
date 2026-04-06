package tts

import "context"

type TTS interface {
	Synthesize(ctx context.Context, text string, lang string) ([]byte, error)
	GetVoiceForLang(lang string) string
}
