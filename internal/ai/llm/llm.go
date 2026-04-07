package llm

// 大语言模型
// Large Language Model

import "context"

type LLM interface {
	Translate(ctx context.Context, sourceText, sourceLang, targetLang string) (string, error)
	TranslateStream(ctx context.Context, sourceText, sourceLang, targetLang string) (<-chan string, error)
}
